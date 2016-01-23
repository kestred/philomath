package semantics

import "github.com/kestred/philomath/ast"

// TypeMap is a map of variable names to types.
// The Reference() method can be used to get a copy-on-write reference to the map.
type TypeMap struct {
	types       map[string]ast.Type
	copyOnWrite bool
}

func MakeTypeMap() TypeMap {
	return TypeMap{make(map[string]ast.Type), false}
}

func (typmap *TypeMap) Get(name string) ast.Type {
	if typ, ok := typmap.types[name]; ok {
		return typ
	} else {
		return nil
	}
}

func (typmap *TypeMap) Set(name string, typ ast.Type) {
	if typmap.copyOnWrite {
		typesCopy := make(map[string]ast.Type)
		for k, v := range typmap.types {
			typesCopy[k] = v
		}

		typmap.types = typesCopy
		typmap.copyOnWrite = false
	}

	typmap.types[name] = typ
}

func (typmap *TypeMap) Reference() TypeMap {
	return TypeMap{typmap.types, true}
}
