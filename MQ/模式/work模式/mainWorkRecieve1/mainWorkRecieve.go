package main

import "github.com/all/Mq/Rabbitmq"

/*
* 消费者
 */
//
func main() {
	//链接mq,返回mq实例
	rabbitmq := Rabbitmq.NewRabbitMQ("testSimple", "", "")
	rabbitmq.ConsumeSimple()
}
