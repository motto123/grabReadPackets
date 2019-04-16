package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"fmt"
	"strings"
	"math/rand"
	"time"
)

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lottertyController{})
	return app
}

type lottertyController struct {
	Ctx iris.Context
}

var userList []string

func main() {
	app := newApp()
	app.Run(iris.Addr(":8888"))
	userList = make([]string, 0)
}

func (c *lottertyController) Get() string {
	return fmt.Sprintf("total %d", len(userList))
}

func (c *lottertyController) PostImport() string {
	usersStr := c.Ctx.FormValue("users")
	users := strings.Split(usersStr, ",")
	cnt1 := len(userList)
	for _, v := range users {
		u := strings.TrimSpace(v)
		if len(u) > 0 {
			userList = append(userList, u)
		}
	}
	cnt2 := len(userList)
	return fmt.Sprintf("total %d import %d", cnt2, cnt2-cnt1)
}

func (c *lottertyController) GetLucky() string {
	luckyUser := ""
	if len(userList) == 0 {
		luckyUser = "no user"
	} else if len(userList) == 1 {
		luckyUser = userList[0]
	} else {
		index := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(len(userList)))
		luckyUser = userList[index]
		userList = append(userList[:index], userList[index+1:]...)
	}
	return luckyUser
}
