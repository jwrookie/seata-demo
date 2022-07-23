## 启动服务

```shell
docker-compose up -d seata-server scripts
```

## 启动聚合服务

```shell
cd at/aggregation
go mod tidy
go run main.go
```

## 启动服务A

```shell
cd at/service_a
go mod tidy
go run main.go
```

## 启动服务B

```shell
cd at/service_b
go mod tidy
go run main.go
```

## 测试事务提交
```shell
curl http://127.0.0.1:7000/testJobCommit
```

#### 可以看到seata_a.a、seata_b.b会各新增一条记录

> 我们在a 与 b 服务的 commit 之后睡眠一会，也可以看到他们各自的undolog（undolog随业务的commit一起提交）
> 如果a 与 b 服务的 commit 都成功，undolog里的数据会立马删除掉

## 测试事务回滚
```shell
curl http://127.0.0.1:7000/testJobRollback
```

#### 可以看到seata_a.a新增了一条记录，然后很快由于事务回滚会被删除掉

> 脏读：a新增了一条记录，c此时读到了这条记录，由于全局事务失败回滚，a新增的记录会被删掉
> 因此产生了脏读

#### seata_b.b则不会新增记录脏读
