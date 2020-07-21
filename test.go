package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	//prjHome := flag.String("prj_home", "", "project home path need")
	//env := flag.String("env", "", "env path")
	flag.Parse()
	fmt.Println(flag.Args())

	r := gin.Default()
	err := r.Run("0.0.0.0:9456")
	if err != nil {
		log.Fatal(err)
	}
}
