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

	e.POST("/expenses", expenses.CreateExpensesHandler)
	e.GET("/expenses/:id", expenses.GetOneExpenses)
	e.PUT("/expenses/:id", expenses.PutExpenses)

	log.Println("sever start at port:", os.Getenv("PORT"))
	log.Fatal(e.Start(os.Getenv("PORT")))

}
