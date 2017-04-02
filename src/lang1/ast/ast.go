package ast
import "strconv"
import "strings"
import "fmt"
import "errors"
import "lang1/token"
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


func EmptyProgram() Program1 {

	var emp []Language1
	emp = make([]Language1,0)
	return emp
}


func NewProgram(ys interface{}, xs interface{}) (Program1,error) {
	y,ok := ys.(Language1)
	if !ok{
		return nil, errors.New("Expected expression")
	}
	qs,ok := xs.(Program1)
	if ok {
			var ps []Language1
			ps = make([]Language1, 1)
      ps[0] = y 
	    qs = append(ps, qs...)
	    return qs,nil
	} else {
		if qs == nil {
			var ps []Language1
			ps = make([]Language1, 1)
			ps[0] = y
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
        n, err := strconv.Atoi(string(num.(*token.Token).Lit))
        if err != nil {
                return Repeat{}, err
        }
        return Repeat {
                N:n,
                Expr: exp.(Language1),
        },nil

}

type Return struct {}

func (e Return) IsLang1(){}
func (e Return) ToString() string { return "todo" }

func NewReturn() (Return, error) {
        return Return{},nil
}

type Proc struct {
	Name string
	Expr Language1

}
func (e Proc) IsLang1(){}
func (e Proc) ToString() string {
	return fmt.Sprintf("proc %s %s", e.Name, e.ToString())
}

type Predicate interface {
        IsPred()
}

type Test struct {
	Type string
	Dir string
}
func (e Test) IsLang1(){}
func (e Test) IsPred(){}
func (e Test) ToString() string {
	return fmt.Sprintf("%s %s",e.Type,e.Dir)
}

func NewTest(tpe interface{}, dir interface{}) (Test, error) {
	tpe_r := string(tpe.(*token.Token).Lit)
	dir_r := string(dir.(*token.Token).Lit)
	return Test{Type:tpe_r,Dir:dir_r,},nil
}

type TestGroup struct {
    Group Predicate 
}

func (e TestGroup) IsLang1(){}
func (e TestGroup) IsPred(){}
func (e TestGroup) ToString() string { return "todo"}


func NewTestGroup(test interface{}) (TestGroup, error){
        return TestGroup{
                Group: getPredicate(test),
        },nil

}


type TestAnd struct {
        Test1 Predicate
        Test2 Predicate
}

func (e TestAnd) IsLang1(){}
func (e TestAnd) IsPred(){}
func (e TestAnd) ToString() string { return "todo"}

type TestOr struct {
        Test1 Predicate
        Test2 Predicate
}

func (e TestOr) IsLang1(){}
func (e TestOr) IsPred(){}
func (e TestOr) ToString() string { return "todo"}

type TestNot struct {
        Test1 Predicate
}

func (e TestNot) IsLang1(){}
func (e TestNot) IsPred(){}
func (e TestNot) ToString() string { return "todo"}


func getPredicate(test interface{}) Predicate {
        switch expr := test.(type) {
        case Test:
                return expr
        case TestAnd:
                return expr
        }
        return Test{}
}


func NewTestNot(test1 interface{}) (TestNot, error){
         test_1 := getPredicate(test1)
         return TestNot{
                 Test1: test_1,
         },nil
}

func NewTestOr(test1 interface{}, test2 interface{}) (TestOr, error) {
         test_1 := getPredicate(test1)
        test_2 := getPredicate(test2)
        return TestOr{
                Test1: test_1,
                Test2: test_2,
        }, nil
}

func NewTestAnd(test1 interface{}, test2 interface{}) (TestAnd, error) {
        test_1 := getPredicate(test1)
        test_2 := getPredicate(test2)
        return TestAnd{
                Test1: test_1,
                Test2: test_2,
        }, nil
}


type If struct {
	Test Predicate
	Ok Language1
	Nok Language1
}

func (e If) ToString() string {
	return "todo"
}

func NewIf(test interface{}, ok interface{}, nok interface{}) (If, error) {
	if nok == nil {
		return If{
			Test: getPredicate(test),
			Ok: ok.(Language1),
			Nok: Group{Prog:make([]Language1,0),},
		},nil
	} else {
		return If{
			Test: getPredicate(test),
			Ok: ok.(Language1),
			Nok: nok.(Language1),
	},nil
	}
}

func (e If) IsLang1(){}
func (e If) IsString() string {
	return "todo" // fmt.Sprintf("if %s { %s } else { %s }", e.Test.ToString(),e.Ok.ToString(), e.Nok.ToString())
}

type Call struct {
	Name string
}
func (e Call) IsLang1(){}
func (e Call) ToString() string {
	return fmt.Sprintf("call %s", e.Name)
}


func NewCall(label interface{}) (Call, error){
	r_label := string(label.(*token.Token).Lit)
	return Call{
		Name: r_label,
	},nil
}
type Nop struct{Label string}

func NewNop(label interface{}) (Nop, error) {
	r_label := string(label.(*token.Token).Lit)
	return Nop{
		Label:r_label,

	},nil
}

func (e Nop) IsLang1() {}
func (e Nop) ToString() string {
	return fmt.Sprintf("nop %s", e.Label)
}

func NewProc(label interface{}, expr interface{}) (Proc,error) {
	r_label := string(label.(*token.Token).Lit)
	return Proc{
		Name: r_label,
		Expr: expr.(Language1),
	},nil
}

type Left struct {}

func NewLeft() (Left, error) {
	return Left{}, nil
}
func (e Left) ToString() string{
	return "left"
}

func (e Left) IsLang1() {}
func NewRight() (Right, error) {
	return Right{},nil
}
type Right struct {}
func (e Right) ToString() string { return "right" }
func (e Right) IsLang1() {}
type Up struct {}
func (e Up) ToString() string { return "up" }
func NewUp() (Up,error) {
	return Up{},nil
}
func (e Up) IsLang1() {}
type Down struct {}
func (e Down) ToString() string { return "down" }
func NewDown() (Down,error){
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
	program = make([]Language1,12)
	program[0] = Repeat{N: 5, Expr: Left{},}
	program[1] = Up{}
	program[2] = Right{}
	program[3] = Right{}
	program[4] = Down{}
	program[5] = Right{}
	program[6] = Left{}
	program[7] = Down{}
  program[8] = Repeat{N: 4, Expr: Group{ Prog: program[3:6], }, }
  program[9] = Repeat{N: 2, Expr: Repeat{N:3 , Expr:Group{Prog: program[1:4],},},}
  program[10] = Repeat{N: 4, Expr: Group{ Prog: program[3:6], }, }
	program[11] = Repeat{N: 2, Expr: Repeat{N:3 , Expr:Group{
		Prog: program[1:4],
	},}}
	return program
}
