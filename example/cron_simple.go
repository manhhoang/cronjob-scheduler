package main

import (
	"gihub.com/cronjob-scheduler/cron"
	"gihub.com/cronjob-scheduler/utility"
)

func main() {
	//cr := cron.New(cron.WithSeconds())
	c := cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	c.Start()
	printfLogger := utility.DefaultLogger
	printfLogger.Info("Test")
}

