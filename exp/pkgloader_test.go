package exp

import (
	"os"
	"testing"

	"golang.org/x/tools/go/types"
)

func TestLoad(t *testing.T) {
	prog, err := Load("github.com/wdamron/astx")
	if err != nil {
		t.Fatal(err)
	}
	typ := prog.Underlying("astx", "StructField")
	if typ == nil {
		t.Fatal("underlying type not found")
	}
	if !IsStruct(typ) {
		t.Fatal("underlying type should be *types.Struct")
	}
	s := typ.(*types.Struct)
	if s.NumFields() != 7 {
		t.Fatal("7 fields should be found in underlying type")
	}
	field := s.Field(0)
	if field.Name() != "Name" {
		t.Error("Name should be found in underlying type's fields")
	}
	if !IsBasic(field.Type()) {
		t.Error("Basic type should be found in underlying type's field (0)")
	}
	basic := field.Type().(*types.Basic)
	if basic.Kind() != types.String {
		t.Error("String type should be found in underlying type's field (0)")
	}
}

func TestLoadDir(t *testing.T) {
	prog, err := LoadDir(os.Getenv("GOPATH") + "/src/github.com/wdamron/astx")
	if err != nil {
		t.Fatal(err)
	}
	typ := prog.Underlying("astx", "StructField")
	if typ == nil {
		t.Fatal("underlying type not found")
	}
	if !IsStruct(typ) {
		t.Fatal("underlying type should be *types.Struct")
	}
	s := typ.(*types.Struct)
	if s.NumFields() != 7 {
		t.Fatal("7 fields should be found in underlying type")
	}
	field := s.Field(0)
	if field.Name() != "Name" {
		t.Error("Name should be found in underlying type's fields")
	}
	if !IsBasic(field.Type()) {
		t.Error("Basic type should be found in underlying type's field (0)")
	}
	basic := field.Type().(*types.Basic)
	if basic.Kind() != types.String {
		t.Error("String type should be found in underlying type's field (0)")
	}
}

func TestLoadFile(t *testing.T) {
	prog, err := LoadFile(os.Getenv("GOPATH") + "/src/github.com/wdamron/astx/astx.go")
	if err != nil {
		t.Fatal(err)
	}
	typ := prog.Underlying("astx", "StructField")
	if typ == nil {
		t.Fatal("underlying type not found")
	}
	if !IsStruct(typ) {
		t.Fatal("underlying type should be *types.Struct")
	}
	s := typ.(*types.Struct)
	if s.NumFields() != 7 {
		t.Fatal("7 fields should be found in underlying type")
	}
	field := s.Field(0)
	if field.Name() != "Name" {
		t.Error("Name should be found in underlying type's fields")
	}
	if !IsBasic(field.Type()) {
		t.Error("Basic type should be found in underlying type's field (0)")
	}
	basic := field.Type().(*types.Basic)
	if basic.Kind() != types.String {
		t.Error("String type should be found in underlying type's field (0)")
	}
}
