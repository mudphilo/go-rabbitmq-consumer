package main

import (
	"fmt"
	"github.com/mudphilo/go-rabbitmq/consumers"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {

	conn := GetRabbitMQConnection()
	consumer := consumers.Consumer{
		DB:       nil,
		RabbitMQ: conn,
	}

	// this is a blocking function, it will not exit/return
	// if run with other APIs run it as goroutine
	consumer.SetupConsumers()

}

func GetRabbitMQConnection() *amqp.Connection {

	host := os.Getenv("RABBITMQ_HOST")
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASSWORD")
	port := os.Getenv("RABBITMQ_PORT")
	vhost := os.Getenv("RABBITMQ_VHOST")

	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port,vhost)

	log.Printf("got rabbitMQ connString %s ",uri)
	conn, err := amqp.Dial(uri)

	if err != nil {

		log.Printf("got error connecting to rabbitMQ %s with %s", err.Error(), uri)
		return nil
	}

	return conn
}
