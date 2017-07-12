package actions

import (
	"fmt"
	_ "log"
	"regexp"
	"strings"
	"strconv"
	"github.com/go-xorm/xorm"
)

// find table name from sql string
func findTableFromSql(sql string) string {
	pattern1 := regexp.MustCompile(`(?i)ALTER TABLE (.+?)\s+`)
	pattern2 := regexp.MustCompile(`(?i)^CREATE INDEX .+\s+ON\s+(.+?)\s+`)
	lines := strings.Split(sql, "\n")
	var table string = ""
	for _, line := range lines {
		r1 := pattern1.FindStringSubmatch(line)
		r2 := pattern2.FindStringSubmatch(line)
		if len(r1) > 0 { 
	   		table = r1[1]
		}   
		if len(r2) > 0 { 
			table = r2[1]
		}   
    	}
	table = strings.Replace(table, "`", "", -1)
	return table
}

func getDatabaseName(o *xorm.Engine) string {
	sql := "select database() as dbname"
	var dbname string
	err := o.DB().QueryRow(sql).Scan(&dbname)
	if err != nil {
		return ""
	}
	return dbname
}

func getTableSize(o *xorm.Engine, db string, table string) uint64 {
	sql := fmt.Sprintf("select round(sum(DATA_LENGTH+INDEX_LENGTH+DATA_FREE)/1024/1024) as size from information_schema.tables where table_schema = '%s' and table_name = '%s'", db, table)
	var size uint64
	var tsize string
	err := o.DB().QueryRow(sql).Scan(&tsize)
	if err != nil {
		size = 0
	}
	size, err = strconv.ParseUint(tsize, 10, 64)
	if err != nil {
		size = 0
	}
	return size
}
