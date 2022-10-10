package auto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/gonutz/w32/v2"
)

var (
	errBlocked   = errors.New("SendInput returned 0, meaning input was blocked")
	errSetCursor = errors.New("SetCursorPos failed")
)

// ClickLeftMouseAt moves the mouse to screen coordinates x,y and clicks the
// left mouse button, i.e. presses and releases it.
func ClickLeftMouseAt(x, y int) error {
	return clickAt(x, y, w32.MOUSEEVENTF_LEFTDOWN, w32.MOUSEEVENTF_LEFTUP)
}

// ClickLeftMouse clicks the left mouse button, i.e. presses and releases it.
func ClickLeftMouse() error {
	return click(w32.MOUSEEVENTF_LEFTDOWN, w32.MOUSEEVENTF_LEFTUP)
}

// PressLeftMouseAt moves the mouse to screen coordinates x,y and presses the
// left mouse button down. Call ReleaseLeftMouse or ReleaseLeftMouseAt to
// release the button.
func PressLeftMouseAt(x, y int) error {
	return mouseInputAt(x, y, w32.MOUSEEVENTF_LEFTDOWN)
}

// PressLeftMouse presses the left mouse button down. Call ReleaseLeftMouse or
// ReleaseLeftMouseAt to release the button.
func PressLeftMouse() error {
	return mouseInput(w32.MOUSEEVENTF_LEFTDOWN)
}

// ReleaseLeftMouseAt moves the mouse to screen coordinates x,y and releases the
// left mouse button. You probably want to press it before, using
// PressLeftMouseAt or PressLeftMouse.
func ReleaseLeftMouseAt(x, y int) error {
	return mouseInputAt(x, y, w32.MOUSEEVENTF_LEFTUP)
}

// ReleaseLeftMouse releases the left mouse button. You probably want to press
// it before, using PressLeftMouseAt or PressLeftMouse.
func ReleaseLeftMouse() error {
	return mouseInput(w32.MOUSEEVENTF_LEFTUP)
}

// ClickRightMouseAt moves the mouse to screen coordinates x,y and clicks the
// right mouse button, i.e. presses and releases it.
func ClickRightMouseAt(x, y int) error {
	return clickAt(x, y, w32.MOUSEEVENTF_RIGHTDOWN, w32.MOUSEEVENTF_RIGHTUP)
}

// ClickRightMouse clicks the right mouse button, i.e. presses and releases it.
func ClickRightMouse() error {
	return click(w32.MOUSEEVENTF_RIGHTDOWN, w32.MOUSEEVENTF_RIGHTUP)
}

// PressRightMouseAt moves the mouse to screen coordinates x,y and presses the
// right mouse button down. Call ReleaseRightMouse or ReleaseRightMouseAt to
// release the button.
func PressRightMouseAt(x, y int) error {
	return mouseInputAt(x, y, w32.MOUSEEVENTF_RIGHTDOWN)
}

// PressRightMouse presses the right mouse button down. Call ReleaseRightMouse or
// ReleaseRightMouseAt to release the button.
func PressRightMouse() error {
	return mouseInput(w32.MOUSEEVENTF_RIGHTDOWN)
}

// ReleaseRightMouseAt moves the mouse to screen coordinates x,y and releases the
// right mouse button. You probably want to press it before, using
// PressRightMouseAt or PressRightMouse.
func ReleaseRightMouseAt(x, y int) error {
	return mouseInputAt(x, y, w32.MOUSEEVENTF_RIGHTUP)
}

// ReleaseRightMouse releases the right mouse button. You probably want to press
// it before, using PressRightMouseAt or PressRightMouse.
func ReleaseRightMouse() error {
	return mouseInput(w32.MOUSEEVENTF_RIGHTUP)
}

// ClickMiddleMouseAt moves the mouse to screen coordinates x,y and clicks the
// middle mouse button, i.e. presses and releases it.
func ClickMiddleMouseAt(x, y int) error {
	return clickAt(x, y, w32.MOUSEEVENTF_MIDDLEDOWN, w32.MOUSEEVENTF_MIDDLEUP)
}

// ClickMiddleMouse clicks the middle mouse button, i.e. presses and releases it.
func ClickMiddleMouse() error {
	return click(w32.MOUSEEVENTF_MIDDLEDOWN, w32.MOUSEEVENTF_MIDDLEUP)
}

