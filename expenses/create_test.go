package expenses

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {

	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	newsMockRows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO expenses .*").WithArgs().WillReturnRows(newsMockRows)

	e := echo.New()

	mockJson := []byte(`{
		"title": "smoothie",
	  "amount": 79,
	  "note": "abcd",
	  "tags": ["food", "beverage"]
		}`)

	expected := `{"id":1,"title":"smoothie","amount":79,"note":"abcd","tags":["food","beverage"]}`

	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(mockJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := Handler{DB: db}

	// Act
	handler.CreateExpense(c)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
}

func TestCreateExpenseFailCase(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	newsMockRows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO expenses .*").WithArgs().WillReturnRows(newsMockRows)

	e := echo.New()

	mockJson := []byte(`{
		"title": "smoothie",
	  "amount": "79",
	  "note": "abcd",
	  "tags": ["food", "beverage"]
		}`)

	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(mockJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := Handler{DB: db}

	// Act
	handler.CreateExpense(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
