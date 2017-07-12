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
)

type View struct {
	AuthRenderBase
}

func (c *View) Get() error {
	id, err := strconv.ParseInt(c.Req().FormValue("id"), 10, 64)
	if err != nil {
		return err
	}

	engine, err := models.GetEngineById(id)
	if err != nil {
		return err
	}

	o := GetOrm(engine)
	if o == nil {
		return fmt.Errorf("get engine %s failed", engine.Name)
	}

	tables, err := o.DBMetas()
	if err != nil {
		return err
	}

	dbName := getDatabaseName(o)
	//get models.User
	UserMsg := c.LoginUser();

	//check user db access
	log.Printf("msg: %#v\n", UserMsg)
	_, ok1 := UserMsg.Database["all"]
	_, ok2 := UserMsg.Database[dbName]
	if !(ok1 || ok2) {
		return errors.New(fmt.Sprintf("You don't have privielges to access database: %s", dbName))
	}

	var records = make([][]*string, 0)
	var columns = make([]*core.Column, 0)
	tb := c.Req().FormValue("tb")
	tb = strings.Replace(tb, `"`, "", -1)
	tb = strings.Replace(tb, `'`, "", -1)
	tb = strings.Replace(tb, "`", "", -1)

	var isTableView = len(tb) > 0

	sql := c.Req().FormValue("sql")
	log.Printf("user: %s, db: %s, origianl sql: %s\n", UserMsg.Name, dbName, sql)
	var table *core.Table
	var pkIdx int
	var isExecute bool
	var affected int64
	var total int
	var countSql string
	var args = make([]interface{}, 0)

	start, _ := strconv.Atoi(c.Req().FormValue("start"))
	limit, _ := strconv.Atoi(c.Req().FormValue("limit"))

	sql2 := SanitiseSQL(sql)
	log.Printf("user: %s, db: %s, after sanitise: %s\n", UserMsg.Name, dbName, sql2)
	if sql != sql2 && sql2 == "" {
		return errors.New(fmt.Sprintf("'%s' is not permit to execute, concat to administrtor", sql))
	}
	sql = sql2

	tableName := findTableFromSql(sql)
	if limit == 0 {
		limit = 20
	}
	if sql != "" || tb != "" {
		if sql != "" {
			isExecute = !strings.HasPrefix(strings.ToLower(sql), "select") && !strings.HasPrefix(strings.ToLower(sql), "show")
		} else if tb != "" {
			countSql = "select count(*) from `" + tb + "`"
			sql = fmt.Sprintf("select * from `"+tb+"` LIMIT %d OFFSET %d", limit, start)
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
				return err
			}
			affected, _ = res.RowsAffected()
		} else {
			if len(countSql) > 0 {
				err = o.DB().QueryRow(countSql).Scan(&total)
				if err != nil {
					return err
				}
				fmt.Println("total records:", total)
			}

			rows, err := o.DB().Query(sql, args...)
			if err != nil {
				return err
			}
			defer rows.Close()

			cols, err := rows.Columns()
			if err != nil {
				return err
			}

			if len(tb) > 0 {
				for _, tt := range tables {
					if tb == tt.Name {
						table = tt
						break
					}
				}
				if table != nil {
					for i, col := range cols {
						c := table.GetColumn(col)
						if len(table.PKColumns()) == 1 && c.IsPrimaryKey {
							pkIdx = i
						}
						columns = append(columns, c)
					}
				}
			} else {
				for _, col := range cols {
					columns = append(columns, &core.Column{
						Name: col,
					})
				}
			}

			for rows.Next() {
				datas := make([]*string, len(columns))
				err = rows.ScanSlice(&datas)
				if err != nil {
					return err
				}
				records = append(records, datas)
			}
		}
	}

	engines, err := models.FindEngines()
	if err != nil {
		return err
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
		"limit":       limit,
		"curPage":     start / limit,
		"totalPage":   pager(total, limit),
		"affected":    affected,
		"pkIdx":       pkIdx,
		"curEngine":   engine.Name,
		"IsLogin":     c.IsLogin(),
	})
}

func pager(total, limit int) int {
	if total%limit == 0 {
		return total / limit
	}
	return total/limit + 1
}

func curPage(start, limit int) int {
	return start / limit
}
