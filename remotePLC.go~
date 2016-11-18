package main

import "./blocks/"

func main() {
	inputTable := [][]string{
		[]string{"var1", "ConstantInput", "666.666"},
		[]string{"var2", "ConstantInput", "669.666"},
		[]string{"var3", "ZeroInput"},
		[]string{"var4", "ScaledInput", "2.0", "1.0", "ConstantInput", "2002"},
	}
	inputs := blocks.ConstructAll(inputTable)

	outputTable := [][]string{
		[]string{"out1", "FileOutput", "test"},
		[]string{"out2", "FileOutput", "stdout"},
	}
	outputs := blocks.ConstructAll(outputTable)

	lineTable := [][]string{
		[]string{"line1", "Line", "var4", "out1"},
		[]string{"line2", "JoinLine", "out2", "var1", "var2", "var3", "var4"},
	}
	lines := blocks.ConstructAll(lineTable)

	for _, v := range inputs {
		v.Update()
	}
	// TODO: where to place logic?
	for _, v := range lines {
		v.Update()
	}
	for _, v := range outputs {
		v.Update()
	}
}
