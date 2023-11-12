package main

import (
	"time"

	"github.com/gonutz/auto"
)

func main() {
	auto.ShowMessage(
		"Clicker",
		`Right-click: Start/stop clicking left forever.
       Escape: Quit.`,
	)
	defer auto.ShowMessage("Clicker", `Quitting.`)

	clicking := false
	go func() {
		for {
			if clicking {
				auto.ClickLeftMouse()
				time.Sleep(time.Nanosecond)
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	auto.SetOnMouseEvent(func(e *auto.MouseEvent) {
		if e.Type == auto.RightMouseDown {
			e.Cancel()
			clicking = !clicking
		}
		if e.Type == auto.RightMouseUp {
			e.Cancel()
		}
	})

	stop := make(chan bool)
	auto.SetOnKeyboardEvent(func(e *auto.KeyboardEvent) {
		if e.Down && e.Key == auto.KeyEscape {
			e.Cancel()
			stop <- true
		}
	})
	<-stop

	auto.SetOnMouseEvent(nil)
	auto.SetOnKeyboardEvent(nil)
}
