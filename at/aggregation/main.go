package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opentrx/seata-golang/v2/pkg/client/tm"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentrx/seata-golang/v2/pkg/client"
	context2 "github.com/opentrx/seata-golang/v2/pkg/client/base/context"
	"github.com/opentrx/seata-golang/v2/pkg/client/base/model"
	"github.com/opentrx/seata-golang/v2/pkg/client/config"
	"github.com/opentrx/seata-golang/v2/pkg/util/log"
)

var proxySvc = &ProxyService{
	Svc: &Svc{},
}

func main() {
	conf := config.InitConfiguration("./conf/client.yml")

	log.Init(conf.Log.LogPath, conf.Log.LogLevel)
	client.Init(conf)

	tm.Implement(proxySvc)

	route := gin.New()

	// 测试提交
	route.GET("/testJobCommit", func(c *gin.Context) {

		if err := proxySvc.TestJob(c, false); err != nil {
			c.String(500, "%s", err.Error())
			return
		}
		c.String(200, "%s", "succ")
	})

	// 测试回滚
	route.GET("/testJobRollback", func(c *gin.Context) {

		if err := proxySvc.TestJob(c, true); err != nil {
			c.String(500, "%s", err.Error())
			return
		}
		c.String(200, "%s", "succ")
	})

	if err := route.Run(":17000"); err != nil {
		panic(err)
	}
}

type Svc struct {
}

func (svc *Svc) TestJob(ctx context.Context, rollback bool) error {
	var err error

	rootContext, ok := ctx.(*context2.RootContext)
	if !ok {
		return errors.New("create rootContext fail")
	}

	// 调用服务A
	if err = svc.reqA(rootContext); err != nil {
		return err
	}
	// 调用服务B
	if err = svc.reqB(rootContext, rollback); err != nil {
		return err
	}

	// commit
	return nil
}

func (svc *Svc) reqA(rootContext *context2.RootContext) error {
	reqA, err := http.NewRequest("POST", "http://localhost:7001/jobA", nil)
	if err != nil {
		return err
	}
	reqA.Header.Set("Content-Type", "application/json")
	reqA.Header.Set("XID", rootContext.GetXID())

	client := &http.Client{}
	resultA, err := client.Do(reqA)
	if err != nil {
		return err
	}
	if resultA.StatusCode != 200 {
		return fmt.Errorf("A err %d", resultA.StatusCode)
	}

	return nil
}
func (svc *Svc) reqB(rootContext *context2.RootContext, rollback bool) error {
	data := struct {
		Rollback bool `json:"rollback"`
	}{
		Rollback: rollback,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	reqB, err := http.NewRequest("POST", "http://localhost:7002/jobB", bytes.NewBuffer(dataJson))
	if err != nil {
		return err
	}
	reqB.Header.Set("Content-Type", "application/json")
	reqB.Header.Set("XID", rootContext.GetXID())

	client := &http.Client{}
	resultB, err := client.Do(reqB)
	if err != nil {
		return err
	}
	if resultB.StatusCode != 200 {
		return fmt.Errorf("B err %d", resultB.StatusCode)
	}

	return nil
}

var methodTransactionInfo = map[string]*model.TransactionInfo{
	"TestJob": {
		TimeOut:     60000000,
		Name:        "TestJob",
		Propagation: model.Required,
	},
}

type ProxyService struct {
	*Svc
	TestJob func(ctx context.Context, rollback bool) error
}

func (svc *ProxyService) GetProxyService() interface{} {
	return svc.Svc
}
func (svc *ProxyService) GetMethodTransactionInfo(methodName string) *model.TransactionInfo {
	return methodTransactionInfo[methodName]
}
