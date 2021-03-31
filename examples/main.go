package main

import (
	"context"
	"fmt"
	"github.com/yarlson/go-cbsd/v2"
)

func main() {
	c := cbsd.NewCBSD()
	bHyves, err := c.BHyve.List(context.Background())
	if err != nil {
		panic(err)
	}

	for _, bHyve := range bHyves {
		fmt.Println(bHyve.JName)
	}
}
