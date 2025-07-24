package config

import "github.com/spf13/viper"

var Configuration *Config

type Config struct {
	AppName         string         `json:"appName"`
	Version         string         `json:"version"`
	Env             string         `json:"env"`
	CpuCount        int            `json:"cpuCount"`
	WorkId          int64          `json:"WorkId"`
	Health          HealthConf     `mapstructure:"health"`
	Restful         RestfulConf    `mapstructure:"restful"`
	Ws              WSConf         `mapstructure:"ws"`
	Tcp             TCPConf        `mapstructure:"tcp"`
	Grpc            GrpcConf       `mapstructure:"grpc"`
	Etcd            EtcdConf       `mapstructure:"etcd"`
	Mongo           MongoConf      `mapstructure:"mongo"`
	Postgresql      PostgresqlConf `mapstructure:"postgresql"`
	Redis           RedisConf      `mapstructure:"redis"`
	JWT             JWTConf        `mapstructure:"jwt"`
	ConnectionTotal int            `json:"connectionTotal"`
	Distribute      DistributeConf `mapstructure:"distribute"`
	Nats            NatsConf       `mapstructure:"nats"`
	MessageNodes    []MessageNode  `mapstructure:"messageNodes"`
	Log             LogConf        `mapstructure:"log"`
	AliYun          AliYunConf     `mapstructure:"aliYun"`
}

type HealthConf struct {
	Port int `json:"port"`
}

type RestfulConf struct {
	Url  string `json:"url"`
	Port int    `json:"port"`
}

type WSConf struct {
	Url  string `json:"url"`
	Port int    `json:"port"`
}

type TCPConf struct {
	Url  string `json:"url"`
	Port int    `json:"port"`
}

type EtcdConf struct {
	Url      string `json:"url"`
	Port     int    `json:"port"`
	DialTime int    `json:"dialTime"`
	Ttl      int    `json:"ttl"`
}

type GrpcConf struct {
	Url  string `json:"url"`
	Port int    `json:"port"`
}

type MongoConf struct {
	Url      string `json:"url"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostgresqlConf struct {
	Url         string `json:"url"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Db          string `json:"db"`
	MaxIdleSize int    `json:"maxIdleSize"`
	MaxOpenSize int    `json:"maxOpenSize"`
}

type RedisConf struct {
	Url         string `json:"url"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Db          int    `json:"db"`
	PoolSize    int    `json:"poolSize"`
	MaxIdleSize int    `json:"maxIdleSize"`
	MinIdleSize int    `json:"minIdleSize"`
}

type JWTConf struct {
	Secret string `json:"secret"`
}

type DistributeConf struct {
	WorkerSize int `json:"workerSize"`
	QueueSize  int `json:"queueSize"`
}

type NatsConf struct {
	Url  string `json:"url"`
	Port int    `json:"port"`
}

type MessageNode struct {
	Key    string `json:"key"`
	Weight int    `json:"weight"`
}

type LogConf struct {
	Level string `json:"level"`
}

type AliYunConf struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SmsSignName     string `json:"smsSignName"`
	SmsEndpoint     string `json:"smsEndpoint"`
}

func InitConfig(configFile string) {
	Configuration = new(Config)
	v := viper.New()
	v.SetConfigFile(configFile)
	v.ReadInConfig()
	v.Unmarshal(&Configuration)
}
