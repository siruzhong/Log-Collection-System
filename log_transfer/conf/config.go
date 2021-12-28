package conf

// Config 总配置结构体
type Config struct {
	ESConfig    `ini:"elasticsearch"`
	KafkaConfig `ini:"kafka"`
}

// ESConfig elasticsearch配置结构体
type ESConfig struct {
	Address     string `ini:"address"`
	Index       string `ini:"index"`
	MaxChanSize int    `ini:"max_chan_size"`
}

// KafkaConfig kafka配置结构体
type KafkaConfig struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}
