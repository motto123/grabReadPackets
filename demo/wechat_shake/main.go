package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"time"
	"fmt"
	"os"
	"log"
)

const (
	giftRealSmal  = iota
	giftRealLarge
	gitCoin
	gitCoupon
)

// 最大号码
const rateMax = 10000

type gift struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	GiftCategory int    `json:"gift_category"`
	Total        int    `json:"total"`
	Surplus      int    `json:"surplus"`
	Data         string `json:"data"`
	Usable       bool   `json:"usable"`   // 是否使用中
	Rate         int    `json:"rate"`     // 中奖概率，万分之N,0-10000
	RateMin      int    `json:"rate_min"` // 大于等于，中奖的最小号码,0-10000
	RateMax      int    `json:"rate_max"` // 小于，中奖的最大号码,0-10000
}

func (g gift) String() string {
	return fmt.Sprintf("{name: %s, total: %d, surplus: %d}", g.Name, g.Total, g.Surplus)
}

type lotteryController struct {
	Ctx iris.Context
}

var logger *log.Logger

// 初始化日志信息
func initLog() {
	f, _ := os.Create("./log/lottery_demo.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

var giftList []*gift

func initGift() {
	giftList = make([]*gift, 4)
	// 1 实物大奖
	g1 := gift{
		Id:           1,
		Name:         "iphone x",
		Data:         "",
		GiftCategory: giftRealLarge,
		Usable:       true,
		Total:        2,
		Surplus:      2,
		Rate:         2,
		RateMax:      0,
		RateMin:      0,
	}
	giftList[0] = &g1
	// 2 实物小奖
	g2 := gift{
		Id:           1,
		Name:         "charge pal",
		Data:         "",
		GiftCategory: giftRealSmal,
		Usable:       true,
		Total:        5,
		Surplus:      5,
		Rate:         100,
		RateMax:      0,
		RateMin:      0,
	}
	giftList[1] = &g2
	// 3 虚拟券，相同的编码
	g3 := gift{
		Id:           1,
		Name:         "$10 coupon",
		Data:         "",
		GiftCategory: gitCoupon,
		Usable:       true,
		Total:        50,
		Surplus:      50,
		Rate:         500,
		RateMax:      0,
		RateMin:      0,
	}
	giftList[2] = &g3
	// 5 虚拟币
	g4 := gift{
		Id:           1,
		Name:         "10 coin",
		Data:         "",
		GiftCategory: gitCoupon,
		Usable:       true,
		Total:        5000,
		Surplus:      5000,
		Rate:         5000,
		RateMax:      0,
		RateMin:      0,
	}
	giftList[3] = &g4

	// 整理奖品数据，把rateMin,rateMax根据rate进行编排
	rateStart := 0
	for _, data := range giftList {
		if !data.Usable {
			continue
		}
		data.RateMin = rateStart
		data.RateMax = data.RateMin + data.Rate
		if data.RateMax >= rateMax {
			// 号码达到最大值，分配的范围重头再来
			data.RateMax = rateMax
			rateStart = 0
		} else {
			rateStart += data.Rate
		}
	}
	fmt.Printf("giftlist=%v\n", giftList)
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	rand.Seed(time.Now().UnixNano())
	initLog()
	initGift()
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8888"))
}

// GET http://localhost:8080/
func (c *lotteryController) Get() string {
	count := 0
	total := 0
	for _, data := range giftList {
		if data.Usable && (data.Total == 0 ||
			(data.Total > 0 && data.Surplus > 0)) {
			count++
			total += data.Surplus
		}
	}
	return fmt.Sprintf("当前有效奖品种类数量: %d，限量奖品总数量=%d\n", count, total)
}

func (c *lotteryController) GetLucky() map[string]interface{} {
	luckyCode := rand.Intn(rateMax) + 1
	m := make(map[string]interface{})
	var ok bool
	var sendData string
	for _, gift := range giftList {
		if !gift.Usable || (gift.Total > 0 && gift.Surplus <= 0) {
			continue
		}
		if luckyCode <= rateMax && gift.Rate >= luckyCode {
			ok, sendData = sendGift(gift)
			saveLuckyData(int32(luckyCode), gift.Id, gift.Name, "", gift.Data, gift.Surplus)
			break
		}
	}
	m["success"] = ok
	m["data"] = sendData
	return m
}

func sendGift(gift *gift) (bool, string) {
	if gift.Total == 0 {
		return true, gift.Name
	} else if gift.Surplus > 0 {
		gift.Surplus -= 1
		return true, gift.Name
	}
	return false, "no prizes"
}

// 记录用户的获奖记录
func saveLuckyData(code int32, id int, name, link, sendData string, left int) {
	logger.Printf("lucky, code=%d, gift=%d, name=%s, link=%s, data=%s, left=%d ", code, id, name, link, sendData, left)
}
