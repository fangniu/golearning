package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"strconv"
	"strings"
	"log"
	"net/http"
	"encoding/json"

	"github.com/prometheus/common/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/bsm/sarama-cluster"
)

var (
	brokers    string
	topics      string
	addr       string
	alarmGroup string
	promGroup  string
)

type Service struct {
	ServiceId   int     `json:"service_id"`
	ServiceName string  `json:"service_name"`
	AttrName    string  `json:"attr_name"`
	Value       float64 `json:"value"`
	Ip          string  `json:"ip"`
}

type Exporter struct {
	alarm       *cluster.Consumer
	prom        *cluster.Consumer
	serviceDesc *prometheus.Desc
}

func NewExporter(topics []string, brokers []string) (*Exporter, error) {
	config := cluster.NewConfig()
	prom, err := cluster.NewConsumer(brokers, promGroup, topics, config)
	if err != nil {
		return nil, err
	}
	alarm, err := cluster.NewConsumer(brokers, alarmGroup, topics, config)
	if err != nil {
		return nil, err
	}

	exporter := Exporter{
		prom:  prom,
		alarm: alarm,
		serviceDesc: prometheus.NewDesc(
			"kafka_service_monitor",
			"服务监控", [] string{"service_id", "service_name", "attr_name", "ip"},
			nil,
		),
	}

	return &exporter, nil
}

func (exporter *Exporter) alert() {
	go func() {
		for err := range exporter.alarm.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	go func() {
		for msg := range exporter.alarm.Messages() {
			fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		}
	}()

}

func (exporter *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- exporter.serviceDesc
}

func (exporter *Exporter) Collect(ch chan<- prometheus.Metric) {
	timeout := true
	for timeout {
		select {
		case msg := <-exporter.prom.Messages():
			var s Service
			err := json.Unmarshal(msg.Value, &s)
			if err != nil {
				log.Println("ERROR:", err)
				break
			}
			ch <- prometheus.MustNewConstMetric(
				exporter.serviceDesc,
				prometheus.GaugeValue,
				s.Value,
				strconv.Itoa(s.ServiceId),
				s.ServiceName,
				s.AttrName,
				s.Ip,
			)

		case <-time.After(time.Millisecond * 50):
			timeout = false
		}
	}
}

func init() {
	flag.StringVar(&brokers, "brokers", "localhost:9092", "Kafka地址，用逗号隔开")
	flag.StringVar(&topics, "topics", "", "Kafka中的Topics，用逗号隔开")
	flag.StringVar(&addr, "addr", ":8000", "prometheus端口")
	flag.StringVar(&promGroup, "promGroup", "prometheus", "prometheus消费组ID")
	flag.StringVar(&alarmGroup, "alarmGroup", "alarm", "alarm消费组ID")
	prometheus.MustRegister(version.NewCollector("kafka_exporter"))
}

func main() {
	flag.Parse()
	if topics == "" {
		fmt.Println("未指定topics！")
		os.Exit(1)
	}

	exporter, err := NewExporter(strings.Split(topics, ","), strings.Split(brokers, ","))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("连接kafka成功: ", brokers, topics)

	defer exporter.alarm.Close()
	defer exporter.prom.Close()
	exporter.alert()
	prometheus.MustRegister(exporter)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
