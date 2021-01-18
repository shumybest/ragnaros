package main

import (
	"github.com/shumybest/ragnaros2"
)

func main() {
	ragnaros.InjectApps(DemoController, func(r *ragnaros.Context) {
		r.Logger.Warn("Welcome to use Ragnaros")
	})
	ragnaros.Start()
}
