package main

import (
	"fmt"
	"go/scanner"
	"go/token"
)

func main() {
	src := []byte("cos(x) + 2i*sin(x) // Euler")

	var s scanner.Scanner
	fileSet := token.NewFileSet()
	file := fileSet.AddFile("", fileSet.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", file.Position(pos), tok, lit)
	}
}
