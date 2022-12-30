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

func middlewareVerifyAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token != "November 10, 2009" {
			return c.JSON(http.StatusUnauthorized, "Token Unauthorized!!!")
		}
		return next(c)
	}
}

func main() {
	h := expenses.InitDb()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewareVerifyAuth)

	e.POST("/expenses", h.CreateExpenses)
	e.GET("/expenses/:id", h.GetOneExpenses)
	e.PUT("/expenses/:id", h.PutExpenses)
	e.GET("/expenses", h.GetAllExpanses)

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
