package maps

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStringBool(t *testing.T) {
	js := `
{
	"foo": true,
	"bar": false,
	"baz": true
}`
	var sb StringBool
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
	expect := StringBool([]StringBoolKV{
		{"bar", false},
		{"baz", true},
		{"foo", true},
	})

	if !reflect.DeepEqual(sb, expect) {
		t.Errorf("got %v; want %v", sb, expect)
	}
}
