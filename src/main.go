package main

import (
	"argovue/app"
	"fmt"
)

var version = "devel"
var commit string
var builddate string

func main() {
	fmt.Printf("ArgoVue, version:%s, commit:%s, builddate:%s\n", version, commit, builddate)
	app.New().SetVersion(version, commit, builddate).Wait()
}
