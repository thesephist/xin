package xin

import (
	"fmt"
	"strings"
)

type formEvaler func(*Frame, []Value) (Value, error)

type DefaultFormValue struct {
	name string
	eval formEvaler
}

func (v DefaultFormValue) String() string {
	return fmt.Sprintf("Default form %s", v.name)
}

func (v DefaultFormValue) Equal(o Value) bool {
	if ov, ok := o.(DefaultFormValue); ok {
		return v.name == ov.name
	}

	return false
}

func loadAllDefaultValues(fr *Frame) {
	// stdin := make(StreamValue)
	// go func() {
	// 	stdinReader := bufio.NewReader(os.Stdin)
	// 	for {
	// 		input, err := stdinReader.ReadString('\n')
	// 		if err == io.EOF {
	// 			break
	// 		}

	// 		stdin <- StringValue(input)
	// 	}
	// }()
	// fr.Scope["os::stdin"] = stdin

	stdout := make(StreamValue)
	go func() {
		for {
			out := <-stdout
			fmt.Printf(out.String())
		}
	}()
	fr.Scope["os::stdout"] = stdout
}

func loadAllDefaultForms(fr *Frame) {
	builtins := map[string]formEvaler{
		"+": addForm,
		"-": subtractForm,
		"*": multiplyForm,
		"/": divideForm,

		">": greaterForm,
		"<": lessForm,
		"=": equalForm,

		"&": andForm,
		"|": orForm,
		"^": xorForm,

		"vec":      vecForm,
		"vec-get":  vecGetForm,
		"vec-set!": vecSetForm,
		"vec-size": vecSizeForm,

		"map":      mapForm,
		"map-get":  mapGetForm,
		"map-set!": mapSetForm,
		"map-del!": mapDelForm,
		"map-size": mapSizeForm,

		"stream": streamForm,
		"->":     sourceForm,
		"<-":     sinkForm,
	}

	for name, evaler := range builtins {
		loadDefaultForm(fr, name, evaler)
	}
}

func loadDefaultForm(fr *Frame, name string, evaler formEvaler) {
	fr.Put(name, DefaultFormValue{
		name: name,
		eval: evaler,
	})
}

type IncorrectNumberOfArgsError struct {
	required int
	given    int
}

func (e IncorrectNumberOfArgsError) Error() string {
	return fmt.Sprintf("Incorrect number of arguments: requires %d but got %d",
		e.required, e.given)
}

type MismatchedArgumentsError struct {
	args []Value
}

func (e MismatchedArgumentsError) Error() string {
	ss := make([]string, len(e.args))
	for i, n := range e.args {
		ss[i] = n.String()
	}
	return fmt.Sprintf("Mismatched arguments: %s", strings.Join(ss, " "))
}

func addForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
	}

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst + cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			return cleanFirst + cleanSecond, nil
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			return cleanFirst + cleanSecond, nil
		}
	case VecValue:
		if cleanSecond, ok := second.(VecValue); ok {
			return VecValue(append(cleanFirst, cleanSecond...)), nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func subtractForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
	}

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst - cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			return cleanFirst - cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func multiplyForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
	}

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst * cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			return cleanFirst * cleanSecond, nil
		}
	case VecValue:
		if cleanSecond, ok := second.(IntValue); ok {
			result := make([]Value, 0)
			copy(result, cleanFirst)
			max := int(cleanSecond)
			for i := 0; i < max; i++ {
				result = append(result, cleanFirst...)
			}
			return VecValue(result), nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func divideForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
	}

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst / cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			return cleanFirst / cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func andForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst & cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func orForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst | cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func xorForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst ^ cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func greaterForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			if cleanFirst > cleanSecond {
				return IntValue(1), nil
			} else {
				return IntValue(0), nil
			}
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			cmp := strings.Compare(string(cleanFirst), string(cleanSecond))
			if cmp == 1 {
				return IntValue(1), nil
			} else {
				return IntValue(0), nil
			}
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}
func lessForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			if cleanFirst < cleanSecond {
				return IntValue(1), nil
			} else {
				return IntValue(0), nil
			}
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			cmp := strings.Compare(string(cleanFirst), string(cleanSecond))
			if cmp == -1 {
				return IntValue(1), nil
			} else {
				return IntValue(0), nil
			}
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func equalForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
	}

	if first.Equal(second) {
		return IntValue(1), nil
	} else {
		return IntValue(0), nil
	}
}
