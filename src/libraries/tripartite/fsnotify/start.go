package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

func main() {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file: ", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error: ", err)
			}
		}
	}()
	_ = watcher.Add("D:\\Learning")
	<-done
}
