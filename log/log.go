package log

import "github.com/B9O2/Inspector/inspect"

var InspLog = inspect.NewInspector("rhchannel", 9999)

func initDecoration() {
	InspLog.SetSeparator("")
	InspLog.SetVisible(false)
}
