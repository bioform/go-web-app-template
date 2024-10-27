package config

type Database struct {
	Dsn string `yaml:"dsn" json:"dsn" mapstructure:"dsn"`
}
