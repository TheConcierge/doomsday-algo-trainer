package main

import (
	"github.com/TheConcierge/doomsday-algo-trainer/core"
	"github.com/TheConcierge/doomsday-algo-trainer/injections/sessionmanager/dummy"
	promptuiInput "github.com/TheConcierge/doomsday-algo-trainer/inputs/cli/promptui/v1"
)

func main() {
	sess := dummy.NewDummySessionManager()

	c := core.NewDoomsdayCore(sess)

	pui := promptuiInput.NewPromptUi(c)

	pui.Start()
}
