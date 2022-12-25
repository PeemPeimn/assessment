package expenses

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetOneExpenseNotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Expect no rows
	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"})

	mock.ExpectPrepare("SELECT (.+) FROM expenses WHERE id=?").
		ExpectQuery().
		WithArgs("1").
		WillReturnRows(newsMockRows)

	req := httptest.NewRequest(http.MethodGet, "/expenses/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := Handler{DB: db}

	// Act
	handler.GetExpenseByID(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
