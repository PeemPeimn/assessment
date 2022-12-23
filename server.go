package main

import (
	"log"
	"os"

	"github.com/PeemPeimn/assessment/expenses"
	"github.com/labstack/echo/v4"
)

func main() {
	expenses.InitDB()

	echoInstance := echo.New()

	echoInstance.POST("/expenses", expenses.CreateExpensesHandler)

	log.Printf("Server started at :%s\n", os.Getenv("PORT"))
	log.Fatal(echoInstance.Start(os.Getenv("PORT")))
	log.Println("bye bye!")

}
