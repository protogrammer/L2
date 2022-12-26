package main

import (
	"server/db"
	"server/interruption"
	"server/server"
)

func main() {
	defer db.Close()
	go server.Run()
	interruption.Wait()
}
