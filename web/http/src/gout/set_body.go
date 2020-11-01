package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"strings"
)

func SetBodyByString() {
	err := gout.POST("example.com").Debug(true).SetBody("string").Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func SetBodyByReader() {
	err := gout.POST("example.com").Debug(true).SetBody(strings.NewReader("io.Reader")).Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func SetBodyByBaseType() {
	err := gout.POST("example.com").Debug(true).SetBody(3.14).Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
