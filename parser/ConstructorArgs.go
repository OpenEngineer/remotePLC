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
func ConstructorArgs(args []string, positionalTemplate []interface{}, optionalTemplate map[string]interface{}) (positional []interface{}, optional map[string]interface{}) {

	if len(args) < len(positionalTemplate) {
		logger.WriteFatal("ConstructorArgs()", errors.New("num args must be greater or equal than "+string(len(positionalTemplate))))
	}

	// loop the positional pointers, and store the args
	//  the positional args must come before the optional args
	for i, p := range positionalTemplate {
		var err error

		switch p.(type) {
		case int64:
			var v int64
			v, err = strconv.ParseInt(args[i], 10, 64)
			positional = append(positional, v)
		case int:
			var v int64
			v, err = strconv.ParseInt(args[i], 10, 64)
			positional = append(positional, int(v))
		case float64:
			var v float64
			v, err = strconv.ParseFloat(args[i], 64)
			positional = append(positional, v)
		case string:
			positional = append(positional, args[i])
		default:
			err = errors.New("type not recognized while parsing positional arguments")
		}
		logger.WriteFatal("ConstructorArgs()", err)
	}

	// loop the remaining arguments
	i := len(positionalTemplate) // already treated
	for i < len(args) {
		var err error
		option := args[i]
		if defaultValue, ok := optionalTemplate[option]; ok {
			i = i + 1

			switch defaultValue.(type) {
			case int64:
				i = i + 1
				if i < len(args) {
					var v int64
					v, err = strconv.ParseInt(args[i], 10, 64)
					optional[option] = v
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case int:
				i = i + 1
				if i < len(args) {
					var v int64
					v, err = strconv.ParseInt(args[i], 10, 64)
					optional[option] = int(v)
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case float64:
				i = i + 1
				if i < len(args) {
					var v float64
					v, err = strconv.ParseFloat(args[i], 64)
					optional[option] = v
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case string:
				i = i + 1
				if i < len(args) {
					optional[option] = args[i]
				} else {
					err = errors.New("premature end of constructor args, arg to option " + option + " not found")
				}
			case bool:
				optional[option] = true
			default:
				err = errors.New("type not recognized while parsing optional arguments")
			}
		} else {
			err = errors.New("optional arg " + args[i] + " not recognized")
		}

		logger.WriteFatal("ConstructorArgs()", err)
	}

	// add the optional arg defaults if the option wasn't specified during construction
	for key, v := range optionalTemplate {
		if _, ok := optional[key]; !ok {
			optional[key] = v
		}
	}

	return positional, optional
}
