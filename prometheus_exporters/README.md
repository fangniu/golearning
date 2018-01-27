# 从kafka获取监控数据转发给prometheus
2组kafka消费者，支持多个partition
* alarm， 实时告警的消费组
* prometheus， 离线消费给prometheus， 如20秒从kafka消费一次
```
go run kafka_exporter.go --brokers 192.168.101.58:9092,192.168.101.58:9093 --topics topic1,topic2
```
