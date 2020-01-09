package xin

// Instructions
const (
	Nop = iota

	// Basic math
	Add
	Sub
	Mul
	Div
	Pow
	Sin
	Cos
	Log // natural

	// Convert to given type
	ConvStr
	ConvInt
	ConvFloat
	ConvBool

	// Concat a string
	CatStr
	// Concat a list
	CatList
	// Index into string
	StrAt
	// get length of str/list
	LenStr
	LenList

	// Call Xin function
	// TODO: not sure if necessary, since function calls
	// are implied by syntax.
	Call
	// Call native function, e.g. 'out'
	CallNative
	// Return function call
	Ret
)

type Position struct {
	FilePath string
	Line     int
	Col      int
}

type Tok interface {
	String() string
	Position() Position
}

type Expr interface {
	String() string
	Pretty() string
	Position() string
	Compile() []Instruction
}

// TODO: not sure what this needs to look like
type Instruction interface {
	Arity()
}
