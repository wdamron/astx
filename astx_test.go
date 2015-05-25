package astx

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	parsed, err := ParseFile("./example_src.go.txt")
	if err != nil {
		t.Fatal(err)
	}
	if parsed == nil {
		t.Fatal("parsed file should not be nil")
	}
	if parsed.Package != "astx" {
		t.Error("should parse package name from example.go.txt")
	}
	if parsed.Path == "" {
		t.Error("should include (non-empty) provided file path (./example.go.txt)")
	}
	if parsed.AbsPath == "" {
		t.Error("should resolve (non-empty) absolute path of provided file path")
	}
	if len(parsed.Imports) != 2 {
		t.Error("should parse (1) import specified in example.go.txt")
	} else {
		imp := parsed.Imports[0]
		if imp.Name != "fmt" {
			t.Error("should parse 'fmt' import specified in example.go.txt")
		}
		if imp.Path != `"fmt"` {
			t.Error("should parse path for 'fmt' import specified in example.go.txt")
		}
		if len(imp.Doc) != 1 {
			t.Error("should parse (1) doc comment above 'fmt' import specified in example.go.txt")
		} else {
			if imp.Doc[0] != "// very useful" {
				t.Error("should parse full doc comment above 'fmt' import specified in example.go.txt")
			}
		}
		if len(imp.Comments) != 1 {
			t.Error("should parse (1) doc comment above 'fmt' import specified in example.go.txt")
		} else {
			if imp.Comments[0] != "// short for format" {
				t.Error("should parse full comment beside 'fmt' import specified in example.go.txt")
			}
		}
		imp = parsed.Imports[1]
		if imp.Path != `"io"` {
			t.Error("should parse path for 'io' import specified in example.go.txt")
		}
	}
	if len(parsed.Structs) != 1 {
		t.Error("should parse (1) struct type defined in example.go.txt")
		return
	}
	s := parsed.Structs[0]
	if s.Name != "Point" {
		t.Error("should parse name for struct type defined in example.go.txt")
	}
	if len(s.Comments) != 1 {
		t.Error("should parse (1) comment for struct type defined in example.go.txt")
	}
	if s.Comments[0] != "// Point is a type of thing" {
		t.Error("should receive full contents of comment for struct type defined in example.go.txt")
	}
	if len(s.Fields) != 4 {
		t.Error("should parse (4) fields for struct type defined in example.go.txt")
	}
	if s.Fields[0].Name != "X" {
		t.Logf("bad field: %#v", s.Fields[0])
		t.Error("should parse names of fields in struct type defined in example.go.txt")
	}
	if s.Fields[0].Type != "int" {
		t.Logf("bad field: %#v", s.Fields[0])
		t.Error("should parse types of fields in struct type defined in example.go.txt")
	}
	if s.Fields[0].Tag.Get("tagz") != "hello" {
		t.Logf("bad field: %#v", s.Fields[0])
		t.Error("should parse tags of fields in struct type defined in example.go.txt")
	}
	if s.Fields[1].Name != "Y" {
		t.Logf("bad field: %#v", s.Fields[1])
		t.Error("should parse names of fields in struct type defined in example.go.txt")
	}
	if s.Fields[1].Type != "io.Reader" {
		t.Logf("bad field: %#v", s.Fields[1])
		t.Error("should parse types of fields in struct type defined in example.go.txt")
	}
	if s.Fields[1].Tag.Get("tagz") != "world" {
		t.Logf("bad field: %#v", s.Fields[1])
		t.Error("should parse tags of fields in struct type defined in example.go.txt")
	}
	if s.Fields[2].Name != "Z" {
		t.Logf("bad field: %#v", s.Fields[2])
		t.Error("should parse names of fields in struct type defined in example.go.txt")
	}
	if s.Fields[2].Type != "[2]*int" {
		t.Logf("bad field: %#v", s.Fields[2])
		t.Error("should parse types of fields in struct type defined in example.go.txt")
	}
	if s.Fields[2].Tag.Get("tagz") != "hello" {
		t.Logf("bad field: %#v", s.Fields[2])
		t.Error("should parse tags of fields in struct type defined in example.go.txt")
	}
	if s.Fields[3].Name != "ZZ" {
		t.Logf("bad field: %#v", s.Fields[3])
		t.Error("should parse names of fields in struct type defined in example.go.txt")
	}
	if s.Fields[3].Type != "map[string]*[SZ]int" {
		t.Logf("bad field: %#v", s.Fields[3])
		t.Error("should parse types of fields in struct type defined in example.go.txt")
	}
	if s.Fields[3].Tag.Get("tagz") != "world" {
		t.Logf("bad field: %#v", s.Fields[3])
		t.Error("should parse tags of fields in struct type defined in example.go.txt")
	}
}
