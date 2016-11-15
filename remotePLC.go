package main

import "./runTime"
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
	fmt.Println(runTime.GetRunTime())
}
