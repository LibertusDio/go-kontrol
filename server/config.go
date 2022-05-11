package main

import (
	"fmt"
)

// Config holds all settings of finportal
type Config struct {
	AppName     string `yaml:"app_name" mapstructure:"app_name"`
	HTTPPort    string `yaml:"port" mapstructure:"port"`
	MySQL       *MySQL `yaml:"mysql" mapstructure:"mysql"`
	Environment string `yaml:"environment" mapstructure:"environment"`
	TokenTTL    int64  `yaml:"token_ttl" mapstructure:"token_ttl"`
	LogLevel    uint8  `yaml:"log_level" mapstructure:"log_level"`
}

// MySQL ...
type MySQL struct {
	Host               string `yaml:"host" mapstructure:"host"`
	Port               int    `yaml:"port" mapstructure:"port"`
	Username           string `yaml:"username" mapstructure:"username"`
	Password           string `yaml:"password" mapstructure:"password"`
	Database           string `yaml:"database" mapstructure:"database"`
	SSLMode            string `yaml:"sslmode" mapstructure:"sslmode"`
	Timeout            int    `yaml:"timeout" mapstructure:"timeout"`
	ConnectionMax      int    `yaml:"connection_max" mapstructure:"connection_max"`
	ConnectionTime     int64  `yaml:"connection_time" mapstructure:"connection_time"`
	ConnectionIdleMax  int    `yaml:"connection_idle_max" mapstructure:"connection_idle_max"`
	ConnectionIdleTime int64  `yaml:"connection_idle_time" mapstructure:"connection_idle_time"`
	Log                bool   `yaml:"log" mapstructure:"log"`
}

// ConnectionString ...
func (c *MySQL) ConnectionString() string {
	if c.Timeout == 0 {
		c.Timeout = 15
	}

	if c.ConnectionIdleMax < 0 {
		c.ConnectionIdleMax = 0
	}

	if c.ConnectionMax <= 0 {
		c.ConnectionMax = 20
	}

	if c.ConnectionIdleTime <= 0 {
		c.ConnectionIdleTime = 10
	}

	if c.ConnectionTime < 0 {
		c.ConnectionTime = 0
	}

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Database)
}

// Postgres ...
type Postgres struct {
	Host        string `yaml:"host" mapstructure:"host"`
	Port        int    `yaml:"port" mapstructure:"port"`
	Username    string `yaml:"username" mapstructure:"username"`
	Password    string `yaml:"password" mapstructure:"password"`
	Database    string `yaml:"database" mapstructure:"database"`
	SSLMode     string `yaml:"sslmode" mapstructure:"sslmode"`
	SSLRootCert string `yaml:"sslrootcert" mapstructure:"sslrootcert"`
	Timeout     int    `yaml:"timeout" mapstructure:"timeout"`
}

// ConnectionString ...
func (c *Postgres) ConnectionString() string {
	if c.Timeout == 0 {
		c.Timeout = 15
	}

	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v sslrootcert=%v connect_timeout=%v", c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode, c.SSLRootCert, c.Timeout)
}

// FormatDSN ...
func (c *Postgres) FormatDSN() string {
	if c.Timeout == 0 {
		c.Timeout = 15
	}

	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v&sslrootcert=%v&connect_timeout=%v", c.Username, c.Password, c.Host, c.Port, c.Database, c.SSLMode, c.SSLRootCert, c.Timeout)
}

// Default config
const ConfigDefault string = `
app_name: "authority"
environment: development
port: "80"
token_ttl: 1800
log_level: 1
mysql: 
 database: kontrol
 host: 127.0.0.1
 password: mysql
 port: 3306
 sslmode: disable
 timeout: 15
 username: root
 connection_max: 20
 connection_time: 300
 connection_idle_max: 0
 connection_idle_time: 10
 log: true
`

// Auto testing config
const ConfigTest string = `
app_name: "authority"
environment: testing
port: "4445"
log_level: 1
token_ttl: 1800
mysql: 
 database: auth_db
 host: localhost
 password: mysql
 port: 5432
 sslmode: disable
 timeout: 15
 username: mysql
 connection_max: 20
 connection_time: 300
 connection_idle_max: 0
 connection_idle_time: 10
 log: true
`
