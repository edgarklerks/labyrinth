package main
import "lang1/ast"
// import "errors"
import "github.com/veandco/go-sdl2/sdl"
import "fmt"
// import "log"
import "os"
import "strconv"
//import "github.com/veandco/go-sdl2/sdl_mixer"
import "github.com/veandco/go-sdl2/sdl_image"
//import "github.com/veandco/go-sdl2/sdl_ttf"
// import "lang1/lexer"
// import "lang1/parser"
// import "io/ioutil"


func main(){
	img.Init(img.INIT_PNG)
	game := InitGame(800,600)
	var input sdl.Event

	/**
	prg_txt,err := ioutil.ReadFile("test.prog")
	if err != nil {
		panic(err)
	}
	s := lexer.NewLexer(prg_txt)
	p := parser.NewParser()
	program, err := p.Parse(s)
	if err != nil {
		fmt.Printf("Failed parse")
		panic(err)
	} else {
		log.Printf("PRINTING TREE")
		for _,pr :=  range program.([]ast.Program1) {
			ast.PrintTree(pr)
		}
	}
	panic(errors.New("STOP"))
**/

	maze_number_arg := os.Args[1]
	maze_number, err := strconv.ParseInt(maze_number_arg,10,64)
	if err != nil {
		panic("Need a number to generate the maze")
	}


	maze := carveMaze(16,12, maze_number)
	game.maze = maze

	rstate := initRobotStateProgram1(0,0,16,12,GO_RIGHT, ast.TestProgram(),maze)


	if(!checkCarving(maze)){
		fmt.Printf("Carving didn't went well\n");
	}

	/** Main loop **/

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
					rstate = rstate.right()
				case t.Keysym.Sym == sdl.GetKeyFromScancode(sdl.SCANCODE_LEFT):
					rstate = rstate.left()
				case t.Keysym.Sym == sdl.GetKeyFromScancode(sdl.SCANCODE_UP):
					rstate = rstate.up()
				case t.Keysym.Sym == sdl.GetKeyFromScancode(sdl.SCANCODE_DOWN):
					rstate = rstate.down()

				case t.Keysym.Sym == sdl.GetKeyFromScancode(sdl.SCANCODE_SPACE):
					rstate.reset()
					rstate = rstate.run_language1(game)
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
