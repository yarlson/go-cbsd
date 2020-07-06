package main

import (
	"fmt"
	"github.com/yarlson/go-cbsd"
)

func main() {
	c := cbsd.NewCBSD()
	bHyves, err := c.BHyve.List()
	if err != nil {
		panic(err)
	}

	for _, bHyve := range bHyves {
		fmt.Println(bHyve.JName)
	}
}
