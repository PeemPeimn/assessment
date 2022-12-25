package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PeemPeimn/assessment/expenses"
	"github.com/labstack/echo/v4"
)

// This function is a middleware function handling Token Authorization.
// According to Postman's tests the validation value is "November 10, 2009".
func authorizationHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "November 10, 2009" {
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, "invalid key value.")
	}
}

func main() {

	db := expenses.InitDB(os.Getenv("DATABASE_URL"))

	handler := expenses.Handler{DB: db}

	echoInstance := echo.New()

	// Use the customized handler.
	echoInstance.Use(authorizationHandler)

	echoInstance.POST("/expenses", handler.CreateExpense)
	echoInstance.GET("/expenses/:id", handler.GetExpenseByID)

	// Start server
	go func() {
		if err := echoInstance.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			echoInstance.Logger.Fatal("shutting down the server.")
		}
	}()

	// Gracefully shut down after interrupt signal is triggered with a 10-second timeout.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := echoInstance.Shutdown(ctx); err != nil {
		echoInstance.Logger.Fatal(err)
	}
	log.Println("shut down gracefully.")

}
