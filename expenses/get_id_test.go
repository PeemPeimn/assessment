package expenses

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetOneExpense(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	newsMockRows := sqlmock.
		NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "smoothie", 99, "unit_test", `["food", "beverage"]`)

	mock.ExpectQuery("SELECT * FROM expenses WHERE id=?").
		WithArgs(1).
		WillReturnRows(newsMockRows)

	expected := `{"id":1,"title":"smoothie","amount":79,"note":"abcd","tags":["food","beverage"]}`

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
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
}
