package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"time"
	"sync"
	"os/exec"
	"io/ioutil"
	"strings"
	"bytes"
	"strconv"
	"encoding/json"
	"database/sql"

	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/go-sql-driver/mysql"
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
}

type Monitor struct {
	bp     client.BatchPoints
	points []*client.Point
}

var (
	config         *GlobalConfig
	cgkConnections []*sql.DB
	infConnection  client.Client
	wg             sync.WaitGroup
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

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
		Database:  config.Influxdb.Database,
	})
	bp.SetRetentionPolicy("agent_retention")
	m.bp = bp
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
	for true {
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
	m.bp.AddPoints(m.points)
	if err := infConnection.Write(m.bp); err != nil {
		log.Fatalln("ERROR fail to report points:", err)
	}
	m.points = nil
}

func stopCollect(ch chan bool)  {
	ch <- true
}

func (m *Monitor) collect(agent *Agent) {
	defer wg.Done()
	ch := make(chan bool)
	go func() {
		defer stopCollect(ch)
		tags := map[string]string{}
		fields := map[string]interface{}{}
		cmd := exec.Command(config.ExecFile, config.LocalIp, agent.ip, strconv.Itoa(agent.port))
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Println("WARN", err)
			return
		}
		results := map[string]string{}
		outStr := strings.TrimSpace(out.String())
		if outStr == "" {
			log.Println("WARN output invalid:", agent.ip, agent.port)
			return
		}
		for _, line := range strings.Split(out.String(), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			slice := strings.SplitN(line, ":", 2)
			if len(slice) != 2 {
				log.Println("WARN split output:", line)
				return
			}
			results[slice[0]] = slice[1]
		}
		for _, tag := range config.InfluxdbTags {
			tags[tag] = results[tag]
		}
		for _, field := range config.InfluxdbFields {
			v, ok := results[field]
			if !ok {
				log.Printf("WARN output format error: %s", results)
				return
			}
			i, err := strconv.Atoi(v)
			if err != nil {
				log.Printf("WARN output format to int: %s[%s]", field, results[field])
				return
			}
			fields[field] = i
		}
		pt, err := client.NewPoint("agent_status", tags, fields, time.Now())
		if err != nil {
			log.Fatalln("ERROR new influxdb point:", err)
			return
		}
		m.points = append(m.points, pt)
	}()

	select {
	case <-ch:
	case <-time.After(time.Duration(config.Timeout) * time.Second):
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
