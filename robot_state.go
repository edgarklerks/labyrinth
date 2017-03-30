package main
import "log"
import "lang1/ast"
import "math/rand"
import "strconv"

const LANGUAGE_1 = 0
const LANGUAGE_2 = 1

type RobotHook interface {
	r_update(rstate RobotState)
}

type RobotState struct {
        x int
        y int
        mx int
        my int
        dir int
	language_type int
	program1 ast.Program1
	code_pointer_register int
	repeat_register int
	maze *Maze
}

func initRobotStateProgram1(x,y,mx,my int, dir int, program1 ast.Program1, maze *Maze) RobotState {
	return RobotState {
		x: x,
		y: y,
		mx: mx,
		my: my,
		dir: dir,
		program1: program1,
		language_type: LANGUAGE_1,
		code_pointer_register: 0,
		repeat_register: 0,
		maze: maze,
	}

}
/** Inject a new program into the robot state **/
func (rstate RobotState) callProgram1(program1 ast.Program1) RobotState {
	return RobotState {
		x:rstate.x,
		y:rstate.y,
		mx:rstate.mx,
		my:rstate.my,
		dir:rstate.dir,
		language_type: rstate.language_type,
		program1: program1,
		code_pointer_register: 0,
		repeat_register: rstate.repeat_register,
		maze: rstate.maze,
	}
}
/** Merge an old and a new state after a call. This needs to preserve the old registers, but take the new positional parameters **/
func (rstate RobotState) restoreState(new_rstate RobotState) RobotState {
	return RobotState {
		x:new_rstate.x,
		y:new_rstate.y,
		mx:new_rstate.mx,
		my:new_rstate.my,
		dir:new_rstate.dir,
		language_type:rstate.language_type,
		program1:rstate.program1,
		code_pointer_register: rstate.code_pointer_register,
		repeat_register: rstate.repeat_register,
		maze: rstate.maze,
	}

}

func (rstate RobotState) collide(dir int) bool {
	maze := rstate.maze
	x := rstate.x
	y := rstate.y
	switch {
	    case dir == GO_LEFT:
		    return maze.grid[x][y] & RightWall == RightWall
	case dir == GO_RIGHT:
		    return maze.grid[x][y] & LeftWall == LeftWall
	case dir == GO_UP:
		    return maze.grid[x][y] & DownWall == DownWall
	case dir == GO_DOWN:
		    return maze.grid[x][y] & UpWall == UpWall

	}

	return false
}

func (rstate RobotState) left() RobotState {
	next_state := RobotState {
		x: nMod(rstate.x - 1, rstate.mx),
		y: rstate.y,
		mx: rstate.mx,
		my: rstate.my,
		dir: GO_LEFT,
		language_type: rstate.language_type,
		program1: rstate.program1,
		code_pointer_register: rstate.code_pointer_register,
		repeat_register: rstate.repeat_register,
		maze: rstate.maze,
	}

	if next_state.collide(GO_LEFT) {
		return rstate
	}
	return next_state
}

func (rstate RobotState) right() RobotState {
	next_state := RobotState {
		x: nMod(rstate.x + 1,rstate.mx),
		y: rstate.y,
		mx: rstate.mx,
		my: rstate.my,
		dir: GO_RIGHT,
		language_type: rstate.language_type,
		program1: rstate.program1,
		code_pointer_register: rstate.code_pointer_register,
		repeat_register: rstate.repeat_register,
		maze: rstate.maze,

	}
	if next_state.collide(GO_RIGHT) {
		return rstate
	}
	return next_state
}

func (rstate RobotState) up() RobotState {
	next_state := RobotState {
		x: rstate.x,
		y: nMod(rstate.y - 1,rstate.my),
		mx: rstate.mx,
		my: rstate.my,
		dir: rstate.dir,
		language_type: rstate.language_type,
		program1: rstate.program1,
		code_pointer_register: rstate.code_pointer_register,
		repeat_register: rstate.repeat_register,
		maze: rstate.maze,

	}
	if next_state.collide(GO_UP) {
		return rstate
	}
	return next_state
}

