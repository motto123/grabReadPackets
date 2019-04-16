package main

import (
	"testing"
	"github.com/kataras/iris/httptest"
	"fmt"
	"sync"
)

func TestLucky(t *testing.T) {
	e := httptest.New(t, newApp())
	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("total 0")

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e.POST("/import").WithFormField("users", fmt.Sprintf("test_u%d", i)).
				Expect().Status(httptest.StatusOK).Body()
		}(i)
	}

	wg.Wait()

	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("total 100")

	e.GET("/lucky").Expect().Status(httptest.StatusOK).Body()

	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("total 99")
}
