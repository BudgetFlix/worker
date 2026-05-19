package rabbitmq

import(
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Ack(msg amqp.Delivery){

	err := msg.Ack(false)

	if err != nil {
		log.Printf("faild to ack message: %v", err)
	}
}

func Nack (msg amqp.Delivery) {
	err := msg.Nack(false,false)

	if err != nil {
		log.Printf("faild to nack messsage: %v",err)
	}

}