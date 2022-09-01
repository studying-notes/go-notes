package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	fmt.Println(strings.Title("Golang is awesome!"))

	caser := cases.Title(language.English)
	fmt.Println(caser.String("here comes o'brian"))
}
