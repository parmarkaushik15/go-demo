package queue

import (
    "github.com/streadway/amqp"
    "log"
    "os" 
    "go-demo-api/model"
    "go-demo-api/service"
    "encoding/json"
)

var RabbitMQChannel *amqp.Channel

func ConnectRabbitMQ() {
    conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
    if err != nil {
        log.Fatal("Failed to connect to RabbitMQ:", err)
    }

    RabbitMQChannel, err = conn.Channel()
    if err != nil {
        log.Fatal("Failed to open a channel:", err)
    }

    log.Printf("RabbitMQ connected")
}

func PublishToRabbitMQ(users []model.User) error {
	queue, err := RabbitMQChannel.QueueDeclare(
		"users_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("RabbitMQ queue declaration error: %v", err)
		return err
	}

    for _, user := range users {
		userJSON, err := json.Marshal(user)
		if err != nil {
			log.Printf("Error marshaling user to JSON: %v", err)
			return err
		}

		// Log the JSON message
		log.Printf("Send to queue message: %s", string(userJSON)) 
		err = RabbitMQChannel.Publish(
			"",
			queue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        userJSON,
			},
		)
		if err != nil {
			log.Printf("RabbitMQ publish error: %v", err)
			return err
		}
	}

	return nil
}   

func ConsumeFromRabbitMQ() error{
    msgs, err := RabbitMQChannel.Consume(
		"users_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("RabbitMQ consume error: %v", err)
		return err
	}

	// for msg := range msgs {
	// 	service.ProcessMessage(string(msg.Body))
	// }

    go func() {
		for d := range msgs {
			// Parse the JSON message
			var user model.User
			err := json.Unmarshal(d.Body, &user)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v", err)
				continue
			}

			// Log the user data
			log.Printf("Received User: %+v", user)

            service.ProcessMessage(user)
			// Process the user data as needed
		}
	}()
	return nil
}

func FormatPointer(p *string) string {
	if p == nil {
		return "NULL"
	}
	return *p
}
