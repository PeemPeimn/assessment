//go:build integration

package expenses

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const url = "postgresql://root:root@db/it-db?sslmode=disable"

func TestITCreateExpense(t *testing.T) {

	// Arrange
	db := InitDB(url)

	handler := Handler{DB: db}

	e := echo.New()

	mockJson := []byte(`{
		"title": "latte",
	  "amount": 99,
	  "note": "integration",
	  "tags": ["coffee", "beverage"]
		}`)

	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(mockJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expected := Expense{0, "latte", 99, "integration", []string{"coffee", "beverage"}}
	got := Expense{}

	// Act
	handler.CreateExpense(c)

	responseJson := rec.Body.String()

	json.Unmarshal([]byte(responseJson), &got)

	// t.Log(got)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEqual(t, expected.ID, got.ID)
	assert.Equal(t, expected.Title, got.Title)
	assert.Equal(t, expected.Amount, got.Amount)
	assert.Equal(t, expected.Note, got.Note)
	assert.Equal(t, expected.Tags, got.Tags)

}
