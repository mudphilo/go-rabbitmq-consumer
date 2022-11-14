package consumers

import (
	"database/sql"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

type Consumer struct {
	DB                *sql.DB
	RabbitMQ		*amqp.Connection
}

//SetupConsumers setup consumers
func (c *Consumer) SetupConsumers()  {

	// lets get rabbitMQ queue to listen

	queues := os.Getenv("JOB_QUEUES")

	queueList := strings.Split(queues,",")

	// how many consumers per queue
	consumersPerQueue := 1

	// let loop through each queue and setup a consumer in separate thread
	for _, queueName := range queueList {

		exchangeType := "direct"
		bindingKey := queueName
		exchange := queueName
		prefetchCount := 200

		x := 0

		for x < consumersPerQueue {

			// run consumer as a separate go routime
			go c.RunConsumer(queueName, bindingKey, exchange, exchangeType, prefetchCount)
			x++
		}

	}

	// we have to keep this function running forever to enable consumers continously listen
	select {}

}

func (c *Consumer) RunConsumer(queueName,bindingKey, exchange, exchangeType string, prefetchCount int) error {

	var err error

	// initialise channel
	channel, err := c.RabbitMQ.Channel()
	if err != nil {

		log.Printf("error starting rabbitMQ channel %s",err.Error())
		return fmt.Errorf("channel: %s", err)
	}

	// set channet properties
	err = channel.Qos(
		prefetchCount, // prefetch count
		0,        // prefetch size
		false,    // global
	)
	if err != nil {

		log.Printf("error setting up rabbitMQ channel properties %s",err.Error())
		return err
	}

	// declare exchange
	if err = channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {

		log.Printf("error exchange declare %s",err.Error())
		return err
	}

	// declare queue
	_, err = channel.QueueDeclare(
		queueName, // name of the queue
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err != nil {
		log.Printf("error queue declare %s",err.Error())
		return err
	}

	// bind to queue
	if err = channel.QueueBind(
		queueName,    // name of the queue
		bindingKey,      // bindingKey
		exchange, // sourceExchange
		false,    // noWait
		nil,      // arguments
	); err != nil {

		log.Printf("error queue bind %s",err.Error())
		return err
	}

	// start consuming from queue
	deliveries, err := channel.Consume(
		queueName, // name
		queueName, // consumerTag,
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)

	if err != nil {

		log.Printf("error starting queue consumer %s",err.Error())
		return err
	}

	// use channels to make run endlessly
	forever := make(chan bool)

	// we route incoming messages to different functions based on the queue names

	switch queueName {

		case "FILES_QUEUE":
			c.FilesProcessor(deliveries)
			break

		case "SMS_QUEUE":
			c.SMSProcessor(deliveries)
			break

	}

	// runs forerver
	<-forever

	return nil
}
