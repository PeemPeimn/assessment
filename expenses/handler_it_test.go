//go:build integration

package expenses

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
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
	  "note": "integration_create",
	  "tags": ["coffee", "beverage"]
		}`)

	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(mockJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expected := Expense{0, "latte", 99, "integration_create", []string{"coffee", "beverage"}}
	got := Expense{}

	// Act
	handler.CreateExpense(c)

	responseJson := rec.Body.String()

	json.Unmarshal([]byte(responseJson), &got)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEqual(t, expected.ID, got.ID)
	assert.Equal(t, expected.Title, got.Title)
	assert.Equal(t, expected.Amount, got.Amount)
	assert.Equal(t, expected.Note, got.Note)
	assert.Equal(t, expected.Tags, got.Tags)

}

func TestITGetExpenseByID(t *testing.T) {

	// Arrange
	db := InitDB(url)

	handler := Handler{DB: db}

	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockExpense := Expense{0, "latte", 99, "integration_getID", []string{"coffee", "beverage"}}

	row := handler.DB.QueryRow(`
		INSERT INTO expenses (title, amount, note, tags) 
		values ($1, $2, $3, $4) 
		RETURNING id
	`, mockExpense.Title, mockExpense.Amount, mockExpense.Note, pq.Array(mockExpense.Tags))

	err := row.Scan(&mockExpense.ID)
	if err != nil {
		t.Fatal("cannot create mock expense.")
	}

	// Set path param value
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(mockExpense.ID))

	expected := Expense{mockExpense.ID, "latte", 99, "integration_getID", []string{"coffee", "beverage"}}
	got := Expense{}

	// Act
	handler.GetExpenseByID(c)

	responseJson := rec.Body.String()

	json.Unmarshal([]byte(responseJson), &got)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Title, got.Title)
	assert.Equal(t, expected.Amount, got.Amount)
	assert.Equal(t, expected.Note, got.Note)
	assert.Equal(t, expected.Tags, got.Tags)

}

func TestITPutExpense(t *testing.T) {

	// Arrange
	db := InitDB(url)

	handler := Handler{DB: db}

	mockJson := []byte(`{
		"title": "latte",
	  "amount": 99,
	  "note": "integration_put",
	  "tags": ["coffee", "beverage"]
		}`)

	req := httptest.NewRequest(http.MethodPut, "/expenses/1", bytes.NewBuffer(mockJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockExpense := Expense{1, "mocha", 99, "mock_put", []string{"abcd", "efgh"}}

	row := handler.DB.QueryRow(`
		INSERT INTO expenses (title, amount, note, tags) 
		values ($1, $2, $3, $4) 
		RETURNING id
	`, mockExpense.Title, mockExpense.Amount, mockExpense.Note, pq.Array(mockExpense.Tags))

	err := row.Scan(&mockExpense.ID)
	if err != nil {
		t.Fatal("cannot create mock expense.")
	}

	// Set path param value
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(mockExpense.ID))

	expected := Expense{mockExpense.ID, "latte", 99, "integration_put", []string{"coffee", "beverage"}}
	got := Expense{}

	// Act
	handler.PutExpense(c)

	responseJson := rec.Body.String()

	json.Unmarshal([]byte(responseJson), &got)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Title, got.Title)
	assert.Equal(t, expected.Amount, got.Amount)
	assert.Equal(t, expected.Note, got.Note)
	assert.Equal(t, expected.Tags, got.Tags)
}

func TestITGetAllExpenses(t *testing.T) {

	// Arrange
	db := InitDB(url)

	handler := Handler{DB: db}

	_, err := handler.DB.Exec("DELETE FROM expenses")
	if err != nil {
		t.Fatal("cannot clear database for testing. " + err.Error())
	}

	req := httptest.NewRequest(http.MethodPut, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockExpenses := []Expense{
		{0, "mocha", 99, "mock_get", []string{"abcd", "efgh"}},
		{0, "latte", 88, "mock_get", []string{"ijkl", "mnop"}},
		{0, "espresso", 77, "mock_get", []string{"qrst", "uvwx"}},
	}

	for i := range mockExpenses {
		row := handler.DB.QueryRow(`
			INSERT INTO expenses (title, amount, note, tags) 
			values ($1, $2, $3, $4) 
			RETURNING id
		`, mockExpenses[i].Title, mockExpenses[i].Amount, mockExpenses[i].Note, pq.Array(mockExpenses[i].Tags))

		err := row.Scan(&mockExpenses[i].ID)
		if err != nil {
			t.Fatal("cannot create mock expense.")
		}
	}
	// t.Log(mockExpenses)

	gotList := []Expense{}

	// Act
	handler.GetAllExpenses(c)

	responseJson := rec.Body.String()

	json.Unmarshal([]byte(responseJson), &gotList)
	// t.Log(gotList)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	for i, expected := range mockExpenses {
		got := gotList[i]

		assert.Equal(t, expected.ID, got.ID)
		assert.Equal(t, expected.Title, got.Title)
		assert.Equal(t, expected.Amount, got.Amount)
		assert.Equal(t, expected.Note, got.Note)
		assert.Equal(t, expected.Tags, got.Tags)
	}
}
