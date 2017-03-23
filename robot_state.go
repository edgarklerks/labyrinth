package main

type RobotState struct {
        x int
        y int
        mx int
        my int
        dir int
}

func (rstate RobotState) collide(dir int, maze *Maze) bool {
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

func (rstate RobotState) left(maze *Maze) RobotState {
	next_state := RobotState {
		x: nMod(rstate.x - 1, rstate.mx),
		y: rstate.y,
		mx: rstate.mx,
		my: rstate.my,
		dir: GO_LEFT,
	}

	if next_state.collide(GO_LEFT, maze) {
		return rstate
	}
	return next_state
}

func (rstate RobotState) right(maze *Maze) RobotState {
	next_state := RobotState {
		x: nMod(rstate.x + 1,rstate.mx),
		y: rstate.y,
		mx: rstate.mx,
		my: rstate.my,
		dir: GO_RIGHT,
	}
	if next_state.collide(GO_RIGHT,maze) {
		return rstate
	}
	return next_state
}

func (rstate RobotState) up(maze *Maze) RobotState {
	next_state := RobotState {
		x: rstate.x,
		y: nMod(rstate.y - 1,rstate.my),
		mx: rstate.mx,
		my: rstate.my,
		dir: rstate.dir,
	}
	if next_state.collide(GO_UP,maze) {
		return rstate
	}
	return next_state
}

func (rstate RobotState) down(maze *Maze) RobotState {
	next_state := RobotState {
		x: rstate.x,
		y: nMod(rstate.y + 1,rstate.my),
		mx: rstate.mx,
		my: rstate.my,
		dir: rstate.dir,
	}
	if next_state.collide(GO_DOWN,maze) {
		return rstate
	}
	return next_state
}
