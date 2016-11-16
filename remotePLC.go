package main

import "./tables/"
import "fmt"

var b = a
var d = getStr()

func main() {
	//fmt.Println("Hello")
	fmt.Println(A)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(tables.GetRunTime())

	inputTable := [][]string{
		[]string{"var1", "ConstantInput", "666.666"},
		[]string{"var2", "ConstantInput", "669.666"},
		[]string{"var3", "ZeroInput"},
	}

	tables.ConstructInputs(inputTable)
	for _, v := range tables.Inputs {
		fmt.Println(v.Get())
	}
}
