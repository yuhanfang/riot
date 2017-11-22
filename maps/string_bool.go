package maps

import (
	"encoding/json"
	"sort"
)

// StringBoolKV represents a key-value pair.
type StringBoolKV struct {
	Key   string
	Value bool
}

// StringBool represents a JSON map. Keys must be unique. The map is
// represented as a list of key-value pairs sorted ascending by key.
type StringBool []StringBoolKV

// Len returns the number of entries.
func (m StringBool) Len() int {
	return len(m)
}

// Less returns true if element i has a smaller key than element j.
func (m StringBool) Less(i, j int) bool {
	return m[i].Key < m[j].Key
}

// Swap swaps the elements i and j.
func (m StringBool) Swap(i, j int) {
	temp := m[i]
	m[i] = m[j]
	m[j] = temp
}

// UnmarshalJSON puts the JSON map into this map representation.
func (m *StringBool) UnmarshalJSON(b []byte) error {
	var val map[string]bool
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}
	var ret []StringBoolKV
	for k, v := range val {
		s := StringBoolKV{Key: k, Value: v}
		ret = append(ret, s)
	}
	sort.Sort(StringBool(ret))
	*m = ret
	return nil
}

// MarshalJSON returns a JSON map corresponding to its representation.
func (m StringBool) MarshalJSON() ([]byte, error) {
	val := make(map[string]bool)
	for _, v := range m {
		val[v.Key] = v.Value
	}
	return json.Marshal(val)
}
