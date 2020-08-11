package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ntt360/gracehttp"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	randPort := rand.Intn(20000) + 10000
	addr := fmt.Sprintf("0.0.0.0:%d", randPort)
	err := gracehttp.ListenAndServe(addr, gin.Default())
	if err != nil {
		log.Fatal(err)
	}
}
