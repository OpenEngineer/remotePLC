package blocks

// only 0.0 and 1.0 are treated, all other numbers are ignored
func ReductionNodeAndOperator(x []float64) float64 {
	y := 1.0

	for _, v := range x {
		if v == 0.0 {
			y = 0.0
			break
		}
	}

	return y
}

var ReductionNodeAndOperatorOk = AddReductionNodeOperator("And", ReductionNodeAndOperator)

func ReductionNodeOrOperator(x []float64) float64 {
	y := 0.0

	for _, v := range x {
		if v == 1.0 {
			y = 1.0
			break
		}
	}

	return y
}

var ReductionNodeOrOperatorOk = AddReductionNodeOperator("Or", ReductionNodeOrOperator)