// PressMiddleMouseAt moves the mouse to screen coordinates x,y and presses the
// middle mouse button down. Call ReleaseMiddleMouse or ReleaseMiddleMouseAt to
// release the button.
func PressMiddleMouseAt(x, y int) error {
	return mouseInputAt(x, y, w32.MOUSEEVENTF_MIDDLEDOWN)
}

// PressMiddleMouse presses the middle mouse button down. Call ReleaseMiddleMouse or
// ReleaseMiddleMouseAt to release the button.
func PressMiddleMouse() error {
	return mouseInput(w32.MOUSEEVENTF_MIDDLEDOWN)
}

// ReleaseMiddleMouseAt moves the mouse to screen coordinates x,y and releases the
// middle mouse button. You probably want to press it before, using
// PressMiddleMouseAt or PressMiddleMouse.
func ReleaseMiddleMouseAt(x, y int) error {
	return mouseInputAt(x, y, w32.MOUSEEVENTF_MIDDLEUP)
}

// ReleaseMiddleMouse releases the middle mouse button. You probably want to press
// it before, using PressMiddleMouseAt or PressMiddleMouse.
func ReleaseMiddleMouse() error {
	return mouseInput(w32.MOUSEEVENTF_MIDDLEUP)
}

// MoveMouseTo move the mouse cursor to the given screen coordinates.
func MoveMouseTo(x, y int) error {
	if !w32.SetCursorPos(x, y) {
		return errSetCursor
	}
	return nil
}

// MoveMouseBy moves the mouse cursor by the given amount of pixels in x and y.
// Positive x moves the cursor right.
// Negative x moves the cursor left.
// Positive y moves the cursor down.
// Negative y moves the cursor up.
func MoveMouseBy(dx, dy int) error {
	x, y, ok := w32.GetCursorPos()
	if !ok {
		return errors.New("GetCursorPos failed")
	}
	if !w32.SetCursorPos(x+dx, y+dy) {
		return errSetCursor
	}
	return nil
}

func clickAt(x, y int, down, up uint32) error {
	if !w32.SetCursorPos(x, y) {
		return errSetCursor
	}
	return click(down, up)
}

func click(down, up uint32) error {
	n := w32.SendInput(
		w32.MouseInput(w32.MOUSEINPUT{Flags: down}),
		w32.MouseInput(w32.MOUSEINPUT{Flags: up}),
	)
	if n == 0 {
		return errBlocked
	}
	return nil
}

func mouseInputAt(x, y int, flags uint32) error {
	if !w32.SetCursorPos(x, y) {
		return errSetCursor
	}
	return mouseInput(flags)
}

func mouseInput(flags uint32) error {
	n := w32.SendInput(
		w32.MouseInput(w32.MOUSEINPUT{Flags: flags}),
	)
	if n == 0 {
		return errBlocked
	}
	return nil
}

// Type will write the given text using Alt+Numpad numbers. It will sleep the
// smallest, non-0 delay between two letters.
func Type(s string) error {
	return TypeWithDelay(s, 1)
}

// TypeWithDelay will write the given text using Alt+Numpad numbers. It will
// sleep the given delay between two letters.
func TypeWithDelay(s string, delay time.Duration) error {
	toScanCode := func(vk uint) uint16 {
		return uint16(w32.MapVirtualKey(vk, w32.MAPVK_VK_TO_VSC))
	}

	const (
		down = 0
		up   = 1
	)

	upDown := func(vk uint) [2]w32.INPUT {
		return [2]w32.INPUT{
			down: w32.KeyboardInput(w32.KEYBDINPUT{
				Scan:  toScanCode(vk),
				Flags: w32.KEYEVENTF_SCANCODE,
			}),

			up: w32.KeyboardInput(w32.KEYBDINPUT{
				Scan:  toScanCode(vk),
				Flags: w32.KEYEVENTF_SCANCODE | w32.KEYEVENTF_KEYUP,
			}),
		}
	}

	alt := upDown(w32.VK_LMENU)
	nums := [][2]w32.INPUT{
		upDown(w32.VK_NUMPAD0),
		upDown(w32.VK_NUMPAD1),
		upDown(w32.VK_NUMPAD2),
		upDown(w32.VK_NUMPAD3),
		upDown(w32.VK_NUMPAD4),
		upDown(w32.VK_NUMPAD5),
		upDown(w32.VK_NUMPAD6),
		upDown(w32.VK_NUMPAD7),
		upDown(w32.VK_NUMPAD8),
		upDown(w32.VK_NUMPAD9),
	}

	keys := []w32.INPUT{alt[down], nums[0][down], nums[0][up]}

	// Unify line breaks to '\r' which is the virtual key code for VK_RETURN.
	s = strings.Replace(s, "\r\n", "\r", -1)
	s = strings.Replace(s, "\n", "\r", -1)

	for _, r := range s {
		if r == '\r' {
			if err := PressKey(w32.VK_RETURN); err != nil {
				return err
			}
		} else if r == '\b' {
			if err := PressKey(w32.VK_BACK); err != nil {
				return err
			}
		} else {
			keys = keys[:3] // Keep Alt down and type 0.
			for _, digit := range fmt.Sprint(int(r)) {
				d := digit - '0'
				keys = append(keys, nums[d][down], nums[d][up])
			}
			keys = append(keys, alt[up])

			if w32.SendInput(keys...) == 0 {
				return errBlocked
			}
		}
		time.Sleep(delay)
	}
	return nil
}

