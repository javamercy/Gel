package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseObjectType(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedType  ObjectType
		expectedValid bool
	}{
		{"valid blob", "blob", ObjectTypeBlob, true},
		{"valid tree", "tree", ObjectTypeTree, true},
		{"valid commit", "commit", ObjectTypeCommit, true},
		{"invalid empty", "", "", false},
		{"invalid random", "random", "", false},
		{"invalid case", "Blob", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objType, valid := ParseObjectType(tt.input)
			assert.Equal(t, tt.expectedType, objType)
			assert.Equal(t, tt.expectedValid, valid)
		})
	}
}

func TestObjectType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		objType  ObjectType
		expected bool
	}{
		{"valid blob", ObjectTypeBlob, true},
		{"valid tree", ObjectTypeTree, true},
		{"valid commit", ObjectTypeCommit, true},
		{"invalid empty", ObjectType(""), false},
		{"invalid random", ObjectType("random"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.objType.IsValid())
		})
	}
}
