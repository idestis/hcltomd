package cmd

import (
	"testing"
)

func TestQuotedType(t *testing.T) {
	data := []byte("variable \"name\" { type = \"string\"  default = \"hcltomd\"}")
	_, err := hclToInterface(data)
	if err != nil {
		t.Errorf("TestHclToJSON failed during parsing. %v", err)
	}
}
