package maps

import (
	"encoding/json"
	"sort"
)

// StringStringKV represents a key-value pair.
type StringStringKV struct {
	Key   string
	Value string
}

// StringString represents a JSON map. Keys must be unique. The map is
// represented as a list of key-value pairs sorted ascending by key.
type StringString []StringStringKV

// Len returns the number of entries.
func (m StringString) Len() int {
	return len(m)
}

// Less returns true if element i has a smaller key than element j.
func (m StringString) Less(i, j int) bool {
	return m[i].Key < m[j].Key
}

// Swap swaps the elements i and j.
func (m StringString) Swap(i, j int) {
	temp := m[i]
	m[i] = m[j]
	m[j] = temp
}

// UnmarshalJSON puts the JSON map into this map representation.
func (m *StringString) UnmarshalJSON(b []byte) error {
	var val map[string]string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}
	var ret []StringStringKV
	for k, v := range val {
		s := StringStringKV{Key: k, Value: v}
		ret = append(ret, s)
	}
	sort.Sort(StringString(ret))
	*m = ret
	return nil
}

// MarshalJSON returns a JSON map corresponding to its representation.
func (m StringString) MarshalJSON() ([]byte, error) {
	val := make(map[string]string)
	for _, v := range m {
		val[v.Key] = v.Value
	}
	return json.Marshal(val)
}
