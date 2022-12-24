package expenses

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type (

	// Handler contains a DB connection
	// and has handling method for requests.
	Handler struct {
		DB *sql.DB
	}

	// Expense is a struct used to represent an expense JSON response.
	// Expense consists of fields below according to the database's table.
	Expense struct {
		ID     int      `json:"id"`
		Title  string   `json:"title"`
		Amount int      `json:"amount"`
		Note   string   `json:"note"`
		Tags   []string `json:"tags"`
	}

	// ErrorResponse represents an error in a JSON response.
	ErrorResponse struct {
		Message string `json:"message"`
	}
)

// CreateExpensesHandler handles HTTP POST request to create a new expense.
// This function receives echo.Context as a parameter
// and returns a JSON response with status code.
func (handler Handler) CreateExpense(c echo.Context) error {

	var expense Expense

	err := c.Bind(&expense)

	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorResponse{Message: "cannot unmarshal request's body. " + err.Error()})
	}

	row := handler.DB.QueryRow(`
		INSERT INTO expenses (title, amount, note, tags) 
		values ($1, $2, $3, $4) 
		RETURNING id
	`, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags))

	err = row.Scan(&expense.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "cannot create the expense. " + err.Error()})
	}

	return c.JSON(http.StatusCreated, expense)
}