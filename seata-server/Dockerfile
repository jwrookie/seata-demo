FROM golang:1.17

ENV TIME_ZONE=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn,direct"
WORKDIR /go/release
COPY . .
RUN cd seata-go && go mod tidy
RUN cd seata-go/cmd/tc && go build -o tc_server
RUN mv seata-go/cmd/tc/tc_server /go/release
RUN chmod +x /go/release/tc_server

ENTRYPOINT ["/go/release/tc_server", "start", "-config", "./config.yml"]
