package main
import "log"
import "lang1/ast"
import "errors"

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
	prev_dir int
	symtab map[string]*ast.Program1
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
    symtab: make(map[string]*ast.Program1, 0),
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
    symtab: rstate.symtab,
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
    symtab: rstate.symtab,
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
    symtab: rstate.symtab,
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
    symtab: rstate.symtab,

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
    symtab: rstate.symtab,

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
    symtab: rstate.symtab,

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
    symtab: rstate.symtab,
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
    symtab: rstate.symtab,
	}
}
/** This is basicly a virtual machine. It will halt, when the code_pointer reaches the end of the program **/
func (rstate RobotState) run_language1(rhook RobotHook) RobotState {
	/** Build symbol table **/

	for _,ins := range rstate.program1 {
		switch expr := (ins).(type) {
		case ast.Proc:
			name := expr.Name
			prog := ast.ToProgram(expr.Expr)
      log.Printf("Converted %s to program %s", expr.Expr.ToString(), name)
			rstate.symtab[name] = &prog
		}

	}


	for rstate.code_pointer_register < len(rstate.program1) {

		/** Lookup the next instruction **/
		ins := rstate.program1[rstate.code_pointer_register]
		log.Printf("Registers: code_pointer_r: %d, repeat_r: %d", rstate.code_pointer_register, rstate.repeat_register)

		switch expr := (ins).(type) {
    case ast.Proc:
            rstate = rstate.next_instruction()
		case ast.Call:
			name := expr.Name
			prog := rstate.symtab[name]
      if prog == nil {
              log.Printf("No such sub program %s", name)
              panic(errors.New("No such sub program"))
      }
			log.Printf("Entering subprogram from Call")
			rstate_new := rstate.callProgram1(*prog)
			rstate_new = rstate_new.run_language1(rhook)
			log.Printf("Exiting subprogram from Call")
			rstate = rstate.restoreState(rstate_new)
		log.Printf("Registers: code_pointer_r: %d, repeat_r: %d", rstate.code_pointer_register, rstate.repeat_register)
			rstate = rstate.next_instruction()
		case ast.If:
			test_result := rstate.evaluatePredicate(expr.Test)

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
