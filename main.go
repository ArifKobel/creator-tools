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
	db.AutoMigrate(&schemas.User{}, &schemas.Video{}, &schemas.Export{})
	app.Use(cors.New())
	authRoutes := app.Group("/auth")
	authRoutes.Post("/send-otp", handlers.SendOTP())
	authRoutes.Post("/verify-otp", handlers.VerifyOTP())
	videoRoutes := app.Group("/video")
	videoRoutes.Post("/create-video", handlers.CreateVideo())
	videoRoutes.Get("/get-videos", handlers.GetVideos())
	videoRoutes.Get("/get-video/:id", handlers.GetVideo())
	videoRoutes.Get("/get-video-file/:id", handlers.GetVideoFile())
	videoRoutes.Get("/get-video-thumbnail/:id", handlers.GetVideoThumbnail())
	videoRoutes.Post("/add-export-url/:id", handlers.AddExportURL())
	videoRoutes.Delete("/delete-video/:id", handlers.DeleteVideo())
	videoRoutes.Get("/get-all-files", func(c fiber.Ctx) error {
		files, err := os.ReadDir("./")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		var filenames []string
		for _, file := range files {
			filenames = append(filenames, file.Name())
		}
		return c.JSON(filenames)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
