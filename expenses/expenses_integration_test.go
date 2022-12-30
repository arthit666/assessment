//go:build integration

package expenses

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExpenses(t *testing.T) {
	body := bytes.NewBufferString(`{
	"title": "beer",
	"amount": 80,
	"note": "leo",
	"tags": ["food", "beverage"]
	}`)
	e := Expenses{}

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.Id)
	assert.Equal(t, "beer", e.Title)
	assert.Equal(t, 80.0, e.Amount)
	assert.Equal(t, "leo", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)
}
func TestGetAllExpenses(t *testing.T) {
	seedExpenses(t)
	var es []Expenses

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&es)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(es), 0)
}

func TestGetExpensesByID(t *testing.T) {
	e := seedExpenses(t)

	latest := Expenses{}
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(e.Id)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.Id, latest.Id)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)
}

func TestUpdateUserByID(t *testing.T) {
	e := seedExpenses(t)

	body := bytes.NewBufferString(`{
		"title": "beer",
		"amount": 80,
		"note": "leo",
		"tags": ["food", "beverage"]
		}`)

	updated := Expenses{}
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(e.Id)), body)
	err := res.Decode(&updated)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, 0, updated.Id)
	assert.Equal(t, "beer", updated.Title)
	assert.Equal(t, 80.0, updated.Amount)
	assert.Equal(t, "leo", updated.Note)
	assert.Equal(t, []string{"food", "beverage"}, updated.Tags)

}

////////////

func seedExpenses(t *testing.T) Expenses {
	e := Expenses{}
	body := bytes.NewBufferString(`{
	"title": "strawberry smoothie",
	"amount": 79,
	"note": "night market promotion discount 10 bath", 
	"tags": ["food", "beverage"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&e)
	if err != nil {
		t.Fatal("can't create uomer:", err)
	}
	return e
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

////////////

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
