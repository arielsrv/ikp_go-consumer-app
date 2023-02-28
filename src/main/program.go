package main

import (
	"github.com/src/main/app"
	"github.com/src/main/app/log"
	_ "github.com/src/resources/docs"
)

// @title Golang Template API
// @description This is a sample golang template api. Have fun.
// @basePath /
// @version v1.
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
