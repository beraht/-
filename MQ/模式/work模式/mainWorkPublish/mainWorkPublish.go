package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/all/Mq/Rabbitmq"
)

/*
* 生产者
 */

func main() {
	//链接mq,返回mq实例
	rabbitmq := Rabbitmq.NewRabbitMQ("testSimple", "", "")
	//生产消息

	for i := 0; i <= 10; i++ {
		rabbitmq.PublishSimple("hello work" + strconv.Itoa(i))
		time.Sleep(time.Second * 1)
		fmt.Println("发送成功", i)
	}

}
