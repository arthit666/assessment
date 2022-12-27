package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	e.GET("/expenses", expenses.GetAllExpanses)

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	fmt.Println("sever start at port:", os.Getenv("PORT"))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	fmt.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("bye bye")
}