// PressKey presses the given key on the keyboard. You can pass key codes
// defined in this package, named Key...
func PressKey(key uint16) error {
	n := w32.SendInput(w32.KeyboardInput(w32.KEYBDINPUT{Vk: key}))
	if n == 0 {
		return errBlocked
	}
	return nil
}

// ReleaseKey releases the given key on the keyboard. You can pass key codes
// defined in this package, named Key...
func ReleaseKey(key uint16) error {
	n := w32.SendInput(w32.KeyboardInput(w32.KEYBDINPUT{
		Vk:    key,
		Flags: w32.KEYEVENTF_KEYUP,
	}))
	if n == 0 {
		return errBlocked
	}
	return nil
}

// TypeKey presses and releases the given key on the keyboard. The value must be
// a virtual keycode like 'A', '1' or VK_RETURN (you can use the constants in
// github.com/gonutz/w32 VK_...).
func TypeKey(key uint16) error {
	n := w32.SendInput(
		w32.KeyboardInput(w32.KEYBDINPUT{
			Vk: key,
		}),
		w32.KeyboardInput(w32.KEYBDINPUT{
			Vk:    key,
			Flags: w32.KEYEVENTF_KEYUP,
		}),
	)
	if n == 0 {
		return errBlocked
	}
	return nil
}

// ForegroundWindow returns the currently active window. If no window is active,
// ForegroundWindow returns an error.
func ForegroundWindow() (Window, error) {
	w := w32.GetForegroundWindow()
	if w == 0 {
		return Window{}, errors.New("no window is active")
	}
	return windowHandleToWindow(w), nil
}

// BringToForeground tries to bring the given window to the front.
func (w *Window) BringToForeground() error {
	if !w32.SetForegroundWindow(w.Handle) {
		return errors.New("SetForegroundWindow failed")
	}
	w.Update()
	return nil
}

// Restore unminimizes a minimized window and unmaximizes a maximized window.
func (w *Window) Restore() {
	w32.ShowWindow(w.Handle, w32.SW_RESTORE)
	w.Update()
}

// Maximize maximizes the given window.
func (w *Window) Maximize() {
	w32.ShowWindow(w.Handle, w32.SW_MAXIMIZE)
	w.Update()
}

// Minimize minimizes the given window.
func (w *Window) Minimize() {
	w32.ShowWindow(w.Handle, w32.SW_MINIMIZE)
	w.Update()
}

// Hide hides the window. Call ShowWindow to show it again.
func (w *Window) Hide() {
	w32.ShowWindow(w.Handle, w32.SW_HIDE)
	w.Update()
}

// Show shows the given window. Call this to show a window that was hidden with
// Hide.
func (w *Window) Show() {
	w32.ShowWindow(w.Handle, w32.SW_SHOW)
	w.Update()
}

// Update updates the state of the window, all fields are queried from the OS
// again. If the state or size of a window changes, Update will poll these
// changes.
func (w *Window) Update() {
	*w = windowHandleToWindow(w.Handle)
}

// Window is a window currently open on you system.
type Window struct {
	// Rectangle is the window's outer boundaries in virtual screen coordinates.
	Rectangle
	// Content is the window's inner boundaries in virtual screen coordinates.
	Content Rectangle
	// Visible is true if the window is a visual window. For background windows
	// Visible is false.
	Visible bool
	// Title is the text currently displayed in the window header.
	Title string
	// ClassName is the name of the class of this window. Multiple windows can
	// have the same class.
	ClassName string
	// Maximized is true if the window is currently maximized.
	Maximized bool
	// Minimized is true if the window is currently minimized.
	Minimized bool
	// Handle is the operating specific window handle.
	Handle w32.HWND
}

