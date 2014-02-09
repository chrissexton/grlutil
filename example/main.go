// © 2014 the grl Authors under the WTFPL. See AUTHORS for the list of authors.

package main

import (
	"flag"
	"fmt"
	"math/rand"

	rl "github.com/chrissexton/grlutil"
)

const (
	FLOOR = iota
	WALL
	COIN
	STAIRS_DOWN
	TORCH

	MAPSIZE = 15
)

type keymap struct {
	up    rune
	down  rune
	left  rune
	right rune
	quit  rune
}

var qwerty keymap = keymap{
	up:    'w',
	down:  's',
	left:  'a',
	right: 'd',
	quit:  'q',
}

var dvorak keymap = keymap{
	up:    ',',
	down:  'o',
	left:  'a',
	right: 'e',
	quit:  'q',
}

var x, y int

var coins, moves, torch, level int

var lvl [MAPSIZE][MAPSIZE]int

/// Generates the dungeon map
func gen(seed int) {
	rand.Seed(int64(seed))

	for j := 0; j < MAPSIZE; j++ {
		for i := 0; i < MAPSIZE; i++ {
			if i == 0 || i == MAPSIZE-1 || j == 0 || j == MAPSIZE-1 || rand.Int()%10 == 0 {
				lvl[i][j] = WALL
			} else if rand.Int()%20 == 0 {
				lvl[i][j] = COIN
			} else if rand.Int()%100 == 0 {
				lvl[i][j] = TORCH
			} else {
				lvl[i][j] = FLOOR
			}
		}
	}

	x = 1 + rand.Int()%(MAPSIZE-2)
	y = 1 + rand.Int()%(MAPSIZE-2)

	x1 := 1 + rand.Int()%(MAPSIZE-2)
	y1 := 1 + rand.Int()%(MAPSIZE-2)
	lvl[x1][y1] = STAIRS_DOWN
}

/// Draws the screen
func draw() {
	rl.Cls()
	rl.Locate(1, MAPSIZE+1)
	rl.SetColor(rl.YELLOW)
	fmt.Printf("Coins: %d\n", coins)
	rl.SetColor(rl.RED)
	fmt.Printf("Torch: %d\n", torch)
	rl.SetColor(rl.MAGENTA)
	fmt.Printf("Moves: %d\n", moves)
	rl.SetColor(rl.GREEN)
	fmt.Printf("Level: %d\n", level)
	rl.Locate(1, 1)
	for j := 0; j < MAPSIZE; j++ {
		for i := 0; i < MAPSIZE; i++ {
			if rl.Abs(x-i)+rl.Abs(y-j) > rl.Min(10, torch/2) {
				fmt.Printf(" ")
			} else if lvl[i][j] == FLOOR {
				rl.SetColor(rl.BLUE)
				fmt.Printf(".")
			} else if lvl[i][j] == WALL {
				rl.SetColor(rl.CYAN)
				fmt.Printf("#")
			} else if lvl[i][j] == COIN {
				rl.SetColor(rl.YELLOW)
				fmt.Printf("o")
			} else if lvl[i][j] == STAIRS_DOWN {
				rl.SetColor(rl.GREEN)
				fmt.Printf("<")
			} else if lvl[i][j] == TORCH {
				rl.SetColor(rl.RED)
				fmt.Printf("f")
			}
		}
		fmt.Printf("\n")
	}
	rl.Locate(x+1, y+1)
	rl.SetColor(rl.WHITE)
	fmt.Printf("@")
}

/// Main loop and input handling
func main() {
	useDvorak := flag.Bool("dvorak", false, "Use dvorak keybindings (,oae)")
	flag.IntVar(&torch, "torches", 30, "[cheat] Set initial number of torches")
	flag.IntVar(&coins, "coins", 0, "[cheat] Set initial number of coins")
	flag.IntVar(&level, "level", 1, "[cheat] Set initial level")
	flag.Parse()

	keys := qwerty
	if *useDvorak {
		keys = dvorak
	}

	rl.HideCursor()
	gen(level)
	rl.SetColor(rl.GREEN)
	if *useDvorak {
		fmt.Printf("Welcome! Use ,oae (the physical WSAD keys) to move.\n")
	} else {
		fmt.Printf("Welcome! Use WASD to move.\n")
	}
	rl.SetColor(rl.CYAN)
	fmt.Printf("Hit any key to start.\n")
	rl.AnyKey()
	draw()
	for true {
		// Input
		if rl.KbHit() {
			var k = rl.GetCh()

			oldx, oldy := x, y
			if k == keys.left {
				x--
				moves++
			} else if k == keys.right {
				x++
				moves++
			} else if k == keys.up {
				y--
				moves++
			} else if k == keys.down {
				y++
				moves++
			} else if k == keys.quit {
				break
			} else {
				continue
			}
			// Collisions
			if lvl[x][y] == WALL {
				x = oldx
				y = oldy
			} else if lvl[x][y] == COIN {
				coins++
				lvl[x][y] ^= COIN
			} else if lvl[x][y] == TORCH {
				torch += 20
				lvl[x][y] ^= TORCH
			} else if lvl[x][y] == STAIRS_DOWN {
				level++
				gen(level)
			}
			// Drawing
			draw()
			// Die
			torch--
			if torch <= 0 {
				break
			}
		}
	}
	rl.ShowCursor()
}
