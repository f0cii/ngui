// Copyright (c) 2014 The ngui authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/nvsoft/ngui

package main

import (
	"github.com/nvsoft/ngui"
)

func main() {
	app := ngui.NewApplication()

	app.CreateBrowser()

	app.Exec()
}
