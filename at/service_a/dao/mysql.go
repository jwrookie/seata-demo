package dao

import (
	"fmt"
	"github.com/opentrx/mysql/v2"
	dialector "service_a/dialector/mysql"
	"strings"
	"time"

	_ "github.com/opentrx/mysql/v2"
	"github.com/opentrx/seata-golang/v2/pkg/client/config"
	"github.com/opentrx/seata-golang/v2/pkg/client/rm"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

type gormLogger struct{}

func (*gormLogger) Printf(format string, v ...interface{}) {
	format = strings.Replace(format, "\n", " ", 1)
	fmt.Println(fmt.Sprintf(format, v))
}

// NewMysql connect to scripts
func NewMysql() {
	var err error

	rm.RegisterTransactionServiceServer(mysql.GetDataSourceManager())
	mysql.RegisterResource(config.GetATConfig().DSN)

	newLogger := logger.New(
		&gormLogger{},
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	mysqlConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: newLogger,
	}

	db, err = gorm.Open(
		dialector.Open(config.GetATConfig().DSN),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			}}, &mysqlConfig)
	if err != nil {
		panic(err)
	}
	DB, err := db.DB()
	if err != nil {
		panic(err)
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(20)
	DB.SetConnMaxLifetime(4 * time.Hour)
}

// Mysql get a connection for scripts
func Mysql() *gorm.DB {
	return db
}

// DisconnectMysql disconnect scripts
func DisconnectMysql() {
	mysqlDB, _ := db.DB()
	if err := mysqlDB.Close(); err != nil {
		panic(err)
	}
}
