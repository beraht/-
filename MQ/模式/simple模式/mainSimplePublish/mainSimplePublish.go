package main

import (
	"fmt"

	"github.com/all/Mq/Rabbitmq"
)

/*
* 生产者
 */

func main() {
	//链接mq,返回mq实例
	rabbitmq := Rabbitmq.NewRabbitMQ("testSimple", "", "")
	//生产一个消息
	rabbitmq.PublishSimple("hello simple..........")
	fmt.Println("发送成功")

}
