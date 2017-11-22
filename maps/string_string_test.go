package maps

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStringString(t *testing.T) {
	js := `
{
	"foo": "a",
	"bar": "b",
	"baz": "c"
}`
	var sb StringString
	err := json.Unmarshal([]byte(js), &sb)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(sb)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, &sb)
	if err != nil {
		t.Fatal(err)
	}

	if len(sb) != 3 {
		t.Fatalf("sb = %v; want 3 elements", sb)
	}
	expect := StringString([]StringStringKV{
		{"bar", "b"},
		{"baz", "c"},
		{"foo", "a"},
	})

	if !reflect.DeepEqual(sb, expect) {
		t.Errorf("got %v; want %v", sb, expect)
	}
}
