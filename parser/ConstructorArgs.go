package parser

import (
	"../logger/"
	"errors"
	"os"
	"strconv"
)

func ParsePositionalArgs(args []string, positional []interface{}) {
	// loop the positional pointers, and store the args
	//  the positional args must come before the optional args
	for i, p := range positional {
		var err error

		switch p := p.(type) {
		case *int64:
			var v int64
			v, err = strconv.ParseInt(args[i], 10, 64)
			*p = v
		case *int:
			var v int64
			v, err = strconv.ParseInt(args[i], 10, 64)
			*p = int(v)
		case *float64:
			var v float64
			v, err = strconv.ParseFloat(args[i], 64)
			*p = v
		case *string:
			*p = args[i]
		case **os.File:
			fname := args[i]
			if fname == "stdin" {
				*p = os.Stdin
			} else {
				*p, err = os.Open(fname)
			}
		default:
			err = errors.New("type not recognized while parsing positional arguments")
		}
		logger.WriteFatal("ConstructorArgs()", err)
	}
}

func ParseOptionalArgs(args []string, optional map[string]interface{}) {
	// loop the args
	i := 0
	for i < len(args) {
		var err error
		option := args[i]
		if p, ok := optional[option]; ok {
			i = i + 1

			// check and assert type
			switch p := p.(type) {
			case *int64:
				i = i + 1
				if i < len(args) {
					var v int64
					v, err = strconv.ParseInt(args[i], 10, 64)
					*p = v
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case *int:
				i = i + 1
				if i < len(args) {
					var v int64
					v, err = strconv.ParseInt(args[i], 10, 64)
					*p = int(v)
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case *float64:
				i = i + 1
				if i < len(args) {
					var v float64
					v, err = strconv.ParseFloat(args[i], 64)
					*p = v
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case *string:
				i = i + 1
				if i < len(args) {
					*p = args[i]
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case *bool:
				*p = true
			default:
				err = errors.New("type not recognized while parsing optional arguments")
			}
		} else {
			err = errors.New("optional arg " + args[i] + " not recognized")
		}

		logger.WriteFatal("ConstructorArgs()", err)
	}
}

// input:
//  - args: list of strings taken directly from constructor file
//  - positional: list of pointers into which to store the parsed args (len(args) >= len(positional))
//  - optional: flags
func ParseArgs(args []string, positional []interface{}, optional map[string]interface{}) {

	if len(args) < len(positional) {
		logger.WriteFatal("ConstructorArgs()", errors.New("num args must be greater or equal than "+string(len(positional))))
	}

	ParsePositionalArgs(args[0:len(positional)], positional)

	ParseOptionalArgs(args[len(positional):], optional)
}

// take a variadic list of pointer arguments and construct a slice containing these pointers
func PositionalArgs(ptrs ...interface{}) (positional []interface{}) {
	for _, ptr := range ptrs {
		positional = append(positional, ptr)
	}
	return positional
}

// take a variadic paired list of strings and pointers and construct a map
// key string first, obj ptr second
func OptionalArgs(objs ...interface{}) (optional map[string]interface{}) {
	// collect the inputs
	var keys []string
	var ptrs []interface{}
	for i, obj := range objs {
		if i%2 == 0 { //string key
			switch obj := obj.(type) {
			case string:
				keys = append(keys, obj)
			default:
				logger.WriteFatal("OptionalArgs()", errors.New("key not of string type"))
			}
		} else { // ptr
			ptrs = append(ptrs, obj)
		}
	}

	// now make the map
	for i, key := range keys {
		optional[key] = ptrs[i]
	}

	return optional
}
