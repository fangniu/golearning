package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"time"
	"sync"
	"io/ioutil"
	"strings"
	"encoding/json"
	"database/sql"

	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"bufio"
	"github.com/golang/protobuf/proto"
	"errors"
)

type Cgk struct {
	Host           string `json:"host"`
	User           string `json:"user"`
	Passwd         string `json:"passwd"`
	Db             string `json:"db"`
	Port           int    `json:"port"`
	Charset        string `json:"charset"`
	ConnectTimeout int    `json:"connect_timeout"`
}

type Influxdb struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type GlobalConfig struct {
	ServerTypeId   int      `json:"server_type_id"`
	LocalIp        string   `json:"local_ip"`
	ExecFile       string   `json:"exec_file"`
	Interval       int      `json:"interval"`
	Timeout        int      `json:"timeout"`
	Cgks           []Cgk    `json:"cgks"`
	Influxdb       Influxdb `json:"influxdb"`
	InfluxdbTags   []string `json:"influxdb_tags"`
	InfluxdbFields []string `json:"influxdb_fields"`
}

type Agent struct {
	ip   string
	port int
	addr string
}

type Monitor struct {
	points []*client.Point
}

var (
	config         *GlobalConfig
	cgkConnections []*sql.DB
	infConnection  client.Client
	wg             sync.WaitGroup
	execTimeout		time.Duration
	requestUdp = []byte{0, 4, 1, 0, 0, 0, 0}
	PackageHeaderError = errors.New("package header error")
	PackageError = errors.New("package error")
)

func initMysqlConns() {
	var err error
	infConnection, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%d", config.Influxdb.Host, config.Influxdb.Port),
		Username: config.Influxdb.Username,
		Password: config.Influxdb.Password,
	})
	if err != nil {
		log.Fatalln("ERROR fail to connect influxdb: ", config.Influxdb.Host, config.Influxdb.Port)
	}
}

func (m *Monitor) initInfluxdb() {
	q := client.NewQuery("CREATE DATABASE "+config.Influxdb.Database, "", "")
	response, err := infConnection.Query(q)
	if err != nil {
		log.Fatalln("ERROR fail to connect influxdb:", config.Influxdb.Host, config.Influxdb.Port)
	}
	if response.Error() != nil {
		log.Fatalln("ERROR fail to create database:", response.Error())
	}

	q = client.NewQuery(fmt.Sprintf("CREATE RETENTION POLICY agent_retention ON %s DURATION 30d REPLICATION 1 DEFAULT",
		config.Influxdb.Database), "", "")
	response, err = infConnection.Query(q)
	if err != nil {
		log.Fatalln("ERROR fail to create retention:", err)
	}
	if response.Error() != nil {
		log.Fatalln("ERROR fail to create retention:", response.Error())
	}

}

func getAgents() [] *Agent {
	var agents [] *Agent
	for _, cgk := range cgkConnections {
		rows, err := cgk.Query("SELECT ip_b, port_b FROM t_server WHERE server_type_id = ?", config.ServerTypeId)
		if err != nil {
			log.Fatalln("ERROR fail to connect cgk:", err)
		}
		defer rows.Close()
		for rows.Next() {
			var a Agent
			err := rows.Scan(&a.ip, &a.port)
			if err != nil {
				log.Fatalln("ERROR fail to parse agent info:", err)
			}
			a.addr = fmt.Sprintf("%v:%v", a.ip, a.port)
			agents = append(agents, &a)
		}
	}

	return agents
}

func (m *Monitor) close() {
	for _, c := range cgkConnections {
		c.Close()
	}
	infConnection.Close()
}

func (m *Monitor) runCollect() {
	for {
		for _, agent := range getAgents() {
			wg.Add(1)
			go m.collect(agent)
		}
		wg.Wait()
		m.report()
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

func (m *Monitor) report() {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
		Database:  config.Influxdb.Database,
	})
	if err != nil {
		log.Fatalln("ERROR fail to create bp:", err)
	}
	bp.SetRetentionPolicy("agent_retention")
	bp.AddPoints(m.points)
	if err = infConnection.Write(bp); err != nil {
		log.Fatalln("ERROR fail to report points:", err)
	}
	m.points = nil
}

func (m *Monitor) getAgentData(agent *Agent) (err error) {
	var conn net.Conn
	var n int
	conn, err = net.Dial("udp", agent.addr)
	if err != nil {
		return
	}
	defer conn.Close()
	reader := bufio.NewReaderSize(conn, 8192)
	_, err = conn.Write(requestUdp)
	if err != nil {
		return
	}
	resp :=  make([]byte, 7)
	n, err = reader.Read(resp)
	if err != nil {
		return
	}

	if n == 7 && resp[1] == 5 && resp[2] == 1 {
		var length int
		for _, v := range resp[3:] {
			length = length*256 + int(v)
		}
		var assr AgentServerStatusResponse
		buf := make([]byte, length)
		n, err = reader.Read(buf)
		if err != nil {
			return
		}
		if n != length {
			err = PackageHeaderError
			return
		}
		err = proto.Unmarshal(buf, &assr)
		if err != nil {
			return
		}
		tags := map[string]string{}
		fields := map[string]interface{}{}
		tags["str_local_ip"] = assr.GetStrLocalIp()
		tags["local_port"] = fmt.Sprintf("%d", assr.GetLocalPort())
		tags["queue_name"] = assr.GetQueueName()
		tags["zk_hosts"] = assr.GetZkHosts()
		fields["queue_ele_size"] = assr.GetQueueEleSize()
		fields["queue_ele_count"] = assr.GetQueueEleCount()
		fields["process_count"] = int(assr.GetProcessCount())
		fields["process_per_sec"] = int(assr.GetProcessPerSec())
		fields["proxy_size"] = assr.GetProxySize()
		fields["attr_size"] = assr.GetAttrSize()
		var pt *client.Point
		pt, err = client.NewPoint("agent_status", tags, fields, time.Now())
		if err != nil {
			return
		}
		m.points = append(m.points, pt)
	}
	return PackageError
}

func (m *Monitor) collect(agent *Agent) {
	defer wg.Done()
	done := make(chan error, 1)

	go func() {
		done <- m.getAgentData(agent)
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Println("WARN failed to get agent data: ", err, agent.addr)
		}
	case <-time.After(execTimeout):
		log.Println("WARN timeout", agent.ip, agent.port)
	}
}

func parseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		log.Fatalln("ERROR config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	data, err := ioutil.ReadFile(cfg)
	if err != nil {
		log.Fatalln("ERROR read config file:", cfg, "fail:", err)
	}
	configContent := strings.TrimSpace(string(data))
	err = json.Unmarshal([]byte(configContent), &config)
	if err != nil {
		log.Fatalln("ERROR parse config file:", cfg, "fail:", err)
	}
	execTimeout = time.Duration(config.Timeout) * time.Second
	log.Println("INFO read config file:", cfg, "successfully")
	for _, cgk := range config.Cgks {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", cgk.User, cgk.Passwd, cgk.Host, cgk.Port, cgk.Db, cgk.Charset)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalln("ERROR Mysql connection error:", cgk, "fail:", err)
		}
		cgkConnections = append(cgkConnections, db)
	}
}

func main() {
	cfg := flag.String("c", "config.json", "configuration file")
	flag.Parse()
	parseConfig(*cfg)
	initMysqlConns()
	m := Monitor{}
	m.initInfluxdb()
	defer m.close()
	m.runCollect()
}
