package main

import (
	"github.com/alexandru197/company/adapters"
)

func main() {
	app := adapters.NewCompanyApp()
	app.Start()
}
