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

func TestGetAllExpenses(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	newsMockRows := sqlmock.
		NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "smoothie", 79, "unit_test", `{food,beverage}`).
		AddRow(2, "latte", 88, "unit_test", `{coffee,drink}`)

	expected := `[{"id":1,"title":"smoothie","amount":79,"note":"unit_test","tags":["food","beverage"]},{"id":2,"title":"latte","amount":88,"note":"unit_test","tags":["coffee","drink"]}]`

	mock.ExpectPrepare("SELECT (.+) FROM expenses").
		ExpectQuery().
		WillReturnRows(newsMockRows)

	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	handler := Handler{DB: db}

	// Act
	handler.GetAllExpenses(c)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
}

func TestGetAllExpensesEmpty(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	newsMockRows := sqlmock.
		NewRows([]string{"id", "title", "amount", "note", "tags"})

	mock.ExpectPrepare("SELECT (.+) FROM expenses").
		ExpectQuery().
		WillReturnRows(newsMockRows)

	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	handler := Handler{DB: db}

	// Act
	handler.GetAllExpenses(c)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "null", strings.TrimSpace(rec.Body.String()))
}
