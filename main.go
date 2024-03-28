package main

import (
	"fmt"
	"os"

	"github.com/ArifKobel/creator-tools/handlers"
	"github.com/ArifKobel/creator-tools/services/database"
	"github.com/ArifKobel/creator-tools/services/database/schemas"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024,
	})
	db, err := database.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.AutoMigrate(&schemas.User{}, &schemas.Video{})
	app.Use(cors.New())
	authRoutes := app.Group("/auth")
	authRoutes.Post("/send-otp", handlers.SendOTP())
	authRoutes.Post("/verify-otp", handlers.VerifyOTP())
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen("0.0.0.0:" + port)
}
