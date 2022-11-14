package consumers

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type FileObject struct {

	Name string `json:"name"`
	Path string `json:"path"`
}

type SMSObject struct {

	SenderID string `json:"sender_id"`
	Msisdn string `json:"msisdn"`
	Message string `json:"message"`
}

func (c *Consumer) FilesProcessor(deliveries <-chan amqp.Delivery)  {

	for d := range deliveries {

		data := FileObject{}

		err := json.Unmarshal([]byte(d.Body), &data)
		if err != nil {

			log.Printf(" got error decoding queue response to struct %s",err.Error())
			d.Nack(false,true) // nack the message, it will be requeued
			continue
		}

		log.Printf("FilesProcessor ==> file name %s | file path %s",data.Name, data.Path)

		d.Ack(false) // acknowledge the message
	}
}

func (c *Consumer) SMSProcessor(deliveries <-chan amqp.Delivery)  {

	for d := range deliveries {

		data := SMSObject{}

		err := json.Unmarshal([]byte(d.Body), &data)
		if err != nil {

			log.Printf(" got error decoding queue response to struct %s",err.Error())
			d.Nack(false,true) // nack the message, it will be requeued
			continue
		}

		log.Printf("SMSProcessor ==> msisdn %s | sender %s | message %s",data.Msisdn, data.SenderID, data.Message)

		d.Ack(false) // acknowledge the message
	}
}