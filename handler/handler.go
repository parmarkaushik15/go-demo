package handler

import (  
	"github.com/gofiber/fiber/v2"
	"go-demo-api/cache"
	"log"
	"strconv"
	"fmt"
	"go-demo-api/model"
	"go-demo-api/database"
)
 

// Get All Users from db
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB.Db
	sizeParam := c.Query("size", "10") // Default page size is 10
	pageParam := c.Query("page", "1")  // Default page number is 1
	filterField := c.Query("filter_field", "")  // Field to filter on
	filterValue := c.Query("filter_value", "")  // Value for the filter

	size, err := strconv.Atoi(sizeParam)
	if err != nil || size <= 0 {
		return c.Status(400).JSON(fiber.Map{"status": "failure", "message": "Invalid size parameter", "data": ""})
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		return c.Status(400).JSON(fiber.Map{"status": "failure", "message": "Invalid page parameter", "data": ""})
	}

	offset := (page - 1) * size
 
	if filterField != "" && filterValue != "" {
		// Check if the filterField is valid (e.g., it matches a known field in the User model)
		validFields := []string{"name", "email", "status"} // Example valid fields
		isValidField := false
		for _, field := range validFields {
			if field == filterField {
				isValidField = true
				break
			}
		}

		if !isValidField {
			return c.Status(400).JSON(fiber.Map{"status": "failure", "message": "Invalid filter field", "data": ""})
		}

		// If filter is valid, apply the filter to the query
		db = db.Where(fmt.Sprintf("%s = ?", filterField), filterValue)
	}

	// Fetch users from the cache
	cachedUsers, err := cache.GetCachedUsersWithPagination(page, size, filterField, filterValue)
	if err != nil {
		log.Printf("Redis fetch error: %v", err)
		return c.Status(500).JSON(fiber.Map{"status": "failure", "message": "internal server error", "data": ""})
	}
	log.Printf("Total users: %v", len(cachedUsers))
	if len(cachedUsers) > 0 {
		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Users Found", "data": cachedUsers})
	} else {
		// Fetch users from the database if not found in cache
		var users []model.User 
		query := db.Unscoped().Offset(offset).Limit(size)

		// Apply filters to the database query if any filter field and value are provided
		if filterField != "" && filterValue != "" {
			switch filterField {
				case "FirstName":
					query = query.Where("first_name = ?", filterValue)
				case "LastName":
					query = query.Where("last_name = ?", filterValue)
				default:
					// Other filter fields can be added here if necessary
			}
		}

		// Fetch filtered users from the database
		if err := query.Find(&users).Error; err != nil {
			log.Printf("PostgreSQL fetch error: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "failure", "message": "internal server error", "data": ""})
		}
	
		for _, user := range users {
			cache.ProcessCache(user)
		}
		log.Printf("Total users: %v", len(users))
		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Users Found", "data": users})
	}
}
