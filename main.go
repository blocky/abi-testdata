package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
)

func exitOnError(msg string, err error) {
	if err != nil {
		localErr := fmt.Errorf("%s: %w", msg, err)
		log.Fatal(localErr)
	}
}

func mustHexDecode(s string) []byte {
	r, err := hex.DecodeString(s)
	exitOnError("hex decoding", err)
	return r
}

type SliceOfByteSliceExample struct {
	Name    string
	Native  [][]byte
	Encoded []byte
}

func makeSliceOfByteSliceExamples() ([]SliceOfByteSliceExample, error) {
	const (
		emptyBytes   = "emptyBytes"
		oneByte      = "oneByte"
		someBytes1   = "someBytes1"
		someBytes2   = "someBytes2"
		alignedBytes = "alignedBytes"
	)

	var values = map[string][]byte{
		emptyBytes:   {},
		oneByte:      []byte("a"),
		someBytes1:   []byte("hello"),
		someBytes2:   []byte("world"),
		alignedBytes: []byte("32-bytes-xxxxxxxxxxxxxxxxxxxxxxx"),
	}

	examples := []SliceOfByteSliceExample{}
	for _, v := range [][]string{
		{emptyBytes},
		{oneByte},
		{someBytes1},
		{someBytes2},
		{alignedBytes},
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
		name := "slice-of-bytes"
		native := [][]byte{}
		for _, vi := range v {
			name = name + "-" + vi
			native = append(native, values[vi])
		}

		encoded, err := SliceOfBytes(native).ABIEncode()
		if err != nil {
			return nil, fmt.Errorf("encoding slice of bytes: %w", err)
		}

		examples = append(examples, SliceOfByteSliceExample{
			Name:    name,
			Native:  native,
			Encoded: encoded,
		})
	}

	return examples, nil
}

type AllIntsExample struct {
	Name    string
	Native  AllInts
	Encoded []byte
}

func makeAllIntsExamples() ([]AllIntsExample, error) {
	examples := []AllIntsExample{}
	for _, v := range []struct {
		name string
		val  AllInts
	}{
		{
			name: "all-ints-xxx",
			val:  AllInts{Val1: 179131, Val2: 137, Val3: 27919810352},
		}, {
			name: "all-ints-0x0",
			val:  AllInts{Val1: 0, Val2: 1, Val3: 0},
		}, {
			name: "all-ints-x0x",
			val:  AllInts{Val1: 10, Val2: 0, Val3: 250},
		},
	} {
		data, err := v.val.ABIEncode()
		if err != nil {
			return nil, fmt.Errorf("encoding: %w", err)
		}

		examples = append(examples, AllIntsExample{
			Name:    v.name,
			Native:  v.val,
			Encoded: data,
		})
	}

	return examples, nil
}

type IntAndBytesExample struct {
	Name    string
	Native  IntAndBytes
	Encoded []byte
}

func makeIntAndBytesExamples() ([]IntAndBytesExample, error) {
	examples := []IntAndBytesExample{}
	for _, v := range []struct {
		name string
		val  IntAndBytes
	}{
		{
			name: "int-an-bytes-xxx",
			val:  IntAndBytes{Int1: 1, Bytes1: []byte("hello"), Bytes2: []byte("world")},
		}, {
			name: "int-an-bytes-0x0",
			val:  IntAndBytes{Int1: 0, Bytes1: []byte("hello"), Bytes2: []byte{}},
		}, {
			name: "all-ints-x0x",
			val:  IntAndBytes{Int1: 1, Bytes1: []byte{}, Bytes2: []byte("world")},
		},
	} {
		data, err := v.val.ABIEncode()
		if err != nil {
			return nil, fmt.Errorf("encoding: %w", err)
		}

		examples = append(examples, IntAndBytesExample{
			Name:    v.name,
			Native:  v.val,
			Encoded: data,
		})
	}

	return examples, nil
}

func main() {
	packageName := flag.String("package-name", "abitestdata", "package name")
	flag.Parse()

	sliceOfBytes, err := makeSliceOfByteSliceExamples()
	exitOnError("creating slice of bytes examples", err)

	allInts, err := makeAllIntsExamples()
	exitOnError("creating all ints examples", err)

	intAndBytes, err := makeIntAndBytesExamples()
	exitOnError("creating int and bytes examples", err)

	renderer := Renderer{
		packageName:  *packageName,
		sliceOfBytes: sliceOfBytes,
		allInts:      allInts,
		intAndBytes:  intAndBytes,
	}
	err = renderer.Render(os.Stdout)
	exitOnError("rendering code", err)
}
