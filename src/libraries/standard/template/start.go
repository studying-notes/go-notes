package main

import (
	"os"
	"strings"
	"text/template"
)

func main() {
	text := "1. {{ title .user }}\n2. {{ .tag | title}}"
	funcMap := template.FuncMap{"title": strings.Title}
	tmpl, _ := template.New("start").Funcs(funcMap).Parse(text)
	data := map[string]string{
		"user": "rustle",
		"tag":  "admin",
	}
	_ = tmpl.Execute(os.Stdout, data)
}
