// Copyright (c) 2014 The ngui authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/nvsoft/ngui

package main

import (
	"os"
	"github.com/nvsoft/ngui"
)

func main() {
	app := ngui.NewEngine()

	// TODO: It should be executable's directory used
	// rather than working directory.
	url, _ := os.Getwd()
	url = "file://" + url + "/example.html"
	app.CreateWindow(url)

	app.Exec()
}

