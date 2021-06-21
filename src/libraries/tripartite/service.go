package main

import (
	"github.com/kardianos/service"
	"os"
	"time"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	// 运行逻辑
	for {
		println("running...")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "A", //服务显示名称
		DisplayName: "A", //服务名称
		Description: "A", //服务描述
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			if err = s.Install(); err != nil {
				panic(err)
			}
			println("服务安装成功")
			if err = s.Start(); err != nil {
				panic(err)
			}
			println("服务启动成功")
		case "start":
			if err = s.Start(); err != nil {
				panic(err)
			}
			println("服务启动成功")
		case "stop":
			if err = s.Stop(); err != nil {
				panic(err)
			}
			println("服务关闭成功")
		case "restart":
			if err = s.Stop(); err != nil {
				panic(err)
			}
			println("服务关闭成功")
			if err = s.Start(); err != nil {
				panic(err)
			}
			println("服务启动成功")
		case "remove", "uninstall":
			if err = s.Stop(); err != nil {
				panic(err)
			}
			println("服务关闭成功")
			if err = s.Uninstall(); err != nil {
				panic(err)
			}
			println("服务卸载成功")
		}
		return
	}

	err = s.Run()

	if err != nil {
		panic(err)
	}
}
