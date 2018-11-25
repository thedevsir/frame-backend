package main

import (
	"os"
	"strings"

	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/config/database"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/thedevsir/frame-backend/config/mail"
	_ "github.com/thedevsir/frame-backend/docs"
	"github.com/thedevsir/frame-backend/routes"
	"github.com/thedevsir/frame-backend/services/storage"
	"github.com/thedevsir/frame-backend/services/validation"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	config.Composer(".env")
	db := &database.Composer{
		Locals:   "resource/locals/locals.json",
		Addrs:    strings.Split(config.DBAddress, ","),
		Database: config.DBName,
		Username: config.DBUsername,
		Password: config.DBPassword,
		Source:   config.DBSource,
	}
	db.Shoot()
	storage.Composer()
	mail.Composer()
}

// @title Frame
// @version 1.0.0
// @description A user system API starter.

// @contact.name Amir Irani
// @contact.email freshmanlimited@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3500
// @BasePath /endpoint

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey AdminApiKeyAuth
// @in header
// @name AdminAuthorization

func main() {
	Run := routes.Composer()
	Run.Use(middleware.Logger())
	Run.Use(middleware.Recover())
	Run.Use(middleware.BodyLimit(config.RoutesBodyLimit))
	Run.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(config.CorsAllowOrigins, ","),
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	if os.Getenv("MODE") != "DEV" {
		Run.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}))
	}
	Run.Validator = &validation.DataValidator{ValidatorData: validator.New()}
	Run.Logger.Fatal(Run.Start(config.Port))
}
