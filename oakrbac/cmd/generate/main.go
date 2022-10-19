package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

var (
	debug  = pflag.BoolP("debug", "d", false, "inject debugging traces")
	output = pflag.StringP("output", "o", "oakacs.gen.go", "the name of the generated role and policy file")
)

func main() {
	pflag.Parse()

	fmt.Println("policy generation tool that uses OakACS schema")
}
