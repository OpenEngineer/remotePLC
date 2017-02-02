package parser

const (
	COMMENT_CHAR = "#"
)

// piggy back on the ConstructorTable functionality
func TokenizeFile(fname string) [][]string {
	var tokens ConstructorTable

	tokens.ReadAppendFile(fname, []string{"\n"})
	tokens.RemoveComments(COMMENT_CHAR)

	return tokens
}
