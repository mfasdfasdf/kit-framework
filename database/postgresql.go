package database

import (
	"fmt"
	"github.com/mfasdfasdf/kit-framework/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var PostgresqlCli *PostgresqlClient

type PostgresqlClient struct {
	Cli *gorm.DB
}

func InitPostgresql() {
	if PostgresqlCli != nil {
		return
	}
	url := config.Configuration.Postgresql.Url
	port := config.Configuration.Postgresql.Port
	username := config.Configuration.Postgresql.Username
	password := config.Configuration.Postgresql.Password
	db := config.Configuration.Postgresql.Db
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=Asia/Shanghai", url, port, username, password, db)
	var client *gorm.DB = nil
	if config.Configuration.Env == "dev" {
		client, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{LogLevel: logger.Info, Colorful: true})})
	} else {
		client, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	}
	sqlDB, _ := client.DB()
	sqlDB.SetMaxIdleConns(config.Configuration.Postgresql.MaxIdleSize)
	sqlDB.SetMaxOpenConns(config.Configuration.Postgresql.MaxOpenSize)
	PostgresqlCli = &PostgresqlClient{Cli: client}
}
