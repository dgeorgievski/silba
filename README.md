# silba
System metrics collector in Go based on cloudfoundry/gosigar. Metrics are encoded using MessagePack, 
and streamed via ZeroMQ. Clients should integrate silba with streaming technologies like Fluentd, or Heka.

... In the meantime InfluxDB has released [telegraf](https://github.com/influxdb/telegraf) a Go agent that implements many of the features planned for silba, especially in terms of plugable architecture. This project will be retired and used only as a playing ground for marshalling serialized data over the network using Go. 
