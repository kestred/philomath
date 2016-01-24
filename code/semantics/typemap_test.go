package semantics

import (
	"testing"

	"github.com/kestred/philomath/code/ast"
	"github.com/stretchr/testify/assert"
)

func TestTypeMap(t *testing.T) {
	// Get/Set
	typmap := MakeTypeMap()
	typmap.Set("a", ast.InferredUnsigned)
	assert.Equal(t, ast.InferredUnsigned, typmap.Get("a"))
	assert.Equal(t, ast.UnknownType, typmap.Get("x"))

	{
		// Get on reference
		refmap := typmap.Reference()
		assert.Equal(t, ast.InferredUnsigned, refmap.Get("a"))
		assert.Equal(t, typmap.types, refmap.types)

		// Set on reference
		refmap.Set("b", ast.InferredSigned)
		assert.Equal(t, ast.InferredSigned, refmap.Get("b"))
		assert.Equal(t, ast.InferredUnsigned, refmap.Get("a"))
		assert.NotEqual(t, typmap.types, refmap.types)

		// Overwriting a value
		refmap.Set("a", ast.InferredFloat)
		assert.Equal(t, ast.InferredFloat, refmap.Get("a"))
	}

	// Get on the original
	assert.Equal(t, ast.InferredUnsigned, typmap.Get("a"))
	assert.Equal(t, ast.UnknownType, typmap.Get("b"))
}
