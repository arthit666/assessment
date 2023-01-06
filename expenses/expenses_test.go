//go:build unit
// +build unit

package expenses

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var jsonPayload = `{
	"title": "test-title",
	"amount": 99,
	"note": "test-note",
	"tags": ["test-tags1","test-tags2"]
}`

func TestGetAllExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "test-title", "99", "test-note", `{"test-tags1","test-tags2"}`)

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").ExpectQuery().WithArgs().WillReturnRows(newsMockRows)

	defer db.Close()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	c := e.NewContext(req, rec)
	expected := "[{\"id\":1,\"title\":\"test-title\",\"amount\":99,\"note\":\"test-note\",\"tags\":[\"test-tags1\",\"test-tags2\"]}]"

	// Act
	err = h.GetAllExpanses(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetOneExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "test-title", "99", "test-note", `{"test-tags1","test-tags2"}`)

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")).ExpectQuery().WithArgs("1").WillReturnRows(newsMockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	expected := "{\"id\":1,\"title\":\"test-title\",\"amount\":99,\"note\":\"test-note\",\"tags\":[\"test-tags1\",\"test-tags2\"]}"

	// Act
	err = h.GetOneExpenses(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestCreateExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id"}).
		AddRow("1")

	db, mock, err := sqlmock.New()

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO expenses (title, amount,note,tags) values ($1, $2, $3, $4)  RETURNING id")).WithArgs("test-title", 99.0, "test-note", pq.Array([]string{"test-tags1", "test-tags2"})).WillReturnRows(newsMockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	h := handler{db}
	c := e.NewContext(req, rec)

	expected := "{\"id\":1,\"title\":\"test-title\",\"amount\":99,\"note\":\"test-note\",\"tags\":[\"test-tags1\",\"test-tags2\"]}"

	// Act
	h.CreateExpenses(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestPutExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "test-title", "99", "test-note", `{"test-tags1","test-tags2"}`)

	db, mock, err := sqlmock.New()

	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1 RETURNING id, title, amount , note, tags")).ExpectQuery().WithArgs("1", "test-title", 99.0, "test-note", pq.Array([]string{"test-tags1", "test-tags2"})).WillReturnRows(newsMockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	expected := "{\"id\":1,\"title\":\"test-title\",\"amount\":99,\"note\":\"test-note\",\"tags\":[\"test-tags1\",\"test-tags2\"]}"

	// Act
	h.PutExpenses(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
