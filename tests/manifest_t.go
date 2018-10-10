package main

import (
	"fmt"

	"github.com/sumorf/ngui"
)

func main() {
	manifest := new(ngui.Manifest)
	manifest.Load()
	fmt.Printf("FirstPage=%v\n", manifest.FirstPage())
	fmt.Printf("LaunchWidth=%v\n", manifest.LaunchWidth())

	//fmt.Printf("%v\n", a)
}
