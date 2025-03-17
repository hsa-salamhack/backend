package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	godotenv.Load()
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") +
		"@tcp(localhost:3306)/messages_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to the database:", err)
	}

	DB = db
	log.Println("✅ Connected to MariaDB")

	db.AutoMigrate(&Message{})
}

type Message struct {
	ID      uint   `gorm:"primaryKey"`
	ChatID  string `gorm:"index"      json:"chat_id"`
	Content string `                  json:"content"`
	Role    string `                  json:"role"`
}
