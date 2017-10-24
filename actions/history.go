package actions

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-xorm/core"
	"github.com/go-xorm/dbweb/models"
	"github.com/tango-contrib/renders"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type History struct {
	AuthRenderBase
}

func get_dbh_backend() (db *sql.DB, err error) {
        var confp string = "./usercfg.conf"
        conf_fh, err := models.Get_config(confp)
        if err != nil {
                log.Printf("read usercfg config file error")
                return nil, errors.New("read usercfg config file error")
        }
        backend_dsn, err := models.Get_backend_dsn(conf_fh)
        if err != nil {
		log.Printf("get backend_dsn error: %v\n", err)
                return nil, err
        }
        backend_dbh, err := dbh(backend_dsn)
        if err != nil {
		log.Printf("get backend_dbh error: %v\n", err)
                return nil, err
        }
        return backend_dbh, nil
}

func (c *History) Get() error {
	//get models.User
	UserMsg := c.LoginUser()

	//check user db access
	log.Printf("get usermsg: %#v\n", UserMsg)

	var records = make([][]*string, 0)
	var columns = make([]*core.Column, 0)
	tb := "dbweb_history"

	//var table *core.Table
	var isExecute bool
	var isTableView = len(tb) > 0
	var affected int64
	var total int
	var countSql string
	var sql string
	//var args = make([]interface{}, 0)

	start, _ := strconv.Atoi(c.Req().FormValue("start"))
	limit, _ := strconv.Atoi(c.Req().FormValue("limit"))

	if limit == 0 {
		limit = 20
	}

	countSql = "select count(*) as total from `" + tb + "`"
	sql = fmt.Sprintf("select user, db, changes, create_time from `"+tb+"` where user = \"%s\" order by create_time desc LIMIT %d OFFSET %d", UserMsg.Name, limit, start)
	if sql != "" || tb != "" {
		isExecute = !strings.HasPrefix(strings.ToLower(sql), "select")

		if isExecute {
			backend_dbh, err := get_dbh_backend()
			if err != nil {
				log.Printf("get backend_dbh error: %v\n", err)
				return c.Render("warn.html", renders.T{"warning": err.Error()})
			}
                	
			res, err := ExecQuery(backend_dbh, sql)
			if err != nil {
				return c.Render("warn.html", renders.T{"warning": err.Error()})
			}

			affected, _ = res.RowsAffected()
		} else {
			backend_dbh, err := get_dbh_backend()
			log.Printf("2: sql: %s\ncountSql: %s\n", sql, countSql)
			if err != nil {
				log.Printf("get backend_dbh error: %v\n", err)
				return c.Render("warn.html", renders.T{"warning": err.Error()})
			}

			if len(countSql) > 0 {
				err = QueryRow(backend_dbh, countSql).Scan(&total)
				if err != nil {
					log.Printf("get queryrow error: %v\n", err)
					return c.Render("warn.html", renders.T{"warning": err.Error()})
				}
			}

			rows, err := Query(backend_dbh, sql)
			if err != nil {
				log.Printf("get rows error: %v\n", err)
				return c.Render("warn.html", renders.T{"warning": err.Error()})
			}
			defer rows.Close()

			cols, err := rows.Columns()
			if err != nil {
				return c.Render("warn.html", renders.T{"warning": err.Error()})
			}

			for _, col := range cols {
				columns = append(columns, &core.Column{
					Name: col,
				})
			}

			for rows.Next() {
				//datas := make([]*string, len(columns))
				var history_user string
				var history_db string
				var history_changes string
				var history_cretime string
				
				err = rows.Scan(&history_user, &history_db, &history_changes, &history_cretime)
				if err != nil {
					log.Printf("get datas err: %v\n", err)
					return c.Render("warn.html", renders.T{"warning": err.Error()})
				}
				records = append(records, []*string{&history_user, &history_db, &history_changes, &history_cretime})
			}
		}
	}

	//engines, err := c.findEngines()
	//if err != nil {
	//	return c.Render("warn.html", renders.T{"warning": err.Error()})
	//}

	return c.Render("history.html", renders.T{
		"records":     records,
		"columns":     columns,
		"tb":          tb,
		"isExecute":   isExecute,
		"isTableView": isTableView,
		"limit":       limit,
		"curPage":     start / limit,
		"totalPage":   pager(total, limit),
		"affected":    affected,
		"IsLogin":     c.IsLogin(),
	})
}

func pager(total, limit int) int {
	if total%limit == 0 {
		return total / limit
	}
	return total / limit + 1
}

func curPage(start, limit int) int {
	return start / limit
}
