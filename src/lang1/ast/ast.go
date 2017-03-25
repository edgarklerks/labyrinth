package ast
import "strconv"
import "strings"
import "log"
import "fmt"
import "errors"
/**
Structure of the first language:
program := <command>*
command := repeat <number> <command> | left | right | up | down
number := [0-9]+
**/
/** This simulates a sumtype **/
type Language1 interface {
        IsLang1()
	ToString() string
}




/**
The robot actions users should delay inside functions. Perhaps we need a bit more information from sdl to do this reliable.
I would say, let it do a step every 100 milliseconds or so
**/

type Program1 []Language1

func EmptyProgram() (Program1,error) {

	log.Printf("EmptyProgram empty program created")
	var emp []Language1
	emp = make([]Language1,0)
	return emp,nil
}


func NewProgram(ys interface{}, xs interface{}) (Program1,error) {
	log.Printf("NewProgram created")
	y,ok := ys.(Language1)
	if !ok{
		return nil, errors.New("Expected expression")
	}
	log.Printf("NewProgram found %s", y.ToString())
	qs,ok := xs.(Program1)
	if ok {
	    qs = append(qs,y)
		log.Printf("NewProgram fine")
		PrintTree(qs)
		log.Printf("NewProgram end")
	    return qs,nil
	} else {
		if qs == nil {
			var ps []Language1
			ps = make([]Language1, 1)
			ps[0] = y
			log.Printf("NewProgram Empty program returned...")
			PrintTree(ps)
			return ps,nil
		}
	}
	return nil, errors.New("Expected program")
}

func ToProgram(q Language1) Program1 {
    new_program := make(Program1, 1)
	new_program[0] = q
	return new_program
}

type Repeat struct {
	N int
	Expr Language1
}
func (e Repeat) ToString() string {
	return fmt.Sprintf("repeat %d %s", e.N, e.Expr.ToString())
}
func (e Repeat) IsLang1() {}
func NewRepeat(num interface{}, exp interface{}) (Repeat,error) {
	log.Printf("New repeat created")
        n, err := strconv.Atoi(num.(string))
        if err != nil {
                return Repeat{}, err
        }
        return Repeat {
                N:n,
                Expr: exp.(Language1),
        },nil

}

type Left struct {}

func NewLeft() (Left, error) {
	log.Printf("New left created")
	return Left{}, nil
}
func (e Left) ToString() string{
	return "left"
}

func (e Left) IsLang1() {}
func NewRight() (Right, error) {
	log.Printf("New right created")
	return Right{},nil
}
type Right struct {}
func (e Right) ToString() string { return "right" }
func (e Right) IsLang1() {}
type Up struct {}
func (e Up) ToString() string { return "up" }
func NewUp() (Up,error) {
	log.Printf("New up created")
	return Up{},nil
}
func (e Up) IsLang1() {}
type Down struct {}
func (e Down) ToString() string { return "down" }
func NewDown() (Down,error){
	log.Printf("New down created")
	return Down{},nil
}
func (e Down) IsLang1() {}
type Group struct {
	Prog Program1
}
func (e Group) ToString() string {
	out := make([]string, len(e.Prog))
	for i, exp := range e.Prog {
		out[i] = exp.ToString()
	}
	return fmt.Sprintf("{ %s }", strings.Join(out, "\n"))
}
func NewGroup(exp interface{}) (Group, error) {
	log.Printf("New group created")
	prog,ok := exp.(Program1)
	if ok {
	    return Group{Prog:prog},nil
	} else {
		lang1,ok := exp.(Language1)
		if ok {
			var prog []Language1
			prog = make([]Language1,1)
			prog[0] = lang1
			return Group{Prog:prog},nil
		} else  {
			return Group{},errors.New("Need to be a subprogram or a expression")
		}
	}
}
func (e Group) IsLang1() {}


func PrintTree(program1 Program1) {
	log.Printf("Outputting program from PrintTree")
	for _,e := range program1 {
		switch expr := e.(type) {
		case Left:
			fmt.Printf("left\n")
		case Right:
			fmt.Printf("right\n")
		case Up:
			fmt.Printf("up\n")
		case Down:
			fmt.Printf("down\n")
		case Repeat:
			fmt.Printf("repeat %d ", expr.N)
			PrintTree(ToProgram(expr.Expr))
		case Group:
			fmt.Printf("{")
			PrintTree(expr.Prog)
			fmt.Printf("}")
		}
	}
}

func TestProgram() Program1 {
	var program []Language1
	program = make([]Language1,10)
	program[0] = Repeat{N: 5, Expr: Left{},}
	program[1] = Up{}
	program[2] = Right{}
	program[3] = Right{}
	program[4] = Down{}
	program[5] = Right{}
	program[6] = Left{}
	program[7] = Down{}
	program[8] = Down{}
	program[9] = Repeat{N: 7, Expr: Repeat{N:3 , Expr:Group{
		Prog: program[1:4],
	},}}
	return program
}
