package main

import (
	"log"
	"go-demo-api/database"
	"go-demo-api/router"
	"go-demo-api/cache"
	"go-demo-api/queue"
	"go-demo-api/model"
	"go-demo-api/util"
	"go-demo-api/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
)

func main() {
	database.Connect()
	cache.ConnectRedis()
	queue.ConnectRabbitMQ() 
	filePath := config.Config("CSV_PATH") 

	var count int64 
	if err := database.DB.Db.Unscoped().Model(&model.User{}).Count(&count).Error; err != nil {
		log.Printf("Error counting user records: %v", err)
		return
	}

	log.Printf("Total users in database (including soft-deleted): %d", count)
	if count == 0 {
		err := util.ReadCSV(filePath)
		if err != nil {
			log.Printf("Error reading CSV: %v\n", err) 
		}
	}
	queue.ConsumeFromRabbitMQ()

	app := fiber.New()
	app.Use(logger.New())

	app.Use(cors.New())

	router.SetupRoutes(app)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	app.Listen(":8080")
}