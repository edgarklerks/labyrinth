package main

import "lang1/ast"
import "log"
import "strconv"
import "math/rand"

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

func (rstate RobotState) test_dir(dir string) bool {
log.Printf("Direction in test_dir: %s", dir)
	if rstate.prev_dir == GO_UNKNOWN {
		log.Printf("Prev_dir == unknown")
	    return false 
	}

	switch {
	case dir == "left" && rstate.prev_dir == GO_LEFT:
		log.Printf("Prev_dir == left")
		return true

	case dir == "right" && rstate.prev_dir == GO_RIGHT:
		log.Printf("Prev_dir == right")
		return true
	case dir == "up" && rstate.prev_dir == GO_UP:
		log.Printf("Prev_dir == up")
		return true
	case dir == "down" && rstate.prev_dir == GO_DOWN:
		log.Printf("Prev_dir == down")
		return true
	}
	return false
}

func (rstate RobotState) test_robot(dir string) bool {
	log.Printf("test_robot has not been implemented, reverting to first branch")
	return true
}

func (rstate RobotState) evaluatePredicate(predicate ast.Predicate) bool {

        switch expr := predicate.(type) {
        case ast.Test:
                test_result := false
                test_type := expr.Type
                test_dir := expr.Dir
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
                case test_type == "prev":
                        test_result = rstate.test_dir(test_dir)

                }
                return test_result
        case ast.TestGroup:
                test1 := expr.Group 
                return rstate.evaluatePredicate(test1)
        case ast.TestAnd:
                test1 := expr.Test1
                test2 := expr.Test2
                return rstate.evaluatePredicate(test1) && rstate.evaluatePredicate(test2)
        case ast.TestOr:
                test1 := expr.Test1
                test2 := expr.Test2
                return rstate.evaluatePredicate(test1) || rstate.evaluatePredicate(test2)
        case ast.TestNot:
                test1 := expr.Test1
                return !rstate.evaluatePredicate(test1)
        }

	return false
}
