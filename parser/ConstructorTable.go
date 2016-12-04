package parser

import (
	"../logger/"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type ConstructorTable [][]string

func (t *ConstructorTable) ReadAppendFile(fname string, lineSeparators []string) (table [][]string) {
	file, err := os.Open(fname)
	defer file.Close()
	logger.WriteError("parseFileTable()", err)

	// get the size of the file
	finfo, errStat := os.Stat(fname)
	logger.WriteError("parseFileTable()", errStat)

	fileBytes := make([]byte, int(finfo.Size()))

	// read the whole file into memory
	_, errRead := file.Read(fileBytes)
	logger.WriteError("parseFileTable()", errRead)

	// now process the string with the lower level readStringTable() function
	t.ReadAppendString(string(fileBytes), lineSeparators)

	return
}

func (t *ConstructorTable) ReadAppendString(str string, lineSeparators []string) {
	lines := []string{str}

	for _, ls := range lineSeparators {
		lines = splitLines(lines, ls)
	}

	// now split each line-string into its words
	for _, l := range lines {
		*t = append(*t, strings.Fields(l))
	}
}

func splitLines(linesIn []string, fs string) (linesOut []string) {
	var split [][]string
	for _, l := range linesIn {
		split = append(split, strings.Split(l, fs))
	}

	for _, s := range split {
		linesOut = append(linesOut, s...)
	}

	return
}

func (t *ConstructorTable) RemoveComments(commentChar string) {
	var tmp [][]string

	for _, row := range *t {
		tmpRow := []string{}
		for _, word := range row {
			if string(word[0]) == commentChar {
				break
			} else if string(word[len(word)-1]) == commentChar {
				tmpRow = append(tmpRow, string(word[0:len(word)-1]))
				break
			} else {
				tmpRow = append(tmpRow, word)
			}
		}

		tmp = append(tmp, tmpRow)
	}

	*t = tmp
}

// remove rows if: len(row) < minNumWords
func (t *ConstructorTable) RemoveEmptyRows(minNumWords int) {
	var tmp [][]string

	for _, row := range *t {
		if len(row) >= minNumWords {
			tmp = append(tmp, row)
		}
	}

	*t = tmp
}

func (t *ConstructorTable) FilterTable(colI int, regexpStr string) (tableOut [][]string) {
	re := regexp.MustCompile(regexpStr)

	var tmpTable [][]string
	for _, row := range *t {
		if re.MatchString(row[colI]) {
			tableOut = append(tableOut, row)
		} else {
			tmpTable = append(tmpTable, row)
		}
	}

	*t = tmpTable

	return
}

func (t *ConstructorTable) GenerateMissingNames(colI int, regexpStr string, suffix string) {
	re := regexp.MustCompile(regexpStr)

	var tmpTable [][]string

	for _, row := range *t {
		var tmpRow []string
		if re.MatchString(row[colI]) {
			tmpRow = []string{slugify(row) + suffix}
			tmpRow = append(tmpRow, row...)
		} else {
			tmpRow = row
		}

		tmpTable = append(tmpTable, tmpRow)
	}

	*t = tmpTable
}

func slugify(words []string) string {
	slug := words[0]

	for i := 1; i < len(words); i++ {
		slug += "_" + words[i]
	}

	return slug
}

func (t *ConstructorTable) AddRow(row []string) {
	*t = append(*t, row)
}

func (t *ConstructorTable) MergeRows(mergeChar string) {
	var tmpTable [][]string

	var tmpRow []string
	for i, row := range *t {
		if len(row) > 0 && row[len(row)-1] == mergeChar {
			if i == len(*t)-1 {
				logger.WriteFatal("MergeRows()", errors.New("last row contains the mergeChar \""+mergeChar+"\", this is illegal"))
			}
			tmpRow = append(tmpRow, row[0:len(row)-1]...)
		} else {
			tmpRow = append(tmpRow, row...)
			tmpTable = append(tmpTable, tmpRow)
			tmpRow = []string{}
		}
	}

	*t = tmpTable
}

func (t *ConstructorTable) Print() {
	for _, row := range *t {
		fmt.Println(row)
	}
}

// the suffix cannot contain characters that might conflict with the regexp
func (t *ConstructorTable) CorrectSuffixes(suffix string, colI int) {
	suffixRegexp := ".*" + suffix + "$"
	re := regexp.MustCompile(suffixRegexp)

	// record all the names:
	var altNames []string
	names := make(map[string]string)

	for _, row := range *t {
		name := row[0]

		altName := name + suffix
		if re.MatchString(name) {
			altName = string(name[0 : len(name)-len(suffix)])
		}

		altNames = append(altNames, altName)
		names[altName] = name
	}

	sort.Strings(altNames)

	// for i >= colI, replace all words that match altNames, by equivalent in names
	var tmpTable [][]string
	for _, row := range *t {
		var tmpRow []string
		tmpRow = row

		for j := colI; j < len(row); j++ {
			k := sort.SearchStrings(altNames, row[j])

			if k < len(altNames) && altNames[k] == row[j] {
				altName := altNames[k]
				tmpRow[j] = names[altName]
			}
		}

		tmpTable = append(tmpTable, tmpRow)
	}

	*t = tmpTable
}

func (t *ConstructorTable) WordToLine(matchWord string) {
	var tmpTable [][]string

	for _, row := range *t {
		matchI := -1
		for i, word := range row {
			if word == matchWord {
				matchI = i
			}
		}

		if matchI >= 0 {
			tmpTable = append(tmpTable, row[0:matchI])
			tmpTable = append(tmpTable, []string{row[matchI]})
			tmpTable = append(tmpTable, row[matchI+1:])
		} else {
			tmpTable = append(tmpTable, row)
		}
	}

	*t = tmpTable
}

func (t *ConstructorTable) SubstituteSingleWordLine(matchWord string, wordI [][]int, line0 []string) {
	var tmpTable [][]string
	tmpTable = *t

	// this can be done in-place
	for i, row := range *t {
		if len(row) == 1 && row[0] == matchWord {
			var tmpRow []string

			for _, wI := range wordI {
				if wI[0] == 0 {
					tmpRow = append(tmpRow, line0[wI[1]])
				} else {
					tmpRow = append(tmpRow, tmpTable[i+wI[0]][wI[1]])
				}
			}

			tmpTable[i] = tmpRow
		}
	}

	*t = tmpTable
}

func (t *ConstructorTable) DetectDuplicates(colI int) {
	nameI := make(map[string]int)

	var tmpTable [][]string
	tmpTable = *t

	for i, row := range *t {
		name := row[colI]
		if _, ok := nameI[name]; ok {
			logger.WriteError("DetectDuplicates()", errors.New("duplicate name detected: "+name+" (vs. \""+tmpTable[nameI[name]][colI]+"\" found earlier)"))
		} else {
			nameI[name] = i
		}
	}
}

func (t *ConstructorTable) ContainsWord(matchWord string, colI int) bool {
	b := false

	for _, row := range *t {
		for i := colI; i < len(row); i++ {
			if row[i] == matchWord {
				b = true
				break
			}
		}

		if b {
			break
		}
	}

	return b
}
