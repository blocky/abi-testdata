package main

import (
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
)

type Renderer struct {
	packageName  string
	sliceOfBytes []SliceOfByteSliceExample
	allInts      []AllIntsExample
	intAndBytes  []IntAndBytesExample
}

func (r *Renderer) Render(w io.Writer) error {
	tmplStr := `// Code generated by abi-testdata DO NOT EDIT.
//
// See github/blocky/abi-testdata/README.md for more information

package {{.PackageName}}

import "encoding/hex"

func hexDecode(s string) []byte {
	r, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return r
}

type SliceOfBytesExample struct {
	name    string
	native  [][]byte
	encoded []byte
}

type AllInts struct {
	Val1 uint64
	Val2 uint64
	Val3 uint64
}

type AllIntsExample struct {
	name    string
	native  AllInts
	encoded []byte
}

type IntAndBytes struct {
	Int1   uint64
	Bytes1 []byte
	Bytes2 []byte
}

type IntAndBytesExample struct {
	name    string
	native  IntAndBytes
	encoded []byte
}

var testData = struct {
	sliceOfBytes []SliceOfBytesExample
	allInts      []AllIntsExample
	intAndBytes  []IntAndBytesExample
}{
	sliceOfBytes: []SliceOfBytesExample{
{{- range .SliceOfBytes }}
		{
			name: "{{.Name}}",
			native: [][]byte{
				{{- range .Native }}
				[]byte("{{toString .}}"),
				{{- end }}
			},
			encoded: hexDecode("{{hexEncode .Encoded}}"),
		},
{{- end }}
	},

	allInts: []AllIntsExample{
{{- range .AllInts }}
		{
			name: "{{.Name}}",
			native: AllInts{
				Val1: {{ .Native.Val1 }},
				Val2: {{ .Native.Val2 }},
				Val3: {{ .Native.Val3 }},
			},
			encoded: hexDecode("{{hexEncode .Encoded}}"),
		},
{{- end }}
	},

	intAndBytes: []IntAndBytesExample{
{{- range .IntAndBytes }}
		{
			name: "{{.Name}}",
			native: IntAndBytes{
				Int1:   {{ .Native.Int1 }},
				Bytes1: []byte("{{ toString .Native.Bytes1 }}"),
				Bytes2: []byte("{{ toString .Native.Bytes2 }}"),
			},
			encoded: hexDecode("{{hexEncode .Encoded}}"),
		},
{{- end }}
	},
}
`
	funcMap := template.FuncMap{
		"hexEncode": func(b []byte) string {
			return hex.EncodeToString(b)
		},
		"toString": func(b []byte) string {
			return string(b)
		},
	}

	tmpl, err := template.New("examples").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parsing the template: %w", err)
	}

	data := map[string]any{
		"PackageName":  r.packageName,
		"SliceOfBytes": r.sliceOfBytes,
		"AllInts":      r.allInts,
		"IntAndBytes":  r.intAndBytes,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	return nil
}