func windowHandleToWindow(window w32.HWND) Window {
	className, _ := w32.GetClassName(window)
	bounds := w32.GetWindowRect(window)
	client := w32.GetClientRect(window)
	clientLeft, clientTop := w32.ClientToScreen(window, 0, 0)
	var placement w32.WINDOWPLACEMENT
	w32.GetWindowPlacement(window, &placement)
	return Window{
		Handle:    window,
		Visible:   w32.IsWindowVisible(window),
		Title:     w32.GetWindowText(window),
		ClassName: className,
		Maximized: placement.ShowCmd == w32.SW_MAXIMIZE,
		Minimized: placement.ShowCmd == w32.SW_SHOWMINIMIZED ||
			placement.ShowCmd == w32.SW_MINIMIZE ||
			placement.ShowCmd == w32.SW_FORCEMINIMIZE,
		Rectangle: Rectangle{
			X:      int(bounds.Left),
			Y:      int(bounds.Top),
			Width:  int(bounds.Width()),
			Height: int(bounds.Height()),
		},
		Content: Rectangle{
			X:      int(clientLeft),
			Y:      int(clientTop),
			Width:  int(client.Width()),
			Height: int(client.Height()),
		},
	}
}

// ClipboardText returns the contents of the clipboard as text. If the clipboard
// is empty or does not contain text it returns "".
func ClipboardText() (string, error) {
	if !w32.OpenClipboard(0) {
		return "", errors.New("OpenClipboard failed")
	}
	defer w32.CloseClipboard()

	data := (*uint16)(unsafe.Pointer(w32.GetClipboardData(w32.CF_UNICODETEXT)))
	if data == nil {
		return "", errors.New("GetClipboardData failed")
	}

	var characters []uint16
	for *data != 0 {
		characters = append(characters, *data)
		data = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(data)) + 2))
	}
	return syscall.UTF16ToString(characters), nil
}

// SetClipboardText sets the contents of the clipboard to the given string.
func SetClipboardText(text string) error {
	if !w32.OpenClipboard(0) {
		return errors.New("OpenClipboard failed")
	}
	defer w32.CloseClipboard()

	if !w32.EmptyClipboard() {
		return errors.New("EmptyClipboard failed")
	}

	data := syscall.StringToUTF16(text)
	clipBuffer := w32.GlobalAlloc(w32.GMEM_DDESHARE, uint32(len(data)*2))
	if clipBuffer == 0 {
		return errors.New("GlobalAlloc failed")
	}
	defer w32.GlobalUnlock(clipBuffer)

	w32.MoveMemory(
		w32.GlobalLock(clipBuffer),
		unsafe.Pointer(&data[0]),
		uint32(len(data)*2),
	)

	if 0 == w32.SetClipboardData(
		w32.CF_UNICODETEXT,
		w32.HANDLE(clipBuffer),
	) {
		return errors.New("SetClipboardData failed")
	}

	return nil
}

// Rectangle is used to desribe monitor and window boundaries.
type Rectangle struct {
	// X is the left-most pixel.
	X int
	// Y is the top-most pixel.
	Y int
	// Width is the width in pixels.
	Width int
	// Height is the height in pixels.
	Height int
}

// Monitor is a single monitor connected to your computer.
type Monitor struct {
	// Rectangle is the outer boundary of the monitor, in virtual screen
	// coordinates. All monitors share this virtual coordinate system. In your
	// operating system settings you can freely move monitors around in this
	// coordinate system to represent the real world layout of your monitors.
	// For example, you might put two monitors side by side or on top of each
	// other.
	Rectangle
	// WorkArea is the monitor area that is not covered by the task bar.
	WorkArea Rectangle
	// Primary is true if this is the current default/primary monitor.
	Primary bool
}

// Windows returns a list of all currently active windows.
func Windows() ([]Window, error) {
	var windows []Window
	if !w32.EnumWindows(func(window w32.HWND) bool {
		windows = append(windows, windowHandleToWindow(window))
		return true
	}) {
		return nil, errors.New("EnumWindows failed")
	}
	return windows, nil
}

// PrimaryMonitor returns the current default/primary monitor.
func PrimaryMonitor() (Monitor, error) {
	m := w32.MonitorFromPoint(0, 0, w32.MONITOR_DEFAULTTOPRIMARY)
	if m == 0 {
		return Monitor{}, errors.New(
			"MonitorFromPoint with MONITOR_DEFAULTTOPRIMARY failed",
		)
	}
	return monitorHandleToMonitor(m)
}

// Monitors returns all monitors currently connected to the computer.
func Monitors() ([]Monitor, error) {
	var monitorHandles []w32.HMONITOR
	if !w32.EnumDisplayMonitors(
		0,
		nil,
		syscall.NewCallback(func(m w32.HMONITOR, _ w32.HDC, _ *w32.RECT, _ w32.LPARAM) uintptr {
			monitorHandles = append(monitorHandles, m)
			return 1
		}),
		0,
	) {
		return nil, errors.New("EnumDisplayMonitors failed")
	}

	monitors := make([]Monitor, len(monitorHandles))
	for i := range monitors {
		m, err := monitorHandleToMonitor(monitorHandles[i])
		if err != nil {
			return nil, err
		}
		monitors[i] = m
	}

	return monitors, nil
}

func monitorHandleToMonitor(monitor w32.HMONITOR) (Monitor, error) {
	var info w32.MONITORINFO
	if !w32.GetMonitorInfo(monitor, &info) {
		return Monitor{}, errors.New("GetMonitorInfo failed")
	}
	return Monitor{
		Rectangle: Rectangle{
			X:      int(info.RcMonitor.Left),
			Y:      int(info.RcMonitor.Top),
			Width:  int(info.RcMonitor.Width()),
			Height: int(info.RcMonitor.Height()),
		},
		WorkArea: Rectangle{
			X:      int(info.RcWork.Left),
			Y:      int(info.RcWork.Top),
			Width:  int(info.RcWork.Width()),
			Height: int(info.RcWork.Height()),
		},
		Primary: info.DwFlags&w32.MONITORINFOF_PRIMARY != 0,
	}, nil
}

// CaptureWindow returns a screen shot of the outer boundaries of the given
// window.
func CaptureWindow(w Window) (image.Image, error) {
	return CaptureScreenRect(w.Rectangle)
}

// CaptureWindowContent returns a screen shot of the inner boundaries of the
// given window.
func CaptureWindowContent(w Window) (image.Image, error) {
	return CaptureScreenRect(w.Content)
}

// CaptureMonitor returns a screen shot of the outer boundaries of the given
// monitor.
func CaptureMonitor(m Monitor) (image.Image, error) {
	return CaptureScreenRect(m.Rectangle)
}

// CaptureScreenRect is a wrapper for CaptureScreen. It allows you to pass a
// Monitor's WorkArea to this function instead of unwrapping the Rectangle
// yourself.
func CaptureScreenRect(r Rectangle) (image.Image, error) {
	return CaptureScreen(r.X, r.Y, r.Width, r.Height)
}

// CaptureScreen returns a screen shot of the given area. The area is given in
// virtual screen coordinates.
func CaptureScreen(x, y, width, height int) (image.Image, error) {
	screenDC := w32.GetDC(0)
	if screenDC == 0 {
		return nil, errors.New("GetDC failed")
	}
	defer w32.ReleaseDC(0, screenDC)

	memDC := w32.CreateCompatibleDC(screenDC)
	if memDC == 0 {
		panic("CreateCompatibleDC failed")
	}
	defer w32.DeleteDC(memDC)

	screenBitmap := w32.CreateCompatibleBitmap(screenDC, width, height)
	if screenBitmap == 0 {
		panic("CreateCompatibleDC failed")
	}
	defer w32.DeleteObject(w32.HGDIOBJ(screenBitmap))

	w32.SelectObject(memDC, w32.HGDIOBJ(screenBitmap))
	blitted := w32.BitBlt(
		memDC,
		0, 0, width, height,
		screenDC,
		x, y,
		w32.SRCCOPY,
	)
	if !blitted {
		panic("BitBlt failed")
	}

	format := w32.BITMAPINFOHEADER{
		BiSize:        uint32(binary.Size(w32.BITMAPINFOHEADER{})),
		BiWidth:       int32(width),
		BiHeight:      int32(-height),
		BiPlanes:      1,
		BiBitCount:    32,
		BiCompression: w32.BI_RGB,
	}

	byteCount := 4 * width * height
	memory := w32.GlobalAlloc(w32.GMEM_MOVEABLE, uint32(byteCount))
	defer w32.GlobalFree(memory)
	memoryPointer := w32.GlobalLock(memory)
	defer w32.GlobalUnlock(memory)

	if 0 == w32.GetDIBits(
		screenDC,
		screenBitmap,
		0,
		uint(height),
		memoryPointer,
		(*w32.BITMAPINFO)(unsafe.Pointer(&format)),
		w32.DIB_RGB_COLORS,
	) {
		panic("GetDIBits failed")
	}

	rawSlice := &reflect.SliceHeader{
		Data: uintptr(memoryPointer),
		Len:  byteCount,
		Cap:  byteCount,
	}
	raw := *(*[]byte)(unsafe.Pointer(rawSlice))
	pixels := make([]byte, 4*width*height)
	copy(pixels, raw)
	runtime.KeepAlive(raw)

	// Windows gives us BRGA, we want RGBA, so we swap 2 of the 4 bytes.
	for i := 0; i < len(pixels); i += 4 {
		pixels[i], pixels[i+2] = pixels[i+2], pixels[i]
	}

	return &image.RGBA{
		Pix:    pixels,
		Stride: 4 * width,
		Rect:   image.Rect(x, y, x+width, y+height),
	}, nil
}

