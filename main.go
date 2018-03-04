package main

import (
	"time"

	"math/rand"

	"github.com/buraksezer/gdoc/cmd"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	cmd.Execute()
}
