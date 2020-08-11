package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ntt360/gracehttp"
	"log"
)

func main() {
	//prjHome := flag.String("prj_home", "", "project home path need")
	//env := flag.String("env", "", "env path")
	flag.Parse()
	fmt.Println(flag.Args())

	err := gracehttp.ListenAndServe("0.0.0.0:9456", gin.Default())
	if err != nil {
		log.Fatal(err)
	}
}