func (rstate RobotState) down() RobotState {
	next_state := RobotState {
		x: rstate.x,
		y: nMod(rstate.y + 1,rstate.my),
		mx: rstate.mx,
		my: rstate.my,
		dir: rstate.dir,
		language_type: rstate.language_type,
		program1: rstate.program1,
		code_pointer_register: rstate.code_pointer_register,
		repeat_register: rstate.repeat_register,
		maze: rstate.maze,

	}
	if next_state.collide(GO_DOWN) {
		return rstate
	}
	return next_state
}


/** Switch for choosing the correct language, the hook is used to update the screen **/
func (rstate RobotState) run(rhook RobotHook) RobotState {
	switch {
	case rstate.language_type == LANGUAGE_1:
		return rstate.run_language1(rhook)
	}

	return rstate


}

func (rstate RobotState) dec_repeat() RobotState {
	return RobotState {
		x: rstate.x,
		y: rstate.y,
		mx: rstate.mx,
		my: rstate.my,
		dir: rstate.dir,
		language_type: rstate.language_type,
		program1: rstate.program1,
		code_pointer_register: rstate.code_pointer_register,
		repeat_register: rstate.repeat_register - 1,
		maze: rstate.maze,
	}

}

func (rstate RobotState) reset() RobotState {
	rstate.code_pointer_register = 0
	rstate.repeat_register = 0
	return rstate
}

