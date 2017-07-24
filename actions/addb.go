package actions

import (
	"fmt"
	"strings"
	"github.com/Unknwon/i18n"
	"github.com/go-xorm/dbweb/models"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
)

type Addb struct {
	AuthRenderBase

	binding.Binder
	xsrf.Checker
	flash.Flash
}

func (c *Addb) Get() error {
	engines, err := c.findEngines()
	if err != nil {
		return err
	}

	return c.Render("add.html", renders.T{
		"dbs":          SupportDBs,
		"flash":        c.Flash.Data(),
		"engines":      engines,
		"XsrfFormHtml": c.XsrfFormHtml(),
		"IsLogin":      c.IsLogin(),
		"isAdd":        true,
	})
}

func (c *Addb) Post() {
	var engine models.Engine
	engine.Name = c.Form("name")
	engine.Driver = c.Form("driver")
	host := c.Form("host")
	port := c.Form("port")
	dbname := c.Form("dbname")
	username := c.Form("username")
	passwd := c.Form("passwd")
	charset := c.Form("charset")
	if charset != strings.ToLower("utf8") && charset != strings.ToLower("utf8mb4") {
		c.Flash.Set("ErrAdd", i18n.Tr(c.CurLang(), "charset must be utf8 or utf8mb4"))
		c.Redirect("/addb")
		return
	}
	
	fmt.Printf("charset: %s\n", charset)
	if engine.Driver == "sqlite3" {
		engine.DataSource = host
	} else {
		engine.DataSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
			username, passwd, host, port, dbname, charset)
	}

	/*if err := c.MapForm(&engine); err != nil {
		c.Flash.Set("ErrAdd", i18n.Tr(c.CurLang(), "err_param"))
		c.Redirect("/addb")
		return
	}*/

	if err := models.AddEngine(&engine); err != nil {
		c.Flash.Set("ErrAdd", i18n.Tr(c.CurLang(), "err_add_failed"))
		c.Redirect("/addb")
		return
	}

	c.Redirect("/")
}
