package actions

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
)

func dbh(dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return db, err
	}
	return db, nil
}

func Query(db *sql.DB, q string) (*sql.Rows, error) {
	return db.Query(q)
}

func QueryRow(db *sql.DB, q string) *sql.Row {
	return db.QueryRow(q)
}

func ExecQuery(db *sql.DB, q string) (sql.Result, error) {
	return db.Exec(q)
}

func insertLog(db *sql.DB, user string, database string, tag string, sql string) bool {
	_, err := ExecQuery(db, fmt.Sprintf("insert into dbweb_history (user, db, tag, changes, create_time) values('%s', '%s', '%s', '%s', now())", user, database, tag, sql_escape(sql)))
	if err != nil {
		log.Printf("insertLog error: %v\n", err)
		return false
	}
	return true
}



func sql_escape(s string) string {
	var j int = 0
	if len(s) == 0 {
		return ""
	}

	tempStr := s[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}
