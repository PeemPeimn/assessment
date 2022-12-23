package main

import (
	"fmt"
	"os"
	"github.com/PeemPeimn/assessment/expenses"
)

func main() {
	expenses.InitDB()
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
}
