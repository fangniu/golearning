# 从kafka获取监控数据转发给prometheus

go run kafka_exporter.go --brokers 192.168.101.58:9092,192.168.101.58:9093 --topics topic1,topic2