// CaptureMonitors returns a screen shot of the outer hull of all the given
// monitors. Depending on your operating system settings this may include blank
// areas which will be transparent in the image. For example, if you have a 1200
// pixel high monitor next to a 1080 pixel high monitor, there will a 1200-1080
// = 120 pixel high area below the smaller monitor that is transparent.
func CaptureMonitors(monitors []Monitor) (image.Image, error) {
	if len(monitors) == 0 {
		return &image.RGBA{}, errors.New("now monitor given")
	}
	hullLeft := monitors[0].X
	hullTop := monitors[0].Y
	hullRight := monitors[0].X + monitors[0].Width
	hullBottom := monitors[0].Y + monitors[0].Height
	for _, m := range monitors {
		left := m.X
		top := m.Y
		right := m.X + m.Width
		bottom := m.Y + m.Height
		if left < hullLeft {
			hullLeft = left
		}
		if right > hullRight {
			hullRight = right
		}
		if top < hullTop {
			hullTop = top
		}
		if bottom > hullBottom {
			hullBottom = bottom
		}
	}
	r := Rectangle{
		X:      hullLeft,
		Y:      hullTop,
		Width:  hullRight - hullLeft,
		Height: hullBottom - hullTop,
	}
	return CaptureScreenRect(r)
}

