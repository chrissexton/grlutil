// © 2014 the grlutil Authors under the WTFPL. See AUTHORS for the list of authors.

package rlutil

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

/*
#include "rlutil.h"
*/
import "C"

// Get character without waiting for Return to be pressed.
func GetCh() rune {
	return rune(C.getch())
}

/// Determines if keyboard has been hit.
func KbHit() bool {
	return C.kbhit() == 0
}

const (
	// Ugly ANSI sequences
	CLS          = "\033[2J"
	BLACK        = "\033[22;30m"
	RED          = "\033[22;31m"
	GREEN        = "\033[22;32m"
	BROWN        = "\033[22;33m"
	BLUE         = "\033[22;34m"
	MAGENTA      = "\033[22;35m"
	CYAN         = "\033[22;36m"
	GREY         = "\033[22;37m"
	DARKGREY     = "\033[01;30m"
	LIGHTRED     = "\033[01;31m"
	LIGHTGREEN   = "\033[01;32m"
	YELLOW       = "\033[01;33m"
	LIGHTBLUE    = "\033[01;34m"
	LIGHTMAGENTA = "\033[01;35m"
	LIGHTCYAN    = "\033[01;36m"
	WHITE        = "\033[01;37m"

	// Keycodes
	KEY_ESCAPE = 0
	KEY_ENTER  = 1
	KEY_SPACE  = 32

	KEY_INSERT = 2
	KEY_HOME   = 3
	KEY_PGUP   = 4
	KEY_DELETE = 5
	KEY_END    = 6
	KEY_PGDOWN = 7

	KEY_UP    = 14
	KEY_DOWN  = 15
	KEY_LEFT  = 16
	KEY_RIGHT = 17

	KEY_F1  = 18
	KEY_F2  = 19
	KEY_F3  = 20
	KEY_F4  = 21
	KEY_F5  = 22
	KEY_F6  = 23
	KEY_F7  = 24
	KEY_F8  = 25
	KEY_F9  = 26
	KEY_F10 = 27
	KEY_F11 = 28
	KEY_F12 = 29

	KEY_NUMDEL  = 30
	KEY_NUMPAD0 = 31
	KEY_NUMPAD1 = 127
	KEY_NUMPAD2 = 128
	KEY_NUMPAD3 = 129
	KEY_NUMPAD4 = 130
	KEY_NUMPAD5 = 131
	KEY_NUMPAD6 = 132
	KEY_NUMPAD7 = 133
	KEY_NUMPAD8 = 134
	KEY_NUMPAD9 = 135
)

// Get a character imediately or fail
func GetChNonBlocking() rune {
	if KbHit() {
		return GetCh()
	} else {
		return 0
	}
}

// Switch printing colors
func SetColor(c string) {
	fmt.Print(c)
}

// Clear the screen
func Cls() {
	fmt.Print("\033[2J\033[H")
}

// Move the cursor to (x,y) location
func Locate(x, y int) {
	buf := fmt.Sprintf("\033[%d;%df", y, x)
	fmt.Print(buf)
}

// Hide the cursor
func HideCursor() {
	fmt.Print("\033[?25l")
}

// Show the cursor
func ShowCursor() {
	fmt.Print("\033[?25h")
}

// Getsize returns the number of rows (lines) and cols (positions
// in each line) in terminal t.
func GetSize() (rows, cols int, err error) {
	var ws winsize
	err = windowrect(&ws, os.Stdin.Fd())
	return int(ws.ws_row), int(ws.ws_col), err
}

type winsize struct {
	ws_row    uint16
	ws_col    uint16
	ws_xpixel uint16
	ws_ypixel uint16
}

func windowrect(ws *winsize, fd uintptr) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

// Get number of columns in the terminal
func TermCols() int {
	_, cols, err := GetSize()
	if err != nil {
		panic(err)
	}
	return cols
}

// Get number of rows in the terminal
func TermRows() int {
	rows, _, err := GetSize()
	if err != nil {
		panic(err)
	}
	return rows
}

// Where's the any key!? This function knows.
func AnyKey() {
	GetCh()
}

// Integer minimum
func Min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

// Integer maximum
func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

// Integer absolute value
func Abs(val int) int {
	if val < 0 {
		return -1 * val
	}
	return val
}
