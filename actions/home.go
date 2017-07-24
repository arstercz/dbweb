package actions

import (
	"github.com/go-xorm/core"
	_ "github.com/go-xorm/dbweb/models"
	"github.com/tango-contrib/renders"
)

type Home struct {
	AuthRenderBase
}

func (c *Home) Get() error {
	engines, err := c.findEngines()
	if err != nil {
		return err
	}

	return c.Render("root.html", renders.T{
		"engines": engines,
		"tables":  []core.Table{},
		"records": [][]string{},
		"columns": []string{},
		"id":      0,
		"ishome":  true,
		"IsLogin": c.IsLogin(),
	})
}
