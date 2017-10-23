package mapping

import (
	//"go/ast"

	"github.com/koki/shorthand/inspect"
)

/*

When mapping to a struct, we map each field from a path into another object.
When mapping to a collection, we map each item from another collection.
When mapping to a pointer, we map from a set of pointers.
  If all are nil, then map to nil.


*/

type AccessType struct {
	Whole TypeWrapper
	Part  TypeWrapper
}

type Access struct {
	AccessType
	Accessor Accessor
}

type Accessor interface {
	getAccessor() Accessor
}

type FieldAccess struct {
	Name string
}

type MapAccess struct {
}

type SliceAccess struct {
}

type PointerAccess struct {
}

type MappingAccess struct {
	Mapping Mapping
}

type AccessPath struct {
	Segments []Access
}

type MappingDefinition interface {
	getMappingDefinition() MappingDefinition
}

type MappedStruct struct {
	Fields map[FieldAccess]AccessPath
}

type MappedMap struct {
}

type MappedSlice struct {
}

type MappedSimplePointer struct {
	Field AccessPath
}

type MappedStructPointer struct {
	Fields map[FieldAccess]AccessPath
}

type Mapping struct {
	From       TypeWrapper
	To         TypeWrapper
	Definition MappingDefinition
}

func IdentityMappingFor(context *inspect.Context) *Mapping {
	/*
		rootType := RootType(context)
		mapping := &Mapping{From: rootType, To: rootType}
			switch root := context.TypeSpec.Type {
				case *
			}
	*/

	return nil
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

func (d *MappedSimplePointer) getMappingDefinition() MappingDefinition {
	return d
}

func (d *MappedStructPointer) getMappingDefinition() MappingDefinition {
	return d
}

func (a *FieldAccess) getAccessor() Accessor {
	return a
}

func (a *MapAccess) getAccessor() Accessor {
	return a
}

func (a *SliceAccess) getAccessor() Accessor {
	return a
}

func (a *PointerAccess) getAccessor() Accessor {
	return a
}
