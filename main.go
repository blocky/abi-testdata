package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mustBeEqual(a, b [][]byte) {
	if len(a) != len(b) {
		log.Fatalf("not equal length %d %d", len(a), len(b))
	}

	for i := range a {
		if !bytes.Equal(a[i], b[i]) {
			log.Fatalf("not equal %v %v", a[i], b[i])
		}
	}
}

func encodeSliceOfByteSlice(input [][]byte) []byte {
	args := abi.Arguments{
		{
			Type: abi.Type{
				T:    abi.SliceTy,               // Slice (dynamic array)
				Elem: &abi.Type{T: abi.BytesTy}, // Element type is bytes (dynamic)
			},
		},
	}

	encoded, err := args.Pack(input)
	exitOnError(err)

	// and make sure that we can decode it back
	x, err := args.Unpack(encoded)
	exitOnError(err)
	got := x[0].([][]byte)
	mustBeEqual(input, got)

	return encoded
}

const (
	emptyBytes   = "emptyBytes"
	oneByte      = "oneByte"
	someBytes1   = "someBytes1"
	someBytes2   = "someBytes2"
	alignedBytes = "alignedBytes"
)

var examples = map[string][]byte{
	emptyBytes:   {},
	oneByte:      []byte("a"),
	someBytes1:   []byte("hello"),
	someBytes2:   []byte("world"),
	alignedBytes: []byte("32-bytes-xxxxxxxxxxxxxxxxxxxxxxx"),
}

func printExample(w io.Writer, name string, native string, encoded []byte) {
	fmt.Fprintf(w, `    { name: "%s", native: %s, encoded: hexDecode("%s") },`+"\n", name, native, hex.EncodeToString(encoded))
}

func printHeader(w io.Writer, packageName string) {
	tmplStr := `
package {{.PackageName}}

import "encoding/hex"

func hexDecode(s string) []byte {
	r, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return r
}

var (
{{- range $key, $value := .Examples }}
	{{$key}} = []byte("{{bytesToString $value}}")
{{- end }}
)

var testData = []struct {
	name    string
	native  [][]byte
	encoded []byte
}{
`
	funcMap := template.FuncMap{
		"bytesToString": func(b []byte) string {
			return string(b)
		},
	}

	tmpl, err := template.New("header").Funcs(funcMap).Parse(tmplStr)
	exitOnError(err)
	tmpl.Execute(w, map[string]any{
		"PackageName": packageName,
		"Examples":    examples,
	},
	)
}

func printFooter(w io.Writer) {
	fmt.Fprintln(w, "}")
}

func printExamplesSliceOfByteSlice(w io.Writer, packageName string) {
	printHeader(w, packageName)

	// empty
	input := [][]byte{}
	enc := encodeSliceOfByteSlice(input)
	printExample(w, "empty", "[][]byte{}", enc)

	// single elements
	for _, v := range []string{
		emptyBytes,
		oneByte,
		someBytes1,
		someBytes2,
		alignedBytes,
	} {
		input := [][]byte{examples[v]}
		enc := encodeSliceOfByteSlice(input)
		name := fmt.Sprintf("one-elt-%s", v)
		native := fmt.Sprintf("[][]byte{%s}", v)
		printExample(w, name, native, enc)
	}

	// multiple element examples
	for _, v := range [][]string{
		{someBytes1, emptyBytes},
		{emptyBytes, someBytes1},
		{emptyBytes, someBytes1, emptyBytes},
		{someBytes1, emptyBytes, someBytes2},
		{someBytes1, someBytes2},
		{someBytes1, alignedBytes},
		{alignedBytes, someBytes1},
		{alignedBytes, someBytes1, alignedBytes},
		{someBytes1, alignedBytes, someBytes2},
	} {
		input := [][]byte{}
		name := "multiple-elts"
		native := "[][]byte{"
		for _, v := range v {
			input = append(input, examples[v])
			name += "-" + v
			native += v + ", "
		}
		enc := encodeSliceOfByteSlice(input)
		native += "}"
		printExample(w, name, native, enc)
	}
	printFooter(w)
}

func main() {
	printExamplesSliceOfByteSlice(os.Stdout, "evmlink")
}
