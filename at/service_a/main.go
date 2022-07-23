package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentrx/mysql/v2"
	"github.com/opentrx/seata-golang/v2/pkg/client"
	"github.com/opentrx/seata-golang/v2/pkg/client/config"
	"github.com/opentrx/seata-golang/v2/pkg/util/log"
	"service_a/dao"
	"strings"
)

type gormLogger struct{}

func (*gormLogger) Printf(format string, v ...interface{}) {
	format = strings.Replace(format, "\n", " ", 1)
	fmt.Println(fmt.Sprintf(format, v))
}

func main() {
	conf := config.InitConfiguration("./conf/client.yml")
	log.Init(conf.Log.LogPath, conf.Log.LogLevel)
	client.Init(conf)

	dao.NewMysql()

	route := gin.New()
	route.POST("/jobA", func(c *gin.Context) {

		tx := dao.Mysql().WithContext(context.WithValue(
			context.Background(),
			mysql.XID,
			c.Request.Header.Get("XID"),
		)).Begin(nil)

		if err := (&dao.ADao{}).Create(tx); err != nil {
			tx.Rollback()
			c.String(500, "%s", err.Error())
			return
		}

		tx.Commit()
		c.String(200, "%s", "succ")
	})

	if err := route.Run(":7001"); err != nil {
		dao.DisconnectMysql()
		panic(err)
	}
	dao.DisconnectMysql()
}
