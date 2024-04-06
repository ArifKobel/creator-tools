package handlers

import (
	auth_service "github.com/ArifKobel/creator-tools/services/auth"
	"github.com/ArifKobel/creator-tools/services/database"
	"github.com/ArifKobel/creator-tools/services/database/schemas"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type OTPRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func SendOTP() fiber.Handler {
	return func(c fiber.Ctx) error {
		var body OTPRequest
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid request",
			})
		}
		db, err := database.Connect()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
		var user schemas.User
		db.Where("email = ?", body.Email).First(&user)
		if user.ID == 0 {
			db.Create(&schemas.User{
				Email: body.Email,
			})
		}
		otp, err := auth_service.GenerateOTP(6)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
		err = auth_service.SaveOtp(body.Email, otp)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
		err = auth_service.SendOtp(body.Email, otp)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
		return c.JSON(fiber.Map{
			"message": "success",
		})
	}
}

func VerifyOTP() fiber.Handler {
	return func(c fiber.Ctx) error {
		var body OTPRequest
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid request",
			})
		}
		db, err := database.Connect()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
		var user schemas.User
		db.Where("email = ?", body.Email).First(&user)
		if user.ID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "user not found",
			})
		}
		if user.Otp != body.Code {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid otp",
			})
		}
		db.Model(&user).Update("otp", "")
		token, err := auth_service.GenerateJwt(jwt.MapClaims{
			"email": user.Email,
			"id":    user.ID,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
		return c.JSON(fiber.Map{
			"message": "success",
			"token":   token,
		})
	}
}

func VerifyToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unauthorized",
			})
		}
		token = token[7:]
		claims, err := auth_service.GetDataFromToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unauthorized",
			})
		}
		return c.JSON(fiber.Map{
			"message": "success",
			"claims":  claims,
		})
	}
}
