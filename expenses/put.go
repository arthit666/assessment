package expenses

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) PutExpenses(c echo.Context) error {
	id := c.Param("id")

	e := Expenses{}
	err := c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Massage: err.Error()})
	}

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Massage: err.Error()})
	}

	e.Id = idInt

	stmt, err := h.DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1 RETURNING id, title, amount , note, tags")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Massage: "can't prepare query expenses statment:" + err.Error()})
	}

	row := stmt.QueryRow(id, e.Title, e.Amount, e.Note, pq.Array(e.Tags))

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
