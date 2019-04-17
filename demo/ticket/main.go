package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"time"
	"fmt"
)

type lottertyController struct {
	kCtx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lottertyController{})
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8888"))
	rand.Seed(time.Now().UnixNano())
}

func (c *lottertyController) Get() string {
	code := rand.Intn(10)
	prize := ""
	if code == 1 {
		prize = "no1"
	} else if code > 1 && code < 4 {
		prize = "no2"
	} else if code > 3 && code < 6 {
		prize = "no3"
	} else {
		prize = "no lucky"
	}
	return prize
}

func (c *lottertyController) GetPrize() string {
	readBalls := make([]int, 7)
	for i := 0; i < len(readBalls)-1; i++ {
		readBalls[i] = rand.Intn(33) + 1
	}
	readBalls[6] = rand.Intn(16) + 1
	return fmt.Sprintf("%v", readBalls)
}
