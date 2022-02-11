//
// Created by Rustle Karl on 2020.11.16 16:04.
//

package main

import (
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"log"
)

func main() {
	wb, err := xlsx.OpenFile("storage/io.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	// 遍历表
	for _, sheet := range wb.Sheets {
		fmt.Println(sheet.Name)
		// 遍历行读取
		for i := 0; i < sheet.MaxRow; i++ {
			// 遍历每行的列读取
			for j := 0; j < sheet.MaxCol; j++ {
				cell, _ := sheet.Cell(i, j)
				fmt.Print(cell.String(), "\t")
			}
			fmt.Println()
		}
		break
	}
}
