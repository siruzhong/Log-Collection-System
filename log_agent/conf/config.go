package conf

// Config 总配置结构体
type Config struct {
	KafkaConfig `ini:"kafka"`
	EtcdConfig  `ini:"etcd"`
}

// KafkaConfig kafka配置结构体
type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int    `ini:"chan_size"`
}

// EtcdConfig etcd配置结构体
type EtcdConfig struct {
	Address string `ini:"address"`
	Key     string `ini:"key"`
}

// LogConfig 日志配置项结构体
type LogConfig struct {
	Path  string `json:"path"`  // 日志文件存放路径
	Topic string `json:"topic"` // 日志文件写入kafka中的topic
}
