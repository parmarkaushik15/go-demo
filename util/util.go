package util

import (
	"encoding/csv"
	"log" 
	"os"
	"go-demo-api/model"
	"go-demo-api/queue"
)

func ReadCSV(filePath string) (error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening CSV file: %v", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file: %v", err)
		return err
	}

	var users []model.User
	for _, record := range records[1:] { // Skip header row
		user := model.User{
			ID:           record[0],
			FirstName:    record[1],
			LastName:     record[2],
			EmailAddress: record[3],
			CreatedAt:    record[4],
		}
		if record[5] != "" {
			user.DeletedAt = &record[5]
		}
		if record[6] != "" {
			user.MergedAt = &record[6]
		}
		if record[7] != "" {
			user.ParentUserID = &record[7]
		}
		users = append(users, user)
	}

	queue.PublishToRabbitMQ(users)
	return nil
}