package god

import (
	"fmt"
	"github.com/ntt360/pmon2/app/model"
)

var MonQueue chan *model.Process

type Monitor struct {
}

func NewMonitor() {
	MonQueue = make(chan *model.Process)

	go runMonitor()
}

func runMonitor() {
	for {
		select {
		case process := <-MonQueue:
			fmt.Println(process.MustJson())
		}
	}
}
