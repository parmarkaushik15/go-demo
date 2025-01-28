package cache

import (
    "encoding/json"
    "context"
    "go-demo-api/model"
    "github.com/go-redis/redis/v8"
    "log"
    "os"
    "fmt" 
)

var RedisClient *redis.Client

func ConnectRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr: os.Getenv("REDIS_URL"),
    })

    _, err := RedisClient.Ping(context.Background()).Result()
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }

    log.Println("Redis Connected")
}


func ProcessCache(user model.User) {
	message, err := json.Marshal(user)
    if err != nil {
        log.Printf("Error marshaling user to JSON: %v", err) 
    } else {
        cacheKey := fmt.Sprintf("user:%s", user.ID)
        err = RedisClient.Set(context.Background(), cacheKey, message, 0).Err()
        if err != nil {
            log.Printf("Redis cache error: %v", err)
        }
        err = RedisClient.RPush(context.Background(), "user", user.ID).Err()
        if err != nil {
            log.Printf("Error adding user ID to Redis list: %v", err)
        }
    }
}

func GetCachedUsersWithPagination(page int, size int, filterField, filterValue string) ([]model.User, error) {
	// Get keys matching the pattern "user:*" from Redis
	keys, err := RedisClient.Keys(context.Background(), "user:*").Result()
	if err != nil {
		log.Printf("Redis fetch error: %v", err)
		return nil, err
	}

	var users []model.User

	// Iterate through all the keys and fetch user data from Redis
	for _, key := range keys {
		userData, err := RedisClient.Get(context.Background(), key).Result()
		if err != nil {
			log.Printf("Error fetching data for key %s: %v", key, err)
			continue
		}

		var user model.User
		// Unmarshal the user data into a User struct
		if err := json.Unmarshal([]byte(userData), &user); err != nil {
			log.Printf("Error unmarshaling data for key %s: %v", key, err)
			continue
		}

		// Apply the filter if applicable
		if filterField != "" && filterValue != "" {
			// Check if the filter matches the user's field and value
			// You may need to add more logic to match your filtering requirements
			switch filterField {
			case "FirstName":
				if user.FirstName != filterValue {
					continue
				}
			case "LastName":
				if user.LastName != filterValue {
					continue
				}
			default:
				// If filter field doesn't match, continue to next user
				continue
			}
		}

		// Add the user to the list if it passes the filter
		users = append(users, user)
	}

	// Calculate the starting and ending indices for pagination after filtering
	start := (page - 1) * size
	end := start + size
	if end > len(users) {
		end = len(users)
	}

	// Slice the users array based on pagination
	paginatedUsers := users[start:end]

	return paginatedUsers, nil
}


func GetCachedUsers() ([]model.User, error) {
    // Get keys matching the pattern "user:*" from Redis
    keys, err := RedisClient.Keys(context.Background(), "user:*").Result()
    if err != nil {
        log.Printf("Redis fetch error: %v", err)
        return nil, err
    }

    var users []model.User
    // Iterate through the keys to fetch and unmarshal the user data
    for _, key := range keys {
        userData, err := RedisClient.Get(context.Background(), key).Result()
        if err != nil {
            log.Printf("Error fetching data for key %s: %v", key, err)
            continue
        }

        var user model.User
        // Unmarshal the user data into a User struct
        if err := json.Unmarshal([]byte(userData), &user); err != nil {
            log.Printf("Error unmarshaling data for key %s: %v", key, err)
            continue
        }
        users = append(users, user)
    }

    return users, nil
}