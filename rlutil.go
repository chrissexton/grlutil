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

func GotoXY(x, y int) {
	Locate(x, y)
}

const (
	// Comparable colors
	BLACK = iota
	BLUE
	GREEN
	CYAN
	RED
	MAGENTA
	BROWN
	GREY
	DARKGREY
	LIGHTBLUE
	LIGHTGREEN
	LIGHTCYAN
	LIGHTRED
	LIGHTMAGENTA
	YELLOW
	WHITE

	// Ugly ANSI sequences
	ANSI_CLS          = "\033[2J"
	ANSI_BLACK        = "\033[22;30m"
	ANSI_RED          = "\033[22;31m"
	ANSI_GREEN        = "\033[22;32m"
	ANSI_BROWN        = "\033[22;33m"
	ANSI_BLUE         = "\033[22;34m"
	ANSI_MAGENTA      = "\033[22;35m"
	ANSI_CYAN         = "\033[22;36m"
	ANSI_GREY         = "\033[22;37m"
	ANSI_DARKGREY     = "\033[01;30m"
	ANSI_LIGHTRED     = "\033[01;31m"
	ANSI_LIGHTGREEN   = "\033[01;32m"
	ANSI_YELLOW       = "\033[01;33m"
	ANSI_LIGHTBLUE    = "\033[01;34m"
	ANSI_LIGHTMAGENTA = "\033[01;35m"
	ANSI_LIGHTCYAN    = "\033[01;36m"
	ANSI_WHITE        = "\033[01;37m"

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

// Get a numeric keycode from the system
func getKey() rune {
	k := GetCh()
	switch k {
	case 0:
		{
			kk := GetCh()
			switch kk {
			case 71:
				return KEY_NUMPAD7
			case 72:
				return KEY_NUMPAD8
			case 73:
				return KEY_NUMPAD9
			case 75:
				return KEY_NUMPAD4
			case 77:
				return KEY_NUMPAD6
			case 79:
				return KEY_NUMPAD1
			case 80:
				return KEY_NUMPAD4
			case 81:
				return KEY_NUMPAD3
			case 82:
				return KEY_NUMPAD0
			case 83:
				return KEY_NUMDEL
			default:
				return kk - 59 + KEY_F1 // Function keys
			}
		}
	case 224:
		{
			kk := GetCh()
			switch kk {
			case 71:
				return KEY_HOME
			case 72:
				return KEY_UP
			case 73:
				return KEY_PGUP
			case 75:
				return KEY_LEFT
			case 77:
				return KEY_RIGHT
			case 79:
				return KEY_END
			case 80:
				return KEY_DOWN
			case 81:
				return KEY_PGDOWN
			case 82:
				return KEY_INSERT
			case 83:
				return KEY_DELETE
			default:
				return kk - 123 + KEY_F1 // Function keys
			}
		}
	case 13:
		return KEY_ENTER
	case 155: // single-character CSI (no idea what that means)
	case 27:
		{
			// Process ANSI escape sequences
			cnt := -1
			if cnt >= 3 && GetCh() == '[' {
				k = GetCh()
				switch k {
				case 'A':
					return KEY_UP
				case 'B':
					return KEY_DOWN
				case 'C':
					return KEY_RIGHT
				case 'D':
					return KEY_LEFT
				}
			} else {
				return KEY_ESCAPE
			}
		}
	default:
		return k
	}
	return -1
}

// Get a character imediately or fail
func GetChNonBlocking() rune {
	if KbHit() {
		return GetCh()
	} else {
		return 0
	}
}

// Convert numeric color enums to ANSI ugliness
func GetANSIColor(c int) string {
	switch c {
	case 0:
		return ANSI_BLACK
	case 1:
		return ANSI_BLUE // non-ANSI
	case 2:
		return ANSI_GREEN
	case 3:
		return ANSI_CYAN // non-ANSI
	case 4:
		return ANSI_RED // non-ANSI
	case 5:
		return ANSI_MAGENTA
	case 6:
		return ANSI_BROWN
	case 7:
		return ANSI_GREY
	case 8:
		return ANSI_DARKGREY
	case 9:
		return ANSI_LIGHTBLUE // non-ANSI
	case 10:
		return ANSI_LIGHTGREEN
	case 11:
		return ANSI_LIGHTCYAN // non-ANSI;
	case 12:
		return ANSI_LIGHTRED // non-ANSI;
	case 13:
		return ANSI_LIGHTMAGENTA
	case 14:
		return ANSI_YELLOW // non-ANSI
	case 15:
		return ANSI_WHITE
	default:
		return ""
	}
}

// Switch printing colors
func SetColor(c int) {
	fmt.Print(GetANSIColor(c))
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
