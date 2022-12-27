package expenses

import (
	"database/sql"
	"net/http"

	"github.com/lib/pq"

	"github.com/labstack/echo/v4"
)

func GetOneExpenses(c echo.Context) error {
	id := c.Param("id")

	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Massage: "can't prepare query expenses statment:" + err.Error()})
	}

	row := stmt.QueryRow(id)

	e := Expenses{}

	err = row.Scan(&e.Id, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Massage: "expenses not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Massage: "can't scan expenses:" + err.Error()})
	}
}

func GetAllExpanses(c echo.Context) error {

	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Massage: "can't prepare query all expenses statment:" + err.Error()})
	}

	rows, err := stmt.Query()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Massage: "can't query all expenses:" + err.Error()})
	}

	es := []Expenses{}

	for rows.Next() {
		e := Expenses{}
		err := rows.Scan(&e.Id, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Massage: "can't scan expenses:" + err.Error()})
		}
		es = append(es, e)
	}
	return c.JSON(http.StatusOK, es)
}
