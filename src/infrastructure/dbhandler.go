package infrastructure


import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"interfaces"
)

type DBHandler struct {
	Conn *sql.DB
}

func (d *DBHandler) Execute(statement string) {
	d.Conn.Exec(statement)
}

func (d *DBHandler) Query(statement string) interfaces.Row {
	rows, err := d.Conn.Query(statement)
	if err != nil {
		fmt.Println(err)
		return new(DBRow)
	}
	row := new(DBRow)
	row.Rows = rows
	return row
}

type DBRow struct {
	Rows *sql.Rows
}

func (r DBRow) Scan(dest ...interface{}) {
	r.Rows.Scan(dest...)
}

func (r DBRow) Next() bool {
	return r.Rows.Next()
}

func NewDBHandler(dsn string) *DBHandler {
	conn, _ := sql.Open("mysql", dsn)
	dbHandler := new(DBHandler)
	dbHandler.Conn = conn
	return dbHandler
}
