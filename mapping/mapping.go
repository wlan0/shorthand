package mapping

import (
	"github.com/golang/glog"
)

/*

When mapping to a struct, we map each field from a path into another object.
When mapping to a collection, we map each item from another collection.
When mapping to a pointer, we map from a set of pointers.
  If all are nil, then map to nil.


*/

type MappingDefinition interface {
	getMappingDefinition() MappingDefinition
}

type MappedFromField struct {
	Field string
}

type MappedStruct struct {
	Fields map[string]*Mapping
}

type MappedMap struct {
	Value *Mapping
}

type MappedSlice struct {
	Value *Mapping
}

type MappedSum struct {
	Choices []*Mapping
}

type Mapping struct {
	From       TypeWrapper
	To         TypeWrapper
	Definition MappingDefinition
}

func IdentityMappingforSum(t *Sum) *MappedSum {
	return &MappedSum{make([]*Mapping, len(t.Choices))}
}

func IdentityMappingForTypeIdent(t *TypeIdent) MappingDefinition {
	switch d := t.Definition.(type) {
	case *StructDefinition:
		return IdentityMappingForStructType(t, d)
	case *WrapperDefinition:
		// Never modify simple wrapper types.
		return nil
	}

	glog.Fatal("Inconceivable!")
	return nil
}

func IdentityMappingForStructType(t *TypeIdent, d *StructDefinition) *MappedStruct {
	fields := make(map[string]*Mapping, len(d.Fields))
	for _, field := range d.Fields {
		fields[field.Name] = &Mapping{
			From:       t,
			To:         field.Type,
			Definition: &MappedFromField{field.Name},
		}
	}

	return &MappedStruct{fields}
}

func IdentityMapping(t TypeWrapper) *Mapping {
	var d MappingDefinition
	switch t := t.(type) {
	case *TypeIdent:
		d := IdentityMappingForTypeIdent(t)
	case *Map:
		d := &MappedMap{}
	case *Slice:
		d := &MappedSlice{}
	case *Sum:
		d := IdentityMappingforSum(t)
	}

	return &Mapping{From: t, To: t, Definition: d}
}

/*
Shrinking recurses like this:
1. Try changing Mapping at this level. This creates a set of child mappings.
2. Try changing the child mappings.
3. Build the return type from the children's return types.
4. If no changes were made, continue. Otherwise, return to 1.

TODO: Deal with branches properly
*/

// Context is used to make sure we create unique type names.
// NOTE: Anonymous structs here?
type Context struct {
}

// ShrinkWrapper is the top-level "generate a mapping" function.
func (c *Context) ShrinkWrapper(t TypeWrapper) *Mapping {
	return c.Shrink(IdentityMapping(t))
}

func (c *Context) Shrink(m *Mapping) *Mapping {
	switch d := m.Definition.(type) {
	case *MappedStruct:
	case *MappedMap:
	case *MappedSlice:
	case *MappedSum:
	}
}

// The best we can do is shrink the values.
func (c *Context) ShrinkMap(d *MappedMap) *Mapping {
	valueMapping := c.ShrinkWrapper(t.Value)
	if valueMapping == nil {
		// Identity mapping (no-op).
		return nil
	}

	return &Mapping{
		From:       t,
		To:         &Map{Key: t.Key, Value: valueMapping.To},
		Definition: &MappedMap{Value: valueMapping},
	}
}

// The best we can do is shrink the values.
func (c *Context) ShrinkSlice(t *Slice) *Mapping {
	valueMapping := c.ShrinkWrapper(t.Value)
	if valueMapping == nil {
		// Identity mapping (no-op).
		return nil
	}

	return &Mapping{
		From:       t,
		To:         &Slice{valueMapping.To},
		Definition: &MappedSlice{valueMapping},
	}
}

func (c *Context) ShrinkSum(t *Sum) *Mapping {
	return nil
}

func (c *Context) ShrinkTypeIdent(t *TypeIdent) *Mapping {
	return nil
}

func (d *MappedFromField) getMappingDefinition() MappingDefinition {
	return d
}

func (d *MappedStruct) getMappingDefinition() MappingDefinition {
	return d
}

func (d *MappedMap) getMappingDefinition() MappingDefinition {
	return d
}

func (d *MappedSlice) getMappingDefinition() MappingDefinition {
	return d
}

func (d *MappedSum) getMappingDefinition() MappingDefinition {
	return d
}
