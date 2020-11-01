package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"time"
)

type Query struct {
	Name     string    `query:"name"`
	Age      int       `query:"age"`
	Weight   float32   `query:"weight"`
	Birthday time.Time `query:"birthday"`
}

func QueryByMap() {
	err := gout.GET("example.com").
		Debug(true).
		SetQuery(gout.H{
			"name":     "user",
			"age":      18,
			"weight":   50.4,
			"birthday": time.Now(),
		}).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func QueryByArray() {
	err := gout.GET("example.com").
		Debug(true).
		SetQuery(gout.A{
			"name", "user",
			"age", 18,
			"weight", 50.4,
			"birthday", time.Now(),
		}).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func QueryByStruct() {
	err := gout.GET("example.com").
		Debug(true).
		SetQuery(Query{
			Name:     "user",
			Age:      18,
			Weight:   50.4,
			Birthday: time.Now(),
		}).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func QueryByString() {
	err := gout.GET("example.com").Debug(true).
		SetQuery("name=user&age=18&weight=50.5&birthday=2020-1-20").Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
