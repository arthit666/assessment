package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/arthit666/assessment/expenses"
)

func main() {
	expenses.InitDb()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	log.Println("sever start at port:")
	log.Fatal(e.Start(os.Getenv("PORT")))

}
