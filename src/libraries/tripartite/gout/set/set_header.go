package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"time"
)

type Header struct {
	Name     string    `SetHeader:"name"`
	Age      int       `SetHeader:"age"`
	Weight   float32   `SetHeader:"weight"`
	Birthday time.Time `SetHeader:"birthday"`
}

func SetHeaderByMap() {
	err := gout.GET("example.com").
		Debug(true).
		SetHeader(gout.H{
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

func SetHeaderByArray() {
	err := gout.GET("example.com").
		Debug(true).
		SetHeader(gout.A{
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

func SetHeaderByStruct() {
	err := gout.GET("example.com").
		Debug(true).
		SetHeader(Header{
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
