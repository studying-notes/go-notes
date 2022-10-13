package main

import (
	"reflect"
)

type User struct {
	X int
	y float64
}

func main() {
	var s = User{X: 1, y: 2.0}
	rValue := reflect.ValueOf(&s).Elem()
	rValueX := rValue.Field(0)
	rValueX.SetInt(100)
	rValueY := rValue.FieldByName("y")
	rValueY.SetFloat(200.0)
	println(s.X, s.y)
}
