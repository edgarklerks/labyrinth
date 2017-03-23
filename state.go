package main

import "github.com/veandco/go-sdl2/sdl"
import "github.com/veandco/go-sdl2/sdl_gfx"
import "github.com/hideo55/go-popcount"

type Game struct {
	window *sdl.Window
	render *sdl.Renderer
	sizex int
	sizey int
}

type Cell uint64

type Maze struct {
	sizex int
	sizey int
	grid [][]Cell
}

const NoWall = 0
const LeftWall = 1
const RightWall = 2
const UpWall = 4
const DownWall = 8
const FullWall = DownWall + UpWall + LeftWall + RightWall


/**
Tests if this is a dead end
**/
func (cell Cell) IsDeadEnd() bool {
	return popcount.Count(uint64(cell)) == 3
}

/**
Tests if this is a passage
**/
func (cell Cell) IsPassage() bool {
	return popcount.Count(uint64(cell)) == 2
}

/**
Tests if this is a T junction
**/
func (cell Cell) IsJunction() bool {
	return popcount.Count(uint64(cell)) == 1
}
/**
Tests if this is a square
**/

func (cell Cell) IsSquare() bool {
	return popcount.Count(uint64(cell)) == 0
}

/**
Is this a part of the maze
**/
func (cell Cell) MazePart() bool {
	return uint64(cell) != FullWall
}
/** Crappy implementation of the signum
function (I don't care to much about it right now, sorry
**/
func signum(a int) int {
	switch {
	case a < 0:
		return -1
	case a >= 0:
		return +1
	}
	return +1
}
/**
The modulus of go returns also negative numbers,
which btw is a good thing. However we need one, which
returns positive numbers solely
**/
func nMod(p int, m int) int {
	modulus := p % m
	if signum(modulus) == -1 {
		return modulus + m
	} else {
		return modulus
	}
}
/**
Create a maze grid consisting of all walls
**/
func generateFullMaze(sizex int, sizey int) *Maze {
	grid := make([][]Cell, sizex)

	for i := 0; i < sizex; i++ {
		grid[i] = make([]Cell, sizey)
		for j := 0; j < sizey; j++ {
			grid[i][j] = FullWall
		    }
	}

	return &Maze {
		sizex: sizex,
		sizey: sizey,
		grid:grid,
	}
}

func InitGame(sizex int, sizey int) *Game {
	sdl.Init(sdl.INIT_EVERYTHING)
	window, err := sdl.CreateWindow("Labyrinth", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, sizex, sizey, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	render,err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	game := Game {
		window: window,
		render: render,
		sizex: sizex,
		sizey: sizey,
	}
	game.render.Clear()
	return &game
}


func (game *Game) CreateBorder(){
	gfx.BoxRGBA(game.render,
		0,0,800,600,0,0,0,255)

}

func (game *Game) CreateHLine(x1 int,x2 int, y int){
	// boxSize := 1
	game.render.SetDrawColor(0,255,0,255)
	game.render.DrawLine(x1,y,x2,y)
	// gfx.BoxRGBA(game.render,x1,y,x2,y+boxSize,0,255,0,255)
}

func (game *Game) CreateVLine(x int, y1 int, y2 int){
	// boxSize := 1
	game.render.SetDrawColor(0,255,0,255)
	game.render.DrawLine(x,y1,x,y2)
	// gfx.BoxRGBA(game.render,x,y1,x+boxSize,y2,0,255,0,255)
}

func (game *Game) Clear(){
        game.render.Clear()
}

func (game *Game) Update(){
	game.render.Present()
}

func (game *Game) Destroy(){
	game.window.Destroy()
}
