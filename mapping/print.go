package mapping

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
)

// These (and init()) are just for formatting the Print methods.
var maxDepth = 40
var indents = make([]string, maxDepth)

func init() {
	for index := range indents {
		indents[index] = strings.Repeat("  ", index)
	}
}

// PrintChoices print a sequence of choices on a line.
func PrintChoices(depth int, tag string, choices []Choice) {
	fmt.Printf("%s%s @[ ", indents[depth], tag)
	for _, choice := range choices {
		switch choice := choice.(type) {
		case *SumChoice:
			fmt.Printf("%d ", choice.Index)
		case *StructChoice:
			var ix string
			if len(choice.Index) == 0 {
				ix = "<anonymous>"
			} else {
				ix = choice.Index
			}

			fmt.Printf("%v ", ix)
		default:
			glog.Fatal("inconceivable")
		}
	}

	fmt.Println("]")
}

// Print the mapping.
func (m *MappedField) Print(depth int) {
	var name string
	if len(m.Name) > 0 {
		name = m.Name
	} else {
		name = "<anonymous>"
	}
	PrintChoices(depth, name, m.SourcePath)
	if m.Atom == nil {
		fmt.Printf("%s<Identity>\n", indents[depth])
	} else {
		m.Atom.Print(depth + 1)
	}
}

// Print the mapping.
func (m *MappedStruct) Print(depth int) {
	PrintChoices(depth, "struct", m.SourcePath)
	for _, field := range m.Fields {
		field.Print(depth + 1)
	}
}

// Print the mapping.
func (m *MappedSlice) Print(depth int) {
	fmt.Printf("%sslice\n", indents[depth])
	if m.Elem == nil {
		fmt.Printf("%s<Identity>\n", indents[depth])
	} else {
		m.Elem.Print(depth)
	}
}

// Print the mapping.
func (m *MappedMap) Print(depth int) {
	fmt.Printf("%smap\n", indents[depth])
	if m.Value == nil {
		fmt.Printf("%s<Identity>\n", indents[depth])
	} else {
		m.Value.Print(depth)
	}
}
