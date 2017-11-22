// Package maps implements marshaling and unmarshaling of common map types,
// which aren't generally supported by datastores.
package maps

//go:generate go run $GOPATH/src/github.com/yuhanfang/riot/maps/generator/generator.go --package maps --key string --value string
//go:generate go run $GOPATH/src/github.com/yuhanfang/riot/maps/generator/generator.go --package maps --key string --value bool
