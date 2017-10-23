package mapping

import ()

/*

When mapping to a struct, we map each field from a path into another object.
When mapping to a collection, we map each item from another collection.
When mapping to a pointer, we map from a set of pointers.
  If all are nil, then map to nil.


*/

// Choice focus on Part of a Whole type.
type Choice interface {
	getChoice() Choice
}

// SumChoice choose a branch of a Sum type.
type SumChoice struct {
	Index int
}

// StructChoice choose a branch (Field) of a Struct type.
type StructChoice struct {
	Index string
}

// MappedField is a field of a new struct type created from an old struct type.
type MappedField struct {
	// SourcePath the sequence of choices from the old root struct to the
	//   old Field this MappedField corresponds to.
	SourcePath []Choice
	// Name (of the Field) should correspond to the last item in SourcePath.
	Name string

	// Struct is non-nil iff the field contains a struct.
	Struct *MappedStruct

	// Atom is non-nil iff the field doesn't contain a struct.
	Atom MappedAtom
}

// MappedAtom is opaque to the shrinking algorithm.
// The atom may be shrinkable on the inside, but fields cannot be pulled from it.
// nil is the Identity mapping (no-op).
type MappedAtom interface {
	getMappedAtom() MappedAtom
}

// MappedStruct a new struct type created from an old struct type.
// MappedStructs can be restructured using the shrinking algorithm.
type MappedStruct struct {
	// SourcePath the sequence of choices from the old root struct to the
	//   old Struct this MappedStruct corresponds to
	SourcePath []Choice
	Fields     []*MappedField
}

// MappedSlice create a new Slice type by mapping elements of an old Slice type.
type MappedSlice struct {
	Elem *MappedAtom
}

// MappedMap create a new Map type by mapping values of an old Map type.
type MappedMap struct {
	Value *MappedAtom
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

// Shrink move struct fields as close to the root struct as possible without
//   creating field-name collisions.
// TODO: When promoting from structs with only one field,
//   use the parent field name instead of the field name for the promoted field.
func (s *MappedStruct) Shrink() {
	for {
		promotionCount := 0

		// Count the occurrences of each field and subfield name.
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

		// Promote the subfields with unique names.
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
	for _, field := range s.Fields {
		if field.Atom != nil {
			continue
		}

		field.Struct.Shrink()
	}
}

// PromoteSubfield move a Field from its struct to the parent of its struct.
//   If its original struct is now empty, delete this struct from its parent.
func (s *MappedStruct) PromoteSubfield(fieldIx int, subfieldIx int) {
	field := s.Fields[fieldIx]
	subfield := field.Struct.Fields[subfieldIx]
	field.Struct.DeleteFieldAt(subfieldIx)
	s.InsertFieldAt(fieldIx, subfield)

	if len(field.Struct.Fields) == 0 {
		s.DeleteFieldAt(fieldIx)
	}
}

// InsertFieldAt insert a new Field into the Struct at a given index.
func (s *MappedStruct) InsertFieldAt(fieldIx int, field *MappedField) {
	s.Fields = append(s.Fields[:fieldIx],
		append([]*MappedField{field}, s.Fields[fieldIx:]...)...)
}

// DeleteFieldAt delete the Field at a given index.
func (s *MappedStruct) DeleteFieldAt(fieldIx int) {
	s.Fields = append(s.Fields[:fieldIx], s.Fields[fieldIx+1:]...)
}

func (c *SumChoice) getChoice() Choice {
	return c
}

func (c *StructChoice) getChoice() Choice {
	return c
}

func (s *MappedStruct) getMappedAtom() MappedAtom {
	return s
}

func (m *MappedSlice) getMappedAtom() MappedAtom {
	return m
}

func (m *MappedMap) getMappedAtom() MappedAtom {
	return m
}
