package mapping

import (
	"github.com/golang/glog"
	"github.com/kr/pretty"
)

/*

When mapping to a struct, we map each field from a path into another object.
When mapping to a collection, we map each item from another collection.
When mapping to a pointer, we map from a set of pointers.
  If all are nil, then map to nil.


*/

type Choice interface {
	getChoice() Choice
}

type SumChoice struct {
	Index int
}
type StructChoice struct {
	Index string
}

type MappedField struct {
	SourcePath []Choice
	Name       string

	// Struct is non-nil if the field contains a struct.
	Struct *MappedStruct

	// Atomic if not a Struct.
	Atom *MappedAtom
}

type MappedStruct struct {
	SourcePath []Choice
	Fields     []*MappedField
}

// MappedAtom is opaque to the shrinking algorithm.
// The atom may be shrinkable on the inside, but fields cannot be pulled from it.
type MappedAtom struct {
}

func incrementFieldCount(fieldCounts map[string]int, fieldName string) {
	count := getFieldCount(fieldCounts, fieldName)
	fieldCounts[fieldName] = count + 1
}

func getFieldCount(fieldCounts map[string]int, fieldName string) int {
	if count, ok := fieldCounts[fieldName]; ok {
		return count
	}

	return 0
}

func (s *MappedStruct) Shrink() {
	for {
		promotionCount := 0
		fieldCounts := map[string]int{}
		for _, field := range s.Fields {
			incrementFieldCount(fieldCounts, field.Name)

			if field.Atom != nil {
				continue
			}

			for _, subfield := range field.Struct.Fields {
				incrementFieldCount(fieldCounts, subfield.Name)
			}
		}

		for fieldIx, field := range s.Fields {
			if field.Atom != nil {
				continue
			}

			for subfieldIx, subfield := range field.Struct.Fields {
				if getFieldCount(fieldCounts, subfield.Name) == 1 {
					promotionCount++
					s.PromoteSubfield(fieldIx, subfieldIx)
				}
			}
		}

		if promotionCount == 0 {
			break
		}
	}

	// Finished shrinking the top level. Shrink the next level down.
	for fieldIx, field := range s.Fields {
		if field.Atom != nil {
			continue
		}

		field.Struct.Shrink()
	}
}

func (s *MappedStruct) PromoteSubfield(fieldIx int, subfieldIx int) {
	field := s.Fields[fieldIx]
	subfield := field.Struct.Fields[subfieldIx]
	field.Struct.DeleteFieldAt(subfieldIx)
	s.InsertFieldAt(fieldIx, subfield)

	if len(field.Struct.Fields) == 0 {
		s.DeleteFieldAt(fieldIx)
	}
}

func (s *MappedStruct) InsertFieldAt(fieldIx int, field *MappedField) {
	s.Fields = append(s.Fields[:fieldIx],
		append([]*MappedField{field}, s.Fields[fieldIx:]...)...)
}

func (s *MappedStruct) DeleteFieldAt(fieldIx int) {
	s.Fields = append(s.Fields[:fieldIx], s.Fields[fieldIx+1:]...)
}

func (c *SumChoice) getChoice() Choice {
	return c
}

func (c *StructChoice) getChoice() Choice {
	return c
}
