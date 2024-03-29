Automate your Windows machine in Go.

    import "github.com/gonutz/auto"

Mouse functions:

    // Replace Left with Right or Middle.
    err := auto.ClickLeftMouseAt(x, y)
    err := auto.ClickLeftMouse()
    err := auto.PressLeftMouseAt(x, y)
    err := auto.PressLeftMouse()
    err := auto.ReleaseLeftMouseAt(x, y)
    err := auto.ReleaseLeftMouse()
    err := auto.MoveMouseTo(x, y)
    err := auto.MoveMouseBy(relativeX, relativeY)
	x, y, err := auto.MousePosition()
	err := auto.MoveMouseWheelBy(dx, dy)

Keyboard functions:

    err := auto.Type("Hello")
    err := auto.TypeWithDelay("Hello", 100 * time.Millisecond)
    err := auto.TypeKey(auto.KeySpace)
    err := auto.PressKey(auto.KeySpace)
    err := auto.ReleaseKey(auto.KeySpace)

Screen shot functions:

    img, err := auto.CaptureMonitor(Monitor)
    img, err := auto.CaptureMonitors([]Monitor)
    img, err := auto.CaptureScreen(x, y, width, height int)
    img, err := auto.CaptureScreenRect(Rectangle)
    img, err := auto.CaptureWindow(Window)
    img, err := auto.CaptureWindowContent(Window)

Monitor functions:

    allMonitors, err := auto.Monitors()
    m, err := auto.PrimaryMonitor()

Window functions:

    allWindows, err := auto.Windows()
    window, err := auto.ForegroundWindow()
    err := window.BringToForeground()
    x, y, width, height, err := window.InnerPosition()
    err := window.SetInnerPosition(x, y, width, height)
    x, y, width, height, err := window.OuterPosition()
    err := window.SetOuterPosition(x, y, width, height)
    window.Restore()
    window.Maximize()
    window.Minimize()
    window.Hide()
    window.Show()
    window.Update()

Global Events:

    auto.SetOnKeyboardEvent(func(*auto.KeyboardEvent))
    auto.SetOnMouseEvent(func(*auto.MouseEvent))
    auto.SetOnClipboardChange(func())

Other OS functions:

    text, err := auto.ClipboardText()
    err := auto.SetClipboardText("Hello")
	ShowMessage(caption, message string)
