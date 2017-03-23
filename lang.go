package main
import "fmt"
import "github.com/veandco/go-sdl2/sdl"
/**
Structure of the first language:
program := <command>*
command := repeat <number> <command> | left | right | up | down
number := [0-9]+
**/
/** This simulates a sumtype **/
type Language1 interface {
	isLang1()
}
/**
The robot actions users should delay inside functions. Perhaps we need a bit more information from sdl to do this reliable.
I would say, let it do a step every 100 milliseconds or so
**/
type ActionsLanguage1 interface {
	goLeft(rstate RobotState, maze *Maze) RobotState
	goRight(rstate RobotState, maze *Maze) RobotState
	goUp(rstate RobotState, maze *Maze) RobotState
	goDown(rstate RobotState, maze *Maze) RobotState
 }

type Program1 []Language1

type Repeat struct {
	n uint
	expr Language1
}
func (e Repeat) isLang1() {}

type Left struct {}
func (e Left) isLang1() {}
type Right struct {}
func (e Right) isLang1() {}
type Up struct {}
func (e Up) isLang1() {}
type Down struct {}
func (e Down) isLang1() {}

func eval1(e Language1, actions ActionsLanguage1, rstate RobotState, maze *Maze) RobotState {
	switch expr := (e).(type) {
		case Left:
		    fmt.Printf("I am going left")
		    return actions.goLeft(rstate,maze)
	    case Right:
		    fmt.Printf("I am going right")
		    return actions.goRight(rstate,maze)
            case Up:
		    fmt.Printf("I am going up")
		    return actions.goUp(rstate,maze)
            case Down:
		    fmt.Printf("I am going down")
		    return actions.goDown(rstate,maze)
            case Repeat:
		    n := expr.n
		    inner := expr.expr
		    var i uint
		    for i = 0; i < n; i++ {
			    rstate = eval1(inner, actions, rstate,maze)
		    }
		    return rstate
	    }
	return rstate
}

type LangAction struct {}

func (la LangAction) goLeft(rs RobotState, maze *Maze) RobotState {
	sdl.Delay(100)
	return rs.left(maze)
}

func (la LangAction) goRight(rs RobotState, maze *Maze) RobotState {
	sdl.Delay(100)
	return rs.right(maze)
}

func (la LangAction) goDown(rs RobotState, maze *Maze) RobotState {
	sdl.Delay(100)
	return rs.down(maze)
}

func (la LangAction) goUp(rs RobotState, maze *Maze) RobotState {
	sdl.Delay(100)
	return rs.up(maze)
}




func testProgram() Program1 {
	var program []Language1
	program = make([]Language1,5)
	program[0] = Left{}
	program[1] = Up{}
	program[2] = Right{}
	program[3] = Right{}
	program[4] = Down{}
	return program
}