func (rstate RobotState) next_instruction() RobotState {
	return RobotState {
		x: rstate.x,
		y: rstate.y,
		mx: rstate.mx,
		my: rstate.my,
		dir: rstate.dir,
		language_type: rstate.language_type,
		program1: rstate.program1,
		code_pointer_register: rstate.code_pointer_register + 1,
		repeat_register: rstate.repeat_register,
		maze: rstate.maze,
	}
}
/** This is basicly a virtual machine. It will halt, when the code_pointer reaches the end of the program **/
func (rstate RobotState) run_language1(rhook RobotHook) RobotState {
	/** Build symbol table **/
	var symtab map[string]*ast.Program1
	symtab = make(map[string]*ast.Program1, 0)

	for index,ins := range rstate.program1 {
		switch expr := (ins).(type) {
		case ast.Proc:
			name := expr.Name
			prog := ast.ToProgram(expr.Expr)
			prog_clone := prog
			symtab[name] = &prog_clone
			rstate.program1[index] = ast.Nop{Label: "Replaced by symbuilder",}
		}

	}


	for rstate.code_pointer_register < len(rstate.program1) {

		/** Lookup the next instruction **/
		ins := rstate.program1[rstate.code_pointer_register]
		log.Printf("Registers: code_pointer_r: %d, repeat_r: %d", rstate.code_pointer_register, rstate.repeat_register)

		switch expr := (ins).(type) {
		case ast.Call:
			name := expr.Name
			prog := symtab[name]
			log.Printf("Entering subprogram from Call")
			rstate_new := rstate.callProgram1(*prog)
			rstate_new = rstate_new.run_language1(rhook)
			rstate = rstate.restoreState(rstate_new)
			rstate = rstate.next_instruction()
		case ast.If:
			/** TODO: implement tests **/
			// dir: left, right, up, down
			test_result := false
			test_type := expr.Test.Type
			test_dir := expr.Test.Dir
			log.Printf("executing: %s %s", test_type, test_dir)
			switch {
			case test_type == "wall":
				test_result = rstate.test_wall(test_dir)
			case test_type == "open":
				test_result = rstate.test_open(test_dir)
			case test_type == "robot":
				test_result = rstate.test_robot(test_dir)
			case test_type == "rand":
				test_result = rstate.random_choice(test_dir)

			}

			var new_prog ast.Program1
			if test_result == true {
				new_prog = ast.ToProgram(expr.Ok)
			} else {
				new_prog = ast.ToProgram(expr.Nok)
			}
			rstate_new := rstate.callProgram1(new_prog)
			rstate_new = rstate_new.run_language1(rhook)
			rstate = rstate.restoreState(rstate_new)
			rstate = rstate.next_instruction()
		case ast.Nop:
			log.Printf("Not doing anything")
			rstate = rstate.next_instruction()
		case ast.Left:
			log.Printf("Running instruction: Left")
			rstate = rstate.left()
			rstate = rstate.next_instruction()
		case ast.Right:
			log.Printf("Running instruction: Right")
			rstate = rstate.right()
			rstate = rstate.next_instruction()
		case ast.Up:
			log.Printf("Running instruction: Up")
			rstate = rstate.up()
			rstate = rstate.next_instruction()
		case ast.Down:
			log.Printf("Running instruction: Down")
			rstate = rstate.down()
			rstate = rstate.next_instruction()
		case ast.Group:
			new_program := expr.Prog
			log.Printf("Entering subprogram from Group")
			rstate_new := rstate.callProgram1(new_program)
			rstate_new = rstate_new.run_language1(rhook)
			rstate = rstate.restoreState(rstate_new)
			log.Printf("Leaving subprogram from Group")
			rstate = rstate.next_instruction()
		case ast.Repeat:
			new_ins := expr.Expr
			repeats := expr.N
			log.Printf("Entering repeat instruction, repeating %d times", repeats)
			/** Set the repeat register **/
			new_program := ast.ToProgram(new_ins)
			log.Printf("Entering subprogram")
			/** Entering a subprogram, we need to load the program
			    This looks stupid, but we can use the same logic for function calls
			**/
			/** Let the program know how many repeats it will run **/
			rstate_new := rstate.callProgram1(new_program)

			for i := 0; i < repeats; i++  {
				rstate_new.repeat_register = repeats - i
				log.Printf("In repeat, Registers: code_pointer_r: %d, repeat_r: %d", rstate_new.code_pointer_register, rstate_new.repeat_register)

				rstate_new = rstate_new.run_language1(rhook)
				rstate_new = rstate_new.reset()
			}

			/** Restore the registers **/
			rstate = rstate.restoreState(rstate_new)
			log.Printf("Leaving subprogram")
			rstate = rstate.next_instruction()

		}
	    rhook.r_update(rstate)
	}
  rstate.reset()
	log.Printf("Program done")

	return rstate
}
/** Various tests for the robot to detect a wall and other features of the maze **/
func (rstate RobotState) test_wall(dir string) bool {
	maze := rstate.maze
	posx := nMod(rstate.x, rstate.mx)
	posy := nMod(rstate.y,rstate.my)

	wall := maze.grid[posx][posy]

	switch {
	case dir == "left" && (wall & LeftWall) == LeftWall:
		log.Printf("Detected left wall")
		return true
	case dir == "right" && (wall & RightWall) == RightWall:
		log.Printf("Detected right wall")
		return true
	case dir == "up" && (wall & UpWall) == UpWall:
		log.Printf("Detected up wall")
		return true
        case dir == "down" && (wall & DownWall) == DownWall:
		log.Printf("Detected down wall")
		return true
	}
	return false
}

func (rstate RobotState) test_open(dir string) bool {
	maze := rstate.maze
	posx := nMod(rstate.x, rstate.mx)
	posy := nMod(rstate.y,rstate.my)

	wall := maze.grid[posx][posy]
	log.Printf("Received dir: %s", dir)

	switch {
	case dir == "left" && (wall & LeftWall) != LeftWall:
		log.Printf("Detected left passage")
		return true
	case dir == "right" && (wall & RightWall) != RightWall:
		log.Printf("Detected right passage")
		return true
	case dir == "up" && (wall & UpWall) !=  UpWall:
		log.Printf("Detected up passage")
		return true
        case dir == "down" && (wall & DownWall) != DownWall:
		log.Printf("Detected down passage")
		return true
	}
	return false

}

func (rstate RobotState) random_choice(times string) bool {
  q,err := strconv.ParseInt(times, 10, 32)
  if err != nil{
          panic(err)
  }
	c := rand.Intn(int(q))

  log.Printf("random_choice 1 / %d: %d == 0?",q, c)
	return c == 0
}

func (rstate RobotState) test_robot(dir string) bool {
	log.Printf("test_robot has not been implemented, reverting to first branch")
	return true
}
