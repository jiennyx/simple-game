package main

import "simplegame.com/simplegame/userservice/server/appx"

func main() {
	app := appx.NewApplication()
	go app.Run()
	app.WaitShutdown()
}
