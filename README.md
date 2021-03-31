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

```

The project is under heavy active development, feel free to create a fork and submit a pull request.
