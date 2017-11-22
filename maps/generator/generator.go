// Helper binary for generating generic map serialization.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	packageName = flag.String("package", "", "package name")
	keyType     = flag.String("key", "", "key type")
	valueType   = flag.String("value", "", "value type")
)

const (
	mapFile = `package {{.Package}}

import (
	"encoding/json"
	"sort"
)

// {{.KVName}} represents a key-value pair.
type {{.KVName}} struct {
	Key   {{.Key}}
	Value {{.Value}}
}

// {{.MapName}} represents a JSON map from {{.Key}} to {{.Value}}. Keys must be unique. The map is
// represented as a list of key-value pairs sorted ascending by key.
type {{.MapName}} []{{.KVName}}

// Len returns the number of entries.
func (m {{.MapName}}) Len() int {
	return len(m)
}

// Less returns true if element i has a smaller key than element j.
func (m {{.MapName}}) Less(i, j int) bool {
	return m[i].Key < m[j].Key
}

// Swap swaps the elements i and j.
func (m {{.MapName}}) Swap(i, j int) {
	temp := m[i]
	m[i] = m[j]
	m[j] = temp
}

// UnmarshalJSON puts the JSON map into this map representation.
func (m *{{.MapName}}) UnmarshalJSON(b []byte) error {
	var val map[{{.Key}}]{{.Value}}
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}
	var ret []{{.KVName}}
	for k, v := range val {
		s := {{.KVName}}{Key: k, Value: v}
		ret = append(ret, s)
	}
	sort.Sort({{.MapName}}(ret))
	*m = ret
	return nil
}

// MarshalJSON returns a JSON map corresponding to its representation.
func (m {{.MapName}}) MarshalJSON() ([]byte, error) {
	val := make(map[{{.Key}}]{{.Value}})
	for _, v := range m {
		val[v.Key] = v.Value
	}
	return json.Marshal(val)
}
`

	testFile = `package {{.Package}}

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test{{.MapName}}(t *testing.T) {
	js := ` + "`" + `
{
	"foo": "a",
	"bar": "b",
	"baz": "c"
}` + "`" + `
	var sb {{.MapName}}
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
	expect := {{.MapName}}([]{{.KVName}}{
		{"bar", "b"},
		{"baz", "c"},
		{"foo", "a"},
	})

	if !reflect.DeepEqual(sb, expect) {
		t.Errorf("got %v; want %v", sb, expect)
	}
}`
)

type params struct {
	Package string
	Key     string
	Value   string
	MapName string
	KVName  string
}

func main() {
	flag.Parse()
	prefix := strings.Title(*keyType) + strings.Title(*valueType)
	kvName := prefix + "KV"
	mapName := prefix + "Map"

	p := params{
		Package: *packageName,
		Key:     *keyType,
		Value:   *valueType,
		MapName: mapName,
		KVName:  kvName,
	}

	mapTemplate, err := template.New("map").Parse(mapFile)
	if err != nil {
		log.Fatal(err)
	}
	testTemplate, err := template.New("test").Parse(testFile)
	if err != nil {
		log.Fatal(err)
	}

	mapFilename := strings.ToLower(fmt.Sprintf("%s_%s_map.go", *keyType, *valueType))

	if _, err = os.Stat(mapFilename); os.IsNotExist(err) {
		mapFile, err := os.Create(mapFilename)
		if err != nil {
			log.Fatal(err)
		}
		defer mapFile.Close()
		err = mapTemplate.Execute(mapFile, p)
		if err != nil {
			log.Fatal(err)
		}
	}

	testFilename := strings.ToLower(fmt.Sprintf("%s_%s_map_test.go", *keyType, *valueType))
	if _, err = os.Stat(testFilename); os.IsNotExist(err) {
		testFile, err := os.Create(testFilename)
		if err != nil {
			log.Fatal(err)
		}
		defer testFile.Close()
		err = testTemplate.Execute(testFile, p)
		if err != nil {
			log.Fatal(err)
		}
	}
}
