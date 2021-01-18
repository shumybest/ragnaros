package main

import (
	"github.com/shumybest/ragnaros2"
	"{{ .App.ProjectName }}/app"
)

func main() {
	ragnaros.InjectApps(app.{{ Export .App.ControllerName }})
	ragnaros.Start()
}
