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

type View struct {
	AuthRenderBase
}

func get_dbh() (db *sql.DB, err error) {
	var confp string = "./usercfg.conf"
	conf_fh, err := models.Get_config(confp)
	if err != nil {
		log.Printf("read usercfg config file error")
		return nil, errors.New("read usercfg config file error")
	}
	backend_dsn, err := models.Get_backend_dsn(conf_fh)
	if err != nil {
		return nil, err
	}
	backend_dbh, err := dbh(backend_dsn)
	if err != nil {
		return nil, err
	}
	return backend_dbh, nil
}



func (c *View) Get() error {
	id, err := strconv.ParseInt(c.Req().FormValue("id"), 10, 64)
	if err != nil {
		return c.Render("warn.html", renders.T{"warning": err.Error()})
	}

	engine, err := models.GetEngineById(id)
	if err != nil {
		return c.Render("warn.html", renders.T{"warning": err.Error()})
	}

	o := GetOrm(engine)
	if o == nil {
		return fmt.Errorf("get engine %s failed", engine.Name)
		return c.Render("warn.html", renders.T{"warning": fmt.Sprintf("get engine %s failed", engine.Name)})
	}

	tables, err := o.DBMetas()
	//tables, err := o.Fialect.GetTables()
	if err != nil {
		return c.Render("warn.html", renders.T{"warning": err.Error()})
	}

	dbName := getDatabaseName(o)
	log.Printf("get dbName: %s\n", dbName)
	//get models.User
	UserMsg := c.LoginUser()

	//check user db access
	log.Printf("get usermsg: %#v\n", UserMsg)
	_, ok1 := UserMsg.Database["all"]
	_, ok2 := UserMsg.Database[dbName]
	if !(ok1 || ok2) {
		return c.Render("warn.html", renders.T{
			"warning": fmt.Sprintf("You don't have privielges to access database: %s", dbName),
		})
	}

	var records = make([][]*string, 0)
	var columns = make([]*core.Column, 0)
	tb := c.Req().FormValue("tb")
	tb = strings.Replace(tb, `"`, "", -1)
	tb = strings.Replace(tb, `'`, "", -1)
	tb = strings.Replace(tb, "`", "", -1)

	var isTableView = len(tb) > 0

	sql := c.Req().FormValue("sql")
	old_sql := sql
	if sql != "" {
		log.Printf("user: %s, db: %s, origianl sql: %s\n", UserMsg.Name, dbName, sql)
	}
	var table *core.Table
	var pkIdx int
	var isExecute bool
	var affected int64
	var args = make([]interface{}, 0)

	sql2 := SanitiseSQL(sql)
	if sql != "" {
		log.Printf("user: %s, db: %s, after sanitise: %s\n", UserMsg.Name, dbName, sql2)
	}
	if sql != sql2 && sql2 == "" {
		return c.Render("warn.html", renders.T{
			"warning": fmt.Sprintf("'%s' is not permit to execute, contact to administrtor", sql),
		})
	}
	sql = sql2

	tableName := findTableFromSql(sql)
	if sql != "" || tb != "" {
		if sql != "" {
			isExecute = !strings.HasPrefix(strings.ToLower(sql), "select") && !strings.HasPrefix(strings.ToLower(sql), "show")
		} else if tb != "" {
			sql = fmt.Sprintf("show create table `"+tb+"`")
			//args = append(args, []interface{}{limit, start}...)
		} else {
			return errors.New("unknow operation")
		}

		if isExecute {
			if tableName != "" {
				tableSize := getTableSize(o, dbName, tableName)
				if tableSize > 200 { // table is biger than 200MB
					return errors.New(fmt.Sprintf("table %s is too big( %d MB) to execute it, contact to Administrator", tableName, tableSize))
				} else {
					log.Printf("db: %s, table: %s, size: %d MB\n", dbName, tableName, tableSize)
				}
			}
			res, err := o.Exec(sql)
			log.Printf("user: %s, db: %s, execute sql: %s\n", UserMsg.Name, dbName, sql)
			if err != nil {
				return c.Render("warn.html", renders.T{"warning": err.Error()})
			}
			backend_dbh, err := get_dbh()
			if err != nil {
				log.Printf("get backend_dbh error: %v\n", err)
			}
                	
			if (insertLog(backend_dbh, UserMsg.Name, dbName, "dbweb online table", old_sql)) {
				log.Printf("insert to dbweb_history ok\n")
			}
			affected, _ = res.RowsAffected()
		} else {
			log.Printf("query sql: %s\n", sql)
			rows, err := o.DB().Query(sql, args...)
			if err != nil {
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
				datas := make([]*string, len(columns))
				err = rows.ScanSlice(&datas)
				if err != nil {
					return c.Render("warn.html", renders.T{"warning": err.Error()})
				}
				records = append(records, datas)
			}
		}
	}

	engines, err := c.findEngines()
	if err != nil {
		return c.Render("warn.html", renders.T{"warning": err.Error()})
	}

	return c.Render("root.html", renders.T{
		"engines":     engines,
		"tables":      tables,
		"table":       table,
		"records":     records,
		"columns":     columns,
		"id":          id,
		"sql":         sql,
		"tb":          tb,
		"isExecute":   isExecute,
		"isTableView": isTableView,
		"affected":    affected,
		"pkIdx":       pkIdx,
		"curEngine":   engine.Name,
		"IsLogin":     c.IsLogin(),
	})
}

