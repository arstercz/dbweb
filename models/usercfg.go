package models

//read and verify user auth
import (
	"log"
	"strings"
	"github.com/arstercz/goconfig"
)

type usermsg struct {
	name string
	email string
	pass string
	secret string
	dbs  map[string]bool
}

func Get_config(conf string) (c *goconfig.ConfigFile, err error) {
	c, err = goconfig.ReadConfigFile(conf)
	if err != nil {
		return c, err
	}
	return c, nil
}

func get_sections(c *goconfig.ConfigFile) ([]string){
	sections := c.GetSections()
	return sections
}

func Get_backend_dsn(c *goconfig.ConfigFile) (dsn string, err error) {
	dsn, err = c.GetString("backend", "dsn")
	if err 	!= nil {
		return dsn, err 
	}   
	return dsn, nil 
}

func get_users(c *goconfig.ConfigFile, section string) (*usermsg, error) {
	email := section
	name, err := c.GetString(section, "name")
	pass, err := c.GetString(section, "pass")
	secret, err := c.GetString(section, "secret")
	db, err := c.GetString(section, "db")
	if err != nil {
		log.Printf("parse cfg for user : %s error: %v", email, err)
		return nil, err
	}
	var dbs = make(map[string]bool)
	if len(db) > 0 {
		array1 := strings.Split(db, ", ")
		for _, v := range array1 {
			dbs[v] = true
		}
	} else {
		dbs["all"] = true
	}
	return &usermsg{name, email, pass, secret, dbs}, nil
}
