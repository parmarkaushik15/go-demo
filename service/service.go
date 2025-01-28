package service

import ( 
    "log"
	"fmt"
    "go-demo-api/database"
	"go-demo-api/model"
	"go-demo-api/cache"
	"gorm.io/gorm"
)

func ProcessMessage(user model.User) {
	db := database.DB.Db

	log.Printf("Processing user: %v", user)

	// Check if a user with the same ID already exists
	var existingUser model.User
	if err := db.First(&existingUser, "id = ?", user.ID).Error; err == nil {
		// User with the same ID already exists, log and return
		log.Printf("User with ID %s already exists. Skipping insertion.", user.ID)
		return
	} else if err != gorm.ErrRecordNotFound {
		// Log error if it is not a "record not found" error
		log.Printf("Error checking for existing user in PostgreSQL: %v", err)
		return
	}  
	// Encrypt the email address
	user.EmailAddress = fmt.Sprintf("PGP_SYM_ENCRYPT('%s', 'encryption_key')", user.EmailAddress)

	log.Printf("Inserting user to PostgreSQL: %v", user)

	// Insert the user record
	if err := db.Create(&user).Error; err != nil {
		log.Printf("Error saving user to PostgreSQL: %v", err)
		return
	}

	// Process the user in the cache
	cache.ProcessCache(user)
}
