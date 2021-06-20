package main

import (
	"crud/pkg/handler"
	"crud/pkg/repository"
	"crud/pkg/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

// @title nix Education API
// @version 0.0.1
// @description Homework for nix educations golang course.

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf(err.Error())
	}

	db, err := repository.NewMysqlDB(&repository.Config{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	r := repository.NewRepository(db)
	s := service.NewService(r)
	h := handler.NewHandler(s)

	e := echo.New()

	h.InitRoutes(e)

	//e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":8080"))
}
