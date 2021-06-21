//
// Created by Rustle Karl on 2020.12.02 08:56.
//

package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()

	c.AddFunc("30 * * * *", func() {
		fmt.Println("Every hour on the half hour")
	})

	c.AddFunc("30 3-6,20-23 * * *", func() {
		fmt.Println(".. in the range 3-6am, 8-11pm")
	})

	c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() {
		fmt.Println("Runs at 04:30 Tokyo time every day")
	})

	c.AddFunc("@hourly", func() {
		fmt.Println("Every hour, starting an hour from now")
	})

	c.AddFunc("@every 1h30m", func() {
		fmt.Println("Every hour thirty, starting an hour thirty from now")
	})

	c.Start()

	//// Funcs are invoked in their own goroutine, asynchronously.
	//// Funcs may also be added to a running Cron
	//c.AddFunc("@daily", func() { fmt.Println("Every day") })
	//
	//c.Stop() // Stop the scheduler (does not stop any jobs already running).
}
