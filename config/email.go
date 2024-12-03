package config

type Email struct {
	Smtp SmtpConfig `yaml:"smtp" json:"smtp" mapstructure:"smtp"`
}

type SmtpConfig struct {
	Host     string `yaml:"host" json:"host" mapstructure:"host"`
	Port     int    `yaml:"port" json:"port" mapstructure:"port"`
	Username string `yaml:"username" json:"username" mapstructure:"username"`
	Password string `yaml:"password" json:"password" mapstructure:"password"`
	From     string `yaml:"from" json:"from" mapstructure:"from"`
	Tls      bool   `yaml:"tls" json:"tls" mapstructure:"tls"`
}
