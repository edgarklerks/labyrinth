package main

import "github.com/veandco/go-sdl2/sdl"
import "github.com/veandco/go-sdl2/sdl_image"

func (game *Game) renderRobot(rstate RobotState) {
        var flip sdl.RendererFlip 
        switch {
          case GO_LEFT == rstate.dir:
                flip = sdl.FLIP_NONE
          case GO_RIGHT == rstate.dir:
                flip = sdl.FLIP_HORIZONTAL
        }
        robot_rect := sdl.Rect { 
                X: int32(10 * nMod(rstate.x,rstate.mx)),
                Y: int32(10 * nMod(rstate.y,rstate.my)),
                W: 9,
                H: 9,
        }

        texture, err := img.LoadTexture(game.render,"img/pacman.png")
        if err != nil {
                panic(err)
        }
        defer texture.Destroy()

        game.render.CopyEx(texture, nil, &robot_rect,0,nil,flip)
}

func (game *Game) renderCell(cell Cell, i int, j int){
	sizex := game.sizex
	sizey := game.sizey
	delta := 5
			/**
A cell is 10 by 10 big, thus to find the middle of our cell we need to multiply by 5 and then
move 5 to right and 5 down
 **/
	centerx := i * 10 + delta
	centery := j * 10 + delta
	for cell != NoWall {
		switch {
		case cell & LeftWall == LeftWall:
			cell = cell ^ LeftWall
			x := nMod(centerx - delta, sizex)
			y1 := nMod(centery - delta, sizey)
			y2 := nMod(centery + delta, sizey)
			game.CreateVLine(x,y1,y2)

		case cell & RightWall == RightWall:
			cell = cell ^ RightWall
			x := nMod(centerx + delta, sizex)
			y1 := nMod(centery - delta, sizey)
			y2 := nMod(centery + delta, sizey)
			game.CreateVLine(x,y1,y2)
		case cell & UpWall == UpWall:
			cell = cell ^ UpWall
			x1 := nMod(centerx - delta, sizex)
			x2 := nMod(centerx + delta, sizex)
			y := nMod(centery - delta, sizey)
			game.CreateHLine(x1,x2,y)
		case cell & DownWall == DownWall:
			cell = cell ^ DownWall
			x1 := nMod(centerx - delta, sizex)
			x2 := nMod(centerx + delta, sizex)
			y := nMod(centery + delta, sizey)
			game.CreateHLine(x1,x2,y)
		}

	}

}
func Min(x, y int64) int64 {
    if x < y {
        return x
    }
    return y
}

func Max(x, y int64) int64 {
    if x > y {
        return x
    }
    return y
}
func (game *Game) renderMaze(maze *Maze) {

	/** This whole scaling shit goes to hell, when the maze gets to big
I need to think it over
On the other hand it could be an artefact of scaling.
Biggest maze is 64 x 48
Perhaps I should just accept it
**/
	screen_sizex := game.sizex
	screen_sizey := game.sizey
	maze_sizex := maze.sizex
	maze_sizey := maze.sizey
	scalex := float32((screen_sizex - 2))/float32(maze_sizex*10)
	scaley := float32((screen_sizey - 2))/float32(maze_sizey*10)
	game.render.SetScale(scalex,scaley)
	game.CreateBorder()

	/** Paint the maze **/
	for i := 0; i < maze_sizex; i++ {
		for j := 0; j < maze_sizey; j++{
			cell := maze.grid[i][j]
			game.renderCell(cell, i, j)

		}
	}


}
