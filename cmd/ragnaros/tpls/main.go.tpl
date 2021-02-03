package main

import (
	"github.com/shumybest/ragnaros"
	"{{ .App.ProjectName }}/app"
)

func main() {
	ragnaros.InjectApps(app.{{ Export .App.ControllerName }})
	ragnaros.Start("{{ .App.ProjectName }}")
}
