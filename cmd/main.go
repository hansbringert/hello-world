package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	box "github.com/hansbringert/hello-world/internal"
)

var version = "development"
var commit = "-"
var date = "-"
var builtBy = "-"

var infoVariables = map[string]string{
	"Version:": version,
	"Commit:":  commit,
	"Date:":    date,
	"BuiltBy:": builtBy,
	"Runtime:": runtime.Version(),
}

func printVersion() {
	for key, value := range infoVariables {
		fmt.Printf("%-10s %s\n", key, value)
	}
}

var rootDir = "./static"

func main() {
	printVersion()
	fileServer := box.FileServer(box.Dir(rootDir))

	//fs := box.MyFileSystem{
	//	Dir: rootDir,
	//}
	//fileServer := http.FileServer(fs)

	http.Handle("/", fileServer)

	// Start server
	if err := http.ListenAndServe(":1959", nil); err != nil {
		log.Fatal(err)
	}
}
