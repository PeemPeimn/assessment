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

func TestPutExpense(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	handler := Handler{DB: db}

	newsMockRows := sqlmock.
		NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "smoothie", 99, "unit_test", `{put_test,beverage}`)

	mock.ExpectPrepare("UPDATE expenses (.+) WHERE (.+) RETURNING (.+)").
		ExpectQuery().WithArgs().WillReturnRows(newsMockRows)

	mockJson := []byte(`{
		"title": "smoothie",
		"amount": 99,
		"note": "unit_test",
		"tags": ["put_test", "beverage"]
		}`)

	expected := `{"id":1,"title":"smoothie","amount":99,"note":"unit_test","tags":["put_test","beverage"]}`

	req := httptest.NewRequest(http.MethodPut, "/expenses/1", bytes.NewBuffer(mockJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Act
	handler.PutExpense(c)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))

}

func TestPutExpenseNotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	handler := Handler{DB: db}

	newsMockRows := sqlmock.
		NewRows([]string{"id", "title", "amount", "note", "tags"})

	mock.ExpectPrepare("UPDATE expenses (.+) WHERE (.+) RETURNING (.+)").
		ExpectQuery().WithArgs().WillReturnRows(newsMockRows)

	mockJson := []byte(`{
		"title": "smoothie",
		"amount": 99,
		"note": "unit_test",
		"tags": ["put_test", "beverage"]
		}`)

	req := httptest.NewRequest(http.MethodPut, "/expenses/1", bytes.NewBuffer(mockJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Act
	handler.PutExpense(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

}
