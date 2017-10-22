package samples

import (
	"go/format"
	"testing"

	"github.com/koki/shorthand/gen"
)

func TestSerialization(t *testing.T) {
	formattedFileContents, err := format.Source([]byte(PodSrc))
	if err != nil {
		t.Errorf("Error formatting PodSrc: %v", err)
	}

	_, f := gen.DeserializeFileAST(string(formattedFileContents))

	buf := gen.SerializeFileAST(f)

	if string(formattedFileContents) != buf.String() {
		t.Errorf("Error in Serialization/Deserialization logic Expected [%s] Found [%s]", string(formattedFileContents), buf.String())
	}
}
