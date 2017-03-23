package main
import "container/heap"
import "math/rand"
import "errors"
import "log"

type Position struct {
	x int
	y int
	sizex int
	sizey int
}

type CarveState struct {
	backtracks *PriorityQueue
	current Position
	maze *Maze
}

const GO_LEFT =  0
const GO_RIGHT = 1
const GO_UP = 2
const GO_DOWN = 3


func (pos Position) doStep(direction int) Position{
	switch {
	case direction == GO_LEFT:
		return Position {
			x: nMod(pos.x - 1, pos.sizex),
			y: pos.y,
			sizex: pos.sizex,
			sizey: pos.sizey,
		}
	case direction == GO_RIGHT:
		return Position {
			x: nMod(pos.x + 1, pos.sizex),
			y: pos.y,
			sizex: pos.sizex,
			sizey: pos.sizey,
		}

	case direction == GO_UP:
		return Position {
			x: pos.x,
			y: nMod(pos.y - 1, pos.sizey),
			sizex: pos.sizex,
			sizey: pos.sizey,
		}
	case direction == GO_DOWN:
		return Position {
			x: pos.x,
			y: nMod(pos.y + 1, pos.sizey),
			sizex: pos.sizex,
			sizey: pos.sizey,
		}


	}
	panic(errors.New("This shouldn't happen :("))
}

func (cs CarveState) nextStep() (dir int, pos Position,err error) {
	var dirs []int
	dirs = rand.Perm(4)
	log.Printf("Directions: %s", dirs)

	for i := 0; i < 4; i++ {
		j := dirs[i]
		pos := cs.current.doStep(j)
		cell := cs.maze.grid[pos.x][pos.y]
		if !cell.MazePart() {
			return j, pos, nil
		}
	}
	log.Printf("No more steps")
	return 0, Position{}, errors.New("No more steps")
}
/** Add a point to the tracking list **/
func (cs CarveState) track(pos Position) {
	log.Printf("Tracking (%d,%d)\n", pos.x,pos.y)
	cs.backtracks.Push(&Item{
          x: pos.x,
          y: pos.y,
          priority: rand.Int(),
  })
}

/** Backtracker, finds a previous eligible point for carving, it uses a heap to randomize the points **/
func (cs CarveState) backtrack() (pos Position, err error){
        pq := cs.backtracks 
        if pq == nil {
                panic(errors.New("Priority queue is nil!"))
        }
        for (*pq).Len() > 0{
            item := heap.Pop(pq).(*Item)
            pos := Position{
                    x: item.x,
                    y: item.y,
                    sizex: cs.maze.sizex,
                    sizey: cs.maze.sizey,

            }
            cell := cs.maze.grid[pos.x][pos.y]
            if(cell.IsPassage() || cell.IsDeadEnd() || cell.IsJunction()){
              return pos, nil
            }
        }
        return Position{}, errors.New("No points to backtrack")
}
/** Update the random seed **/
func updateSeed(seed int64){
	rand.Seed(seed)
}

/** Check certain invariants in the maze **/
func checkCarving(maze *Maze) bool {
	sizex := maze.sizex
	sizey := maze.sizey
	ok := true
	for i := 0; i < sizex; i++ {
		for j := 0; j < sizey; j++ {
			cell := maze.grid[i][j]
			ok := ok && (cell.IsDeadEnd() || cell.IsPassage() || cell.IsJunction())
			if !ok {
				log.Printf("Found wrong cell: %s at (%d,%d)\n", cell, i,j)
			}
		}
	}
	return ok



}

/** Carve a maze with a recursive backtracker. It is reasonable simple and creats aesthetically pleasant mazes 
   I made a small change regarding the backtracking, the points are stored in a priority queue with randomized priorities
**/
func carveMaze(sizex int, sizey int, maze_number int64) *Maze{
	updateSeed(maze_number)
	maze := generateFullMaze(sizex,sizey)
	startPosition := Position{
		x: 0,
		y: 0,
		sizex: sizex,
		sizey: sizey,
	}
	carveState := CarveState{
		backtracks: NewQueue(32),
		current:startPosition,
		maze: maze,
	}


	carveState.maze = maze

	for {

	    dir,next,error := carveState.nextStep()

	    if(error != nil){
		    log.Printf("needs to backtrack\n")
		    backpos,err := carveState.backtrack()
		    /** No further points to backtrack **/
		    if err != nil {
			    log.Print("No points available\n")
			    /** We are done, no points to backtrack **/
			    return carveState.maze
		    }
		    log.Printf("Choosing next point: (%d,%d)\n", backpos.x,backpos.y)

		    carveState.current = backpos
	    } else {
		log.Printf("Step: (%d,%d) -> (%d,%d)\n", carveState.current.x, carveState.current.y, next.x,next.y)
		/** Here we did a step, we need to carve and update the current position **/
		    old_pos := carveState.current
		    carveState.track(old_pos)
		    carveState.current = next

		    /** carve the shit out of it. If we move left, the the previous cell needs to be carved on the left side,
                        However the next cell needs to be carved to the right side:
                       *--**--*
                       |  <-  |
                       *--**--*
                       <- carving direction
                       * wall or removed wall
		    **/
		    new_cell := carveState.maze.grid[next.x][next.y]
		    old_cell := carveState.maze.grid[old_pos.x][old_pos.y]
		    switch {
			case dir == GO_LEFT:
			    carveState.maze.grid[next.x][next.y] = new_cell ^ RightWall
			    carveState.maze.grid[old_pos.x][old_pos.y] = old_cell ^ LeftWall
			case dir == GO_RIGHT:
			    carveState.maze.grid[next.x][next.y] = new_cell ^ LeftWall
			    carveState.maze.grid[old_pos.x][old_pos.y] = old_cell ^ RightWall
			case dir == GO_UP:
			    carveState.maze.grid[next.x][next.y] = new_cell ^ DownWall
			    carveState.maze.grid[old_pos.x][old_pos.y] = old_cell ^ UpWall
			case dir == GO_DOWN:
			    carveState.maze.grid[next.x][next.y] = new_cell ^ UpWall
			    carveState.maze.grid[old_pos.x][old_pos.y] = old_cell ^ DownWall
		    }
	    }
	}
	return maze

}
