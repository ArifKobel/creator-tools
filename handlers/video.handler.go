package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ArifKobel/creator-tools/services"
	auth_service "github.com/ArifKobel/creator-tools/services/auth"
	"github.com/ArifKobel/creator-tools/services/database"
	"github.com/ArifKobel/creator-tools/services/database/schemas"
	"github.com/gofiber/fiber/v3"
)

func SplitFilenameAndExtension(filename string) (string, string) {
	splitted := strings.Split(filename, ".")
	return strings.Join(splitted[:len(splitted)-1], "."), splitted[len(splitted)-1]
}

type Request struct {
	Language string `query:"language"`
}

func CreateVideo() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		data, err := auth_service.GetDataFromToken(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		file, err := c.FormFile("video")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Video is required",
			})
		}
		userID := data["id"].(float64)
		db, err := database.Connect()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		video := schemas.Video{
			UserID:    uint(userID),
			Filename:  file.Filename,
			CreatedAt: services.GetCurrentTime(),
		}
		result := db.Create(&video)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		leadingText, err := services.RandStringRunes(10)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		filename, extension := SplitFilenameAndExtension(file.Filename)
		inputFilePath := fmt.Sprintf("uploads/%d/%s-%s.%s", int(userID), leadingText, base64.StdEncoding.EncodeToString([]byte(filename)), extension)
		outputFilePath := fmt.Sprintf("uploads/%d/tmp/%s-%s.wav", int(userID), leadingText, base64.StdEncoding.EncodeToString([]byte(filename)))
		thumbnailPath := fmt.Sprintf("uploads/%d/%s-%s.jpg", int(userID), leadingText, base64.StdEncoding.EncodeToString([]byte(filename)))
		os.MkdirAll(fmt.Sprintf("uploads/%d", int(userID)), os.ModePerm)
		c.SaveFile(file, inputFilePath)
		os.MkdirAll(fmt.Sprintf("uploads/%d/tmp", int(userID)), os.ModePerm)
		services.ConvertToWAV(inputFilePath, outputFilePath)
		duration, err := services.Mp4Duration(outputFilePath)
		os.Remove(outputFilePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		err = services.GenerateThumbnail(inputFilePath, thumbnailPath)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		db.Model(&video).Updates(map[string]interface{}{
			"duration":       duration,
			"filepath":       inputFilePath,
			"thumbnail_path": thumbnailPath,
		})
		language := c.Query("language")
		subtitles, err := services.GenerateSubTitles(fmt.Sprintf("uploads/%d/tmp/%s-%s.wav", int(userID), leadingText, file.Filename), language)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		jsonSubtitles, err := json.Marshal(subtitles)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		db.Model(&video).Update("subtitles", jsonSubtitles)
		return c.JSON(fiber.Map{
			"message":   "File uploaded successfully",
			"subtitles": subtitles,
			"id":        video.ID,
		})
	}
}

func GetVideos() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		page := c.Query("page")
		if page == "" {
			page = "1"
		}
		itemsPerPage := c.Query("itemsPerPage")
		if itemsPerPage == "" {
			itemsPerPage = "5"
		}
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		data, err := auth_service.GetDataFromToken(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		userID := data["id"].(float64)
		db, err := database.Connect()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		var videos []schemas.Video
		pageInt, _ := strconv.Atoi(page)                 // Convert page from string to int
		itemsPerPageInt, _ := strconv.Atoi(itemsPerPage) // Convert itemsPerPage from string to int
		db.Where("user_id = ?", userID).Offset((pageInt - 1) * itemsPerPageInt).Limit(itemsPerPageInt).Find(&videos)
		return c.JSON(fiber.Map{
			"videos": videos,
		})
	}
}

func GetVideo() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		data, err := auth_service.GetDataFromToken(authHeader)
		if data == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		videoID := c.Params("id")
		db, err := database.Connect()
		var video schemas.Video
		err = db.Model(&schemas.Video{}).Preload("Exports").Find(&video, "id = ?", videoID).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		return c.JSON(fiber.Map{
			"video": video,
		})
	}
}

func DeleteVideo() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		data, err := auth_service.GetDataFromToken(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		userID := data["id"].(float64)
		videoID := c.Params("id")
		db, err := database.Connect()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		var video schemas.Video
		db.Where("user_id = ? AND id = ?", userID, videoID).First(&video)
		os.Remove(video.Filepath)
		os.Remove(video.ThumbnailPath)
		db.Delete(&video)
		return c.JSON(fiber.Map{
			"message": "Video deleted successfully",
		})
	}
}

func GetVideoFile() fiber.Handler {
	return func(c fiber.Ctx) error {
		videoID := c.Params("id")
		db, err := database.Connect()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		videoID = strings.Replace(videoID, ".mp4", "", -1)
		var video schemas.Video
		err = db.Where("id = ?", videoID).First(&video).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		if video.Filepath == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "File not found",
			})
		}
		return c.SendFile(video.Filepath)
	}
}

func GetVideoThumbnail() fiber.Handler {
	return func(c fiber.Ctx) error {
		videoID := c.Params("id")
		db, err := database.Connect()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		var video schemas.Video
		db.Where("id = ?", videoID).First(&video)
		fmt.Println(video.ThumbnailPath)
		return c.SendFile(video.ThumbnailPath)
	}
}

func AddExportURL() fiber.Handler {
	return func(c fiber.Ctx) error {
		videoID := c.Params("id")
		db, err := database.Connect()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		var video schemas.Video
		db.Where("id = ?", videoID).First(&video)
		url := c.Body()
		export := schemas.Export{
			VideoID: video.ID,
			URL:     string(url),
		}
		db.Create(&export)
		return c.JSON(fiber.Map{
			"message": "URL added successfully",
		})
	}
}
