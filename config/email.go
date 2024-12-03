package config

import "github.com/wneessen/go-mail"

type Email struct {
	Smtp SmtpConfig `yaml:"smtp" json:"smtp" mapstructure:"smtp"`
}

type SmtpConfig struct {
	Host     string `yaml:"host" json:"host" mapstructure:"host"`
	Port     int    `yaml:"port" json:"port" mapstructure:"port"`
	Username string `yaml:"username" json:"username" mapstructure:"username"`
	Password string `yaml:"password" json:"password" mapstructure:"password"`
	From     string `yaml:"from" json:"from" mapstructure:"from"`
	Tls      string `yaml:"tls" json:"tls" mapstructure:"tls"`
	Auth     string `yaml:"auth" json:"auth" mapstructure:"auth"`
}

func (c *SmtpConfig) TlsPolicy() mail.TLSPolicy {
	switch c.Tls {
	case "mandatory":
		return mail.TLSMandatory
	case "opportunistic":
		return mail.TLSOpportunistic
	case "no_tls":
		return mail.NoTLS
	default:
		return mail.NoTLS
	}
}

func (c *SmtpConfig) SMTPAuth() mail.SMTPAuthType {
	switch c.Auth {
	case "plain":
		return mail.SMTPAuthPlain
	case "plain_noenc":
		return mail.SMTPAuthPlainNoEnc
	case "login":
		return mail.SMTPAuthLogin
	case "login_noenc":
		return mail.SMTPAuthLoginNoEnc
	case "cram_md5":
		return mail.SMTPAuthCramMD5
	case "scram_sha1":
		return mail.SMTPAuthSCRAMSHA1
	case "scram_sha1_plus":
		return mail.SMTPAuthSCRAMSHA1PLUS
	case "scram_sha256":
		return mail.SMTPAuthSCRAMSHA256
	case "scram_sha256_plus":
		return mail.SMTPAuthSCRAMSHA256PLUS
	case "xoauth2":
		return mail.SMTPAuthXOAUTH2
	case "noauth":
		return mail.SMTPAuthNoAuth
	default:
		return mail.SMTPAuthNoAuth
	}
}
