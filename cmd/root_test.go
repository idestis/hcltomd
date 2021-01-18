package cmd

import "testing"

func TestHclToJSON(t *testing.T) {
	data := []byte("variable \"name\" { type = \"string\"}")
	_, err := hclToJSON(data)
	if err != nil {
		t.Errorf("TestHclToJSON failed during parsing")
	}
}

// func TestUnqotedType(t *testing.T) {
// 	data := []byte("variable \"name\" { type = string}")
// 	_, err := hclToJSON(data)
// 	if err != nil {
// 		t.Errorf("TestUnqotedType failed during parsing unquoted type.")
// 	}
// }
