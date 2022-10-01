package main

import "debug/elf"

func main() {
	info, err := elf.Open("main")
	if err != nil {
		panic(err)
	}
	defer info.Close()
	for _, section := range info.Sections {
		println(section.Name)
	}
}
