package parser

import (
	"../logger/"
	"errors"
	"strconv"
)

// input:
//  - args: list of strings taken directly from constructor file
//  - positional: list of pointers into which to store the parsed args (len(args) >= len(positional))
//  - optional: flags
func ConstructorArgs(args []string, positional []*interface{}, optional map[string]*interface{}) {

	if len(args) < len(positional) {
		logger.WriteFatal("ConstructorArgs()", errors.New("num args must be greater or equal than "+string(len(positional))))
	}

	// loop the positional pointers, and store the args
	//  the positional args must come before the optional args
	for i, ptr := range positional {
		var err error

		p := *ptr
		switch p.(type) {
		case int64:
			var v int64
			v, err = strconv.ParseInt(args[i], 10, 64)
			*ptr = v
		case int:
			var v int64
			v, err = strconv.ParseInt(args[i], 10, 64)
			*ptr = int(v)
		case float64:
			var v float64
			v, err = strconv.ParseFloat(args[i], 64)
			*ptr = v
		case string:
			*ptr = args[i]
		default:
			err = errors.New("type not recognized while parsing positional arguments")
		}
		logger.WriteFatal("ConstructorArgs()", err)
	}

	// loop the remaining arguments
	i := len(positional) // already treated
	for i < len(args) {
		var err error
		option := args[i]
		if ptr, ok := optional[option]; ok {
			i = i + 1

			defaultValue := *ptr

			switch defaultValue.(type) {
			case int64:
				i = i + 1
				if i < len(args) {
					var v int64
					v, err = strconv.ParseInt(args[i], 10, 64)
					*ptr = v
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case int:
				i = i + 1
				if i < len(args) {
					var v int64
					v, err = strconv.ParseInt(args[i], 10, 64)
					*ptr = int(v)
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case float64:
				i = i + 1
				if i < len(args) {
					var v float64
					v, err = strconv.ParseFloat(args[i], 64)
					*ptr = v
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case string:
				i = i + 1
				if i < len(args) {
					*ptr = args[i]
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case bool:
				*ptr = true
			default:
				err = errors.New("type not recognized while parsing optional arguments")
			}
		} else {
			err = errors.New("optional arg " + args[i] + " not recognized")
		}

		logger.WriteFatal("ConstructorArgs()", err)
	}
}
