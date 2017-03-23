package main

import "github.com/veandco/go-sdl2/sdl"
import "fmt"
import "os"
import "strconv"
import "github.com/hideo55/go-popcount"
//import "github.com/veandco/go-sdl2/sdl_mixer"
import "github.com/veandco/go-sdl2/sdl_image"
//import "github.com/veandco/go-sdl2/sdl_ttf"



func main(){
  img.Init(img.INIT_PNG)
	game := InitGame(800,600)
  var input sdl.Event

	maze_number_arg := os.Args[1]
	maze_number, err := strconv.ParseInt(maze_number_arg,10,64)
	if err != nil {
		panic("Need a number to generate the maze")
	}


	maze := carveMaze(16,12, maze_number)

  rstate := RobotState {
          x: 0,
          y: 0,
          mx: 16,
          my: 12,
          dir: GO_RIGHT,
  }

	fmt.Printf("Full maze: %s\n", maze.grid)
	fmt.Printf("Blub: %d\n", nMod(-2,  3))

	fmt.Printf("Pop count testing:\n")
	for i:=0; i < 15; i++ {
		fmt.Printf("%d: %d\n", uint64(i), popcount.Count(uint64(i)))
	}
	if(!checkCarving(maze)){
		fmt.Printf("Carving didn't went well\n");
	} else {
		fmt.Printf("Carving went fine\n")
	}
  quit := false
  for !quit {
	  for input = sdl.PollEvent(); input != nil; input = sdl.PollEvent() {
		  switch t := input.(type) {
		  case *sdl.QuitEvent:
			  fmt.Printf("Quit event\n")
			  quit = true
		  case *sdl.KeyDownEvent:
			  switch {
			  case t.Keysym.Sym == 'q':
				  fmt.Printf("Quit event\n")
				  quit = true
			  case t.Keysym.Sym == sdl.GetKeyFromScancode(sdl.SCANCODE_RIGHT):
				  rstate = rstate.right(maze)
			  case t.Keysym.Sym == sdl.GetKeyFromScancode(sdl.SCANCODE_LEFT):
				  rstate = rstate.left(maze)
			  case t.Keysym.Sym == sdl.GetKeyFromScancode(sdl.SCANCODE_UP):
				  rstate = rstate.up(maze)
			  case t.Keysym.Sym == sdl.GetKeyFromScancode(sdl.SCANCODE_DOWN):
				  rstate = rstate.down(maze)
			  default:
				  fmt.Printf("Quit event: %d\n", t.Keysym.Sym)

			  }

		  }
	  }
	  game.Clear()
	  game.renderMaze(maze)
	  game.renderRobot(rstate)
	  game.Update()
	  sdl.Delay(10)
  }


	game.Destroy()
}
