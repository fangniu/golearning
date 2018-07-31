package main

import (
	"flag"
	"os"
	"time"
	"strconv"
	"strings"
	"log"
	"fmt"
	"sync"
	"encoding/json"

	"github.com/bsm/sarama-cluster"
	"github.com/Shopify/sarama"
	"github.com/influxdata/influxdb/client/v2"
	"io/ioutil"
)

type Service struct {
	ServiceId   int     `json:"service_id"`
	ServiceName string  `json:"service_name"`
	AttrName    string  `json:"attr_name"`
	Value       float64 `json:"value"`
	Ip          int     `json:"repotor_ip"`
	OPType      int     `json:"op_type"`
	BeginTime   int64   `json:"begin_time"`
}

func (s *Service) ip() string {
	return fmt.Sprintf("%d.%d.%d.%d", s.Ip%C256A, s.Ip/C256A%C256A, s.Ip/C256B%C256A, s.Ip/C256C%C256A)
}

type Kafka struct {
	Brokers []string `json:"brokers"`
	Topics  []string `json:"topics"`
	Oldest  bool     `json:"oldest"`
	GroupId string   `json:"group_id"`
}

type InfluxDB struct {
	Addr        string `json:"addr"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Database    string `json:"database"`
	Measurement string `json:"measurement"`
}

type Config struct {
	Kafka    Kafka    `json:"kafka"`
	InfluxDB InfluxDB `json:"influxdb"`
	Limit    int      `json:"limit"`
}

type Transfer struct {
	influxdbClient client.Client
	bps            []client.BatchPoints
	kafkaClient    *cluster.Consumer
	buffer         chan *client.Point
	debug          bool
	curMsg         *sarama.ConsumerMessage
	lastMsg        *sarama.ConsumerMessage
}

func (t *Transfer) connect() {
	var err error
	err = t.connectInf()
	if err != nil {
		log.Fatal("ERROR 连接influxdb失败")
	}
	kafkaConfig := cluster.NewConfig()
	//kafkaConfig.Group.Mode = cluster.ConsumerModePartitions
	if config.Kafka.Oldest {
		kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	t.kafkaClient, err = cluster.NewConsumer(config.Kafka.Brokers, config.Kafka.GroupId, config.Kafka.Topics, kafkaConfig)
	if err != nil {
		log.Fatal("ERROR 连接kafka失败： ", err)
	}
	log.Println("连接kafka成功:", config.Kafka.Brokers)
}

func (t *Transfer) connectInf() (err error){
	t.influxdbClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.InfluxDB.Addr,
		Username: config.InfluxDB.Username,
		Password: config.InfluxDB.Password,
	})
	if err != nil {
		log.Println("ERROR 连接influxdb失败： ", err)
		return
	}
	q := client.NewQuery("CREATE DATABASE "+config.InfluxDB.Database, "", "")
	response, err := t.influxdbClient.Query(q)
	if err != nil {
		log.Println("ERROR fail to connect influxdb:", config.InfluxDB.Addr)
		return
	}
	if response.Error() != nil {
		log.Println("ERROR fail to create database:", response.Error())
		err = response.Error()
		return
	}

	q = client.NewQuery(fmt.Sprintf("CREATE RETENTION POLICY service_retention ON %s DURATION 30d REPLICATION 1 DEFAULT",
		config.InfluxDB.Database), "", "")
	response, err = t.influxdbClient.Query(q)
	if err != nil {
		log.Println("ERROR fail to create retention:", err)
		return
	}
	if response.Error() != nil {
		log.Println("ERROR fail to create retention:", response.Error())
		err = response.Error()
		return
	}
	log.Println("连接influxdb成功:", config.InfluxDB.Addr)
	return
}


func (t *Transfer) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go t.consume()
	go t.markOffset()
	go t.combinePoints()
	wg.Wait()
}

func (t *Transfer) consume() {
	for msg := range t.kafkaClient.Messages() {
		var s Service
		err := json.Unmarshal(msg.Value, &s)
		if err != nil {
			log.Println("WARN Unmarshal:", string(msg.Value))
			continue
		}
		tags := map[string]string{
			"service_name": s.ServiceName,
			"service_id":   strconv.Itoa(s.ServiceId),
			"attr_name":    s.AttrName,
			"op_type":      strconv.Itoa(s.OPType),
			"repotor_ip":   s.ip(),
		}
		values := map[string]interface{}{
			"value": s.Value,
		}
		pt, err := client.NewPoint(config.InfluxDB.Measurement, tags, values, time.Unix(s.BeginTime, 0))
		if err != nil {
			log.Fatalln("ERROR new influxdb point:", err)
			return
		}
		t.curMsg = msg
		t.buffer <- pt
	}
}

func (t *Transfer) markOffset() {
	for {
		time.Sleep(time.Second * 2)
		if t.curMsg != nil && t.curMsg != t.lastMsg {
			t.kafkaClient.MarkOffset(t.curMsg, "")
			t.lastMsg = t.curMsg
		}
	}

}

func (t *Transfer) combinePoints() {
	for {
		point := <-t.buffer
		var points []*client.Point
		points = append(points, point)
		length := len(t.buffer)
		for i:=0; i < length && i < config.Limit - 1; i++ {
			points = append(points, <-t.buffer)
		}
		t.write(points)
	}
}

func (t *Transfer) write(points []*client.Point) {
	if points == nil {
		return
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.InfluxDB.Database,
		Precision: "s",
	})
	if err != nil {
		log.Fatalln("ERROR influxdb 创建BatchPoints失败", err)
	}
	bp.AddPoints(points)
	start := time.Now()
	if err := t.influxdbClient.Write(bp); err != nil {
		log.Println("ERROR fail to write points:", err)
		t.reconnectInf()
	}
	if t.debug {
		log.Println("INFO writed points count: ", len(points), time.Now().Sub(start))
	}
}

func (t *Transfer) reconnectInf() {
	for {
		err := t.connectInf()
		if err != nil {
			time.Sleep(time.Second * 10)
		} else {
			return
		}
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
}

var (
	config Config
)

const C256A = 256
const C256B = C256A * C256A
const C256C = C256B * C256A

func main() {
	cfg := flag.String("c", "config.json", "configuration file")
	debug := flag.Bool("d", false, "debug mode")
	flag.Parse()
	parseConfig(*cfg)

	transfer := Transfer{
		debug:  *debug,
		buffer: make(chan *client.Point, config.Limit*10),
	}
	transfer.connect()
	transfer.run()
}
