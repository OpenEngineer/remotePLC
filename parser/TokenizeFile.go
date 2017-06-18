package parser

import (
	"../logger/"
	"strconv"
)

const (
	COMMENT_CHAR = "#"
)

// piggy back on the ConstructorTable functionality
func TokenizeFile(fname string) [][]string {
	var tokens ConstructorTable

	tokens.ReadAppendFile(fname, []string{"\n"})
	tokens.RemoveComments(COMMENT_CHAR)
	tokens.RemoveEmptyRows(1)

	return tokens
}

func VectorizeFile(fname string) []string {
	tokens := TokenizeFile(fname)

	vector := []string{}

	for _, row := range tokens {
		for _, token := range row {
			vector = append(vector, token)
		}
	}

	return vector
}

func VectorizeFileFloats(fname string) []float64 {
	vectorStr := VectorizeFile(fname)

	vector := make([]float64, len(vectorStr))

	for i, str := range vectorStr {
		v, e := strconv.ParseFloat(str, 64)

		logger.WriteError("VectorizeFileFloats()", e)

		vector[i] = v
	}

	return vector
}
