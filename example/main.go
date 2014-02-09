// © 2014 the grlutil Authors under the WTFPL. See AUTHORS for the list of authors.

package main

import (
        "fmt"
        "math/rand"

        rlutil "github.com/chrissexton/grlutil"
)

const (
        FLOOR = iota
        WALL
        COIN
        STAIRS_DOWN
        TORCH
        MAPSIZE = 15
)

var x, y int
var coins, moves, torch, level int = 0, 0, 30, 1
var lvl [MAPSIZE][MAPSIZE]int

/// Generates the dungeon map
func gen(seed int) {
        rand.Seed(int64(seed))
        for j := 0; j < MAPSIZE; j++ {
                for i := 0; i < MAPSIZE; i++ {
                        if i == 0 || i == MAPSIZE-1 || j == 0 || j == MAPSIZE-1 || rand.Int()%10 == 0 {
                                lvl[i][j] = 1
                        } else if rand.Int()%20 == 0 {
                                lvl[i][j] = COIN
                        } else if rand.Int()%100 == 0 {
                                lvl[i][j] = TORCH
                        } else {
                                lvl[i][j] = 0
                        }
                }
        }
        var randcoord int
        x = 1 + rand.Int()%MAPSIZE - 2
        y = 1 + rand.Int()%MAPSIZE - 2
        lvl[randcoord][randcoord] = STAIRS_DOWN
}

/// Draws the screen
func draw() {
        rlutil.Cls()
        rlutil.Locate(1, MAPSIZE+1)
        rlutil.SetColor(rlutil.YELLOW)
        fmt.Printf("Coins: %d\n", coins)
        rlutil.SetColor(rlutil.RED)
        fmt.Printf("Torch: %d\n", torch)
        rlutil.SetColor(rlutil.MAGENTA)
        fmt.Printf("Moves: %d\n", moves)
        rlutil.SetColor(rlutil.GREEN)
        fmt.Printf("Level: %d\n", level)
        rlutil.Locate(1, 1)
        abs := func(val int) int {
                if val < 0 {
                        return -1 * val
                }
                return val
        }
        min := func(x, y int) int {
                if x < y {
                        return x
                }
                return y
        }
        for j := 0; j < MAPSIZE; j++ {
                for i := 0; i < MAPSIZE; i++ {
                        if abs(x-i)+abs(y-j) > min(10, torch/2) {
                                fmt.Printf(" ")
                        } else if lvl[i][j] == 0 {
                                rlutil.SetColor(rlutil.BLUE)
                                fmt.Printf(".")
                        } else if lvl[i][j]&WALL != 0 {
                                rlutil.SetColor(rlutil.CYAN)
                                fmt.Printf("#")
                        } else if lvl[i][j]&COIN != 0 {
                                rlutil.SetColor(rlutil.YELLOW)
                                fmt.Printf("o")
                        } else if lvl[i][j]&STAIRS_DOWN != 0 {
                                rlutil.SetColor(rlutil.GREEN)
                                fmt.Printf("<")
                        } else if lvl[i][j]&TORCH != 0 {
                                rlutil.SetColor(rlutil.RED)
                                fmt.Printf("f")
                        }
                }
                fmt.Printf("\n")
        }
        rlutil.Locate(x+1, y+1)
        rlutil.SetColor(rlutil.WHITE)
        fmt.Printf("@")
}

/// Main loop and input handling
func main() {
        rlutil.HideCursor()
        gen(level)
        rlutil.SetColor(2)
        fmt.Printf("Welcome! Use WASD to move.\n")
        rlutil.SetColor(6)
        fmt.Printf("Hit any key to start.\n")
        rlutil.AnyKey()
        draw()
        for true {
                // Input
                if rlutil.KbHit() {
                        var k = rlutil.GetCh()

                        oldx, oldy := x, y
                        if k == 'a' {
                                x--
                                moves++
                        } else if k == 'd' {
                                x++
                                moves++
                        } else if k == 'w' {
                                y--
                                moves++
                        } else if k == 's' {
                                y++
                                moves++
                        } else if k == 27 {
                                break
                        }
                        // Collisions
                        if lvl[x][y]&WALL != 0 {
                                x = oldx
                                y = oldy
                        } else if lvl[x][y]&COIN != 0 {
                                coins++
                                lvl[x][y] ^= COIN
                        } else if lvl[x][y]&TORCH != 0 {
                                torch += 20
                                lvl[x][y] ^= TORCH
                        } else if lvl[x][y]&STAIRS_DOWN != 0 {
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
        rlutil.ShowCursor()
}
