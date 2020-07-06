package Rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

/*
	simple 模式
*/
//url 格式: amqp://账户:密码@mq服务器地址:端口号/vhost
const MQUEL = "amqp://root:123@111.229.163.26:5672"

//结构体信息
type RabbitMQ struct {
	//保存的链接
	conn *amqp.Connection
	//管道
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	//链接信息
	Mqurl string
}

//创建mq的实例,返回RabbitMQ结构体
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQUEL,
	}

	var err error
	//创建rabbitmq链接.
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接失败")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取chnnel失败")
	return rabbitmq
}

//断开channel 和 connection,节省资源
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, msg string) {
	if err != nil {
		log.Fatal("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}

//simple 模式
func NewRabbitMQSimple(queueName string) *RabbitMQ {

	rabbitmq := NewRabbitMQ(queueName, "", "")

	return rabbitmq
}

//simple 模式 生产
func (r *RabbitMQ) PublishSimple(msg string) {
	//1.申请队列,如果队列不存在,会自动创建,存在则跳过
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否支持持久化 durable:
		false,
		// 是否为自动删除 autoDelete:
		false,
		// 是否具有排他性 exclusive:
		false,
		// 是否阻塞 noWait:
		false,
		// 额外属性 args:
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true,根据exchange类型和roukey规则,
		// 如果无法找到符合条件的队列,那么会会把发送消息返回给发送者 mandatory :
		false,
		////如果为true,当exchange发送消息队列到队列系统后发现队列上没有绑定消费者,则会把消息发还给发送者
		//immediate:
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
}

//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//接收消息
	msgs, err := r.channel.Consume(
		r.QueueName, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}
	//消费
	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