// Key... constants are keys you can pass to TypeKey, PressKey and ReleaseKey.
const (
	KeyA                  = 'A'
	KeyB                  = 'B'
	KeyC                  = 'C'
	KeyD                  = 'D'
	KeyE                  = 'E'
	KeyF                  = 'F'
	KeyG                  = 'G'
	KeyH                  = 'H'
	KeyI                  = 'I'
	KeyJ                  = 'J'
	KeyK                  = 'K'
	KeyL                  = 'L'
	KeyM                  = 'M'
	KeyN                  = 'N'
	KeyO                  = 'O'
	KeyP                  = 'P'
	KeyQ                  = 'Q'
	KeyR                  = 'R'
	KeyS                  = 'S'
	KeyT                  = 'T'
	KeyU                  = 'U'
	KeyV                  = 'V'
	KeyW                  = 'W'
	KeyX                  = 'X'
	KeyY                  = 'Y'
	KeyZ                  = 'Z'
	Key0                  = '0'
	Key1                  = '1'
	Key2                  = '2'
	Key3                  = '3'
	Key4                  = '4'
	Key5                  = '5'
	Key6                  = '6'
	Key7                  = '7'
	Key8                  = '8'
	Key9                  = '9'
	KeyLeftButton         = w32.VK_LBUTTON
	KeyRightButton        = w32.VK_RBUTTON
	KeyMiddleButton       = w32.VK_MBUTTON
	KeyXButton1           = w32.VK_XBUTTON1
	KeyXButton2           = w32.VK_XBUTTON2
	KeyCancel             = w32.VK_CANCEL
	KeyBackspace          = w32.VK_BACK
	KeyTab                = w32.VK_TAB
	KeyClear              = w32.VK_CLEAR
	KeyEnter              = w32.VK_RETURN
	KeyShift              = w32.VK_SHIFT
	KeyControl            = w32.VK_CONTROL
	KeyAlt                = w32.VK_MENU
	KeyPause              = w32.VK_PAUSE
	KeyCapsLock           = w32.VK_CAPITAL
	KeyImeKana            = w32.VK_KANA
	KeyImeHangul          = w32.VK_HANGUL
	KeyImeOn              = w32.VK_IME_ON
	KeyImeJunja           = w32.VK_JUNJA
	KeyImeFinal           = w32.VK_FINAL
	KeyImeHanja           = w32.VK_HANJA
	KeyImeKanji           = w32.VK_KANJI
	KeyImeOff             = w32.VK_IME_OFF
	KeyEscape             = w32.VK_ESCAPE
	KeyImeConvert         = w32.VK_CONVERT
	KeyImeNonConvert      = w32.VK_NONCONVERT
	KeyImeAccept          = w32.VK_ACCEPT
	KeyImeModeChange      = w32.VK_MODECHANGE
	KeySpace              = w32.VK_SPACE
	KeyPageUp             = w32.VK_PRIOR
	KeyPageDown           = w32.VK_NEXT
	KeyEnd                = w32.VK_END
	KeyHome               = w32.VK_HOME
	KeyLeft               = w32.VK_LEFT
	KeyUp                 = w32.VK_UP
	KeyRight              = w32.VK_RIGHT
	KeyDown               = w32.VK_DOWN
	KeySelect             = w32.VK_SELECT
	KeyPrint              = w32.VK_PRINT
	KeyExecute            = w32.VK_EXECUTE
	KeyPrintScreen        = w32.VK_SNAPSHOT
	KeyInsert             = w32.VK_INSERT
	KeyDelete             = w32.VK_DELETE
	KeyHelp               = w32.VK_HELP
	KeyLeftWin            = w32.VK_LWIN
	KeyRightWin           = w32.VK_RWIN
	KeyApps               = w32.VK_APPS
	KeySleep              = w32.VK_SLEEP
	KeyNum0               = w32.VK_NUMPAD0
	KeyNum1               = w32.VK_NUMPAD1
	KeyNum2               = w32.VK_NUMPAD2
	KeyNum3               = w32.VK_NUMPAD3
	KeyNum4               = w32.VK_NUMPAD4
	KeyNum5               = w32.VK_NUMPAD5
	KeyNum6               = w32.VK_NUMPAD6
	KeyNum7               = w32.VK_NUMPAD7
	KeyNum8               = w32.VK_NUMPAD8
	KeyNum9               = w32.VK_NUMPAD9
	KeyMultiply           = w32.VK_MULTIPLY
	KeyPlus               = w32.VK_ADD
	KeySeparator          = w32.VK_SEPARATOR
	KeyMinus              = w32.VK_SUBTRACT
	KeyDecimal            = w32.VK_DECIMAL
	KeyDivide             = w32.VK_DIVIDE
	KeyF1                 = w32.VK_F1
	KeyF2                 = w32.VK_F2
	KeyF3                 = w32.VK_F3
	KeyF4                 = w32.VK_F4
	KeyF5                 = w32.VK_F5
	KeyF6                 = w32.VK_F6
	KeyF7                 = w32.VK_F7
	KeyF8                 = w32.VK_F8
	KeyF9                 = w32.VK_F9
	KeyF10                = w32.VK_F10
	KeyF11                = w32.VK_F11
	KeyF12                = w32.VK_F12
	KeyF13                = w32.VK_F13
	KeyF14                = w32.VK_F14
	KeyF15                = w32.VK_F15
	KeyF16                = w32.VK_F16
	KeyF17                = w32.VK_F17
	KeyF18                = w32.VK_F18
	KeyF19                = w32.VK_F19
	KeyF20                = w32.VK_F20
	KeyF21                = w32.VK_F21
	KeyF22                = w32.VK_F22
	KeyF23                = w32.VK_F23
	KeyF24                = w32.VK_F24
	KeyNumLock            = w32.VK_NUMLOCK
	KeyScrollLock         = w32.VK_SCROLL
	KeyOemNecEqual        = w32.VK_OEM_NEC_EQUAL
	KeyOemFjJisho         = w32.VK_OEM_FJ_JISHO
	KeyOemFjMasshou       = w32.VK_OEM_FJ_MASSHOU
	KeyOemFjTouroku       = w32.VK_OEM_FJ_TOUROKU
	KeyOemFjLoya          = w32.VK_OEM_FJ_LOYA
	KeyOemFjRoya          = w32.VK_OEM_FJ_ROYA
	KeyLeftShift          = w32.VK_LSHIFT
	KeyRightShift         = w32.VK_RSHIFT
	KeyLeftControl        = w32.VK_LCONTROL
	KeyRightControl       = w32.VK_RCONTROL
	KeyLeftAlt            = w32.VK_LMENU
	KeyRightAlt           = w32.VK_RMENU
	KeyBrowserBack        = w32.VK_BROWSER_BACK
	KeyBrowserForward     = w32.VK_BROWSER_FORWARD
	KeyBrowserRefresh     = w32.VK_BROWSER_REFRESH
	KeyBrowserStop        = w32.VK_BROWSER_STOP
	KeyBrowserSearch      = w32.VK_BROWSER_SEARCH
	KeyBrowserFavorites   = w32.VK_BROWSER_FAVORITES
	KeyBrowserHome        = w32.VK_BROWSER_HOME
	KeyVolumeMute         = w32.VK_VOLUME_MUTE
	KeyVolumeDown         = w32.VK_VOLUME_DOWN
	KeyVolumeUp           = w32.VK_VOLUME_UP
	KeyMediaNextTrack     = w32.VK_MEDIA_NEXT_TRACK
	KeyMediaPreviousTrack = w32.VK_MEDIA_PREV_TRACK
	KeyMediaStop          = w32.VK_MEDIA_STOP
	KeyMediaPlayPause     = w32.VK_MEDIA_PLAY_PAUSE
	KeyLaunchMail         = w32.VK_LAUNCH_MAIL
	KeyLaunchMediaSelect  = w32.VK_LAUNCH_MEDIA_SELECT
	KeyLaunchApp1         = w32.VK_LAUNCH_APP1
	KeyLaunchApp2         = w32.VK_LAUNCH_APP2
	KeyOemPlus            = w32.VK_OEM_PLUS
	KeyOemComma           = w32.VK_OEM_COMMA
	KeyOemMinus           = w32.VK_OEM_MINUS
	KeyOemPeriod          = w32.VK_OEM_PERIOD
	KeyOem1               = w32.VK_OEM_1
	KeyOem2               = w32.VK_OEM_2
	KeyOem3               = w32.VK_OEM_3
	KeyOem4               = w32.VK_OEM_4
	KeyOem5               = w32.VK_OEM_5
	KeyOem6               = w32.VK_OEM_6
	KeyOem7               = w32.VK_OEM_7
	KeyOem8               = w32.VK_OEM_8
	KeyOemAx              = w32.VK_OEM_AX
	KeyOem102             = w32.VK_OEM_102
	KeyIcoHelp            = w32.VK_ICO_HELP
	KeyIco00              = w32.VK_ICO_00
	KeyImeProcessKey      = w32.VK_PROCESSKEY
	KeyIcoClear           = w32.VK_ICO_CLEAR
	KeyUnicodePacket      = w32.VK_PACKET
	KeyOemReset           = w32.VK_OEM_RESET
	KeyOemJump            = w32.VK_OEM_JUMP
	KeyOemPa1             = w32.VK_OEM_PA1
	KeyOemPa2             = w32.VK_OEM_PA2
	KeyOemPa3             = w32.VK_OEM_PA3
	KeyOemWsControl       = w32.VK_OEM_WSCTRL
	KeyOemCuSel           = w32.VK_OEM_CUSEL
	KeyOemAttn            = w32.VK_OEM_ATTN
	KeyOemFinish          = w32.VK_OEM_FINISH
	KeyOemCopy            = w32.VK_OEM_COPY
	KeyOemAuto            = w32.VK_OEM_AUTO
	KeyOemEnlw            = w32.VK_OEM_ENLW
	KeyOemNBackTab        = w32.VK_OEM_BACKTAB
	KeyAttn               = w32.VK_ATTN
	KeyCrSel              = w32.VK_CRSEL
	KeyExSel              = w32.VK_EXSEL
	KeyErEof              = w32.VK_EREOF
	KeyPlay               = w32.VK_PLAY
	KeyZoom               = w32.VK_ZOOM
	KeyNoName             = w32.VK_NONAME
	KeyPa1                = w32.VK_PA1
	KeyOemClear           = w32.VK_OEM_CLEAR
)
