package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srv-cashpay/report/routes"
)

func main() {
	// Inisialisasi Echo
	e := routes.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(CORSMiddleware())

	// Sertifikat Let's Encrypt
	certFile := "/certs/fullchain.pem"
	keyFile := "/certs/privkey.pem"

	// Jalankan HTTPS
	if err := e.StartTLS(":2345", certFile, keyFile); err != nil {
		log.Fatal("StartTLS error: ", err)
	}
}

// CORSMiddleware custom
func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")

			if c.Request().Method == "OPTIONS" {
				return c.NoContent(204)
			}
			return next(c)
		}
	}
}
