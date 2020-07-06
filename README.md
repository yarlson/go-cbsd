# go-cbsd
![Build](https://github.com/yarlson/go-cbsd/workflows/Build/badge.svg)

go-cbsd is a Go client library for managing [CBSD](https://github.com/cbsd/cbsd).

## Installation

```shell script
go get -u github.com/yarlson/go-cbsd
```

## Example
```go
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

```