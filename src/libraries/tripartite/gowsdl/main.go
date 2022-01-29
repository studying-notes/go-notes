/*
 * @Date: 2021.08.25 9:20
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2021.08.25 9:20
 */

package main

import (
	"fmt"
	"github.com/hooklift/gowsdl/soap"
	"github/fujiawei-dev/go-notes/src/libraries/tripartite/gowsdl/gen"
)

//func ExampleBasicUsage() {
//	client := soap.NewClient("http://10.1.76.75/eecmisws/services/eecmisws?wsdl")
//	service := gen.NewStockQuotePortType(client)
//	reply, err := service.GetLastTradePrice(&gen.TradePriceRequest{})
//	if err != nil {
//		log.Fatalf("could't get trade prices: %v", err)
//	}
//	log.Println(reply)
//}
//
//func ExampleWithOptions() {
//	client := soap.NewClient(
//		"http://svc.asmx",
//		soap.WithTimeout(time.Second*5),
//		soap.WithBasicAuth("usr", "psw"),
//		soap.WithTLS(&tls.Config{InsecureSkipVerify: true}),
//	)
//	service := gen.NewStockQuotePortType(client)
//	reply, err := service.GetLastTradePrice(&gen.TradePriceRequest{})
//	if err != nil {
//		log.Fatalf("could't get trade prices: %v", err)
//	}
//	log.Println(reply)
//}

func main() {
	client := soap.NewClient("http://10.1.76.75/eecmisws/services/eecmisws?wsdl")
	service := gen.NewEecmiswsPortType(client)

	loginName := "loginName"
	password := "password"

	loginResponse, err := service.Login(&gen.Login{LoginName: &loginName, Password: &password})
	if err != nil {
		panic(err)
	}

	sessionId := loginResponse.Return_
	if sessionId != nil {
		fmt.Println(sessionId)
	}
}
