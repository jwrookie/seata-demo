package main

import (
	"context"
	"net/http"
	"service_b/dao"

	"github.com/opentrx/seata-golang/v2/pkg/client"
	"github.com/opentrx/seata-golang/v2/pkg/client/config"
	"github.com/opentrx/seata-golang/v2/pkg/util/log"

	"github.com/gin-gonic/gin"
	"github.com/opentrx/mysql/v2"
)

func main() {
	conf := config.InitConfiguration("./conf/client.yml")
	log.Init(conf.Log.LogPath, conf.Log.LogLevel)
	client.Init(conf)

	dao.NewMysql()

	route := gin.New()
	route.POST("/jobB", func(c *gin.Context) {
		type data struct {
			Rollback bool `json:"rollback"`
		}
		var d data
		if err := c.ShouldBindJSON(&d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 模拟发生错误
		if d.Rollback {
			c.String(500, "%s", "fail")
			return
		}

		tx := dao.Mysql().WithContext(context.WithValue(
			context.Background(),
			mysql.XID,
			c.Request.Header.Get("XID"),
		)).Begin(nil)

		if err := (&dao.BDao{}).Create(tx); err != nil {
			tx.Rollback()
			c.String(500, "%s", err.Error())
			return
		}

		tx.Commit()
		c.String(200, "%s", "succ")
	})

	if err := route.Run(":7002"); err != nil {
		dao.DisconnectMysql()
		panic(err)
	}
	dao.DisconnectMysql()
}
