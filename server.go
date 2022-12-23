package main

import (
	"log"
	"net/http"
	"os"

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
		return c.JSON(http.StatusUnauthorized, "Invalid key value.")
	}
}

func main() {
	expenses.InitDB()

	echoInstance := echo.New()

	// Use the customized handler.
	echoInstance.Use(authorizationHandler)

	echoInstance.POST("/expenses", expenses.CreateExpensesHandler)

	log.Printf("Server started at :%s\n", os.Getenv("PORT"))
	log.Fatal(echoInstance.Start(os.Getenv("PORT")))
	log.Println("bye bye!")

}
