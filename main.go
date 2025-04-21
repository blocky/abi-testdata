package main

import (
	"encoding/hex"
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

//
//	func exampleDecode(data []byte) error {
//		tuple := TWAP{}
//		return DecodeTuple(data,
//			DecodeBytes(&tuple.BaseToken),
//			DecodeBytes(&tuple.QuoteToken),
//			DecodeUint64(&tuple.Price),
//			DecodeUint64(&tuple.ChainID),
//			DecodeUint64(&tuple.Timestamp),
//		)
//
//		// Final goal once all the pieces are in place
//		// err := NewTupleDecoder().
//		// 	Bytes(&tuple.BaseToken).
//		// 	Bytes(&tuple.QuoteToken).
//		// 	Uint64(&tuple.Price).
//		// 	Uint64(&tuple.ChainID).
//		// 	Uint64(&tuple.Timestamp).
//		// 	Decode(data)
//
// }
//
//	func exampleEncode() ([]byte, error) {
//		tuple := TWAP{}
//
//		return EncodeTuple(
//			EncodeBytes(&tuple.BaseToken),
//			EncodeBytes(&tuple.QuoteToken),
//			EncodeUint64(&tuple.Price),
//			EncodeUint64(&tuple.ChainID),
//			EncodeUint64(&tuple.Timestamp),
//		)
//
//		// final goal once all the pieces are in place
//		enc, err := NewTupleEncoder().
//			Bytes(&tuple.BaseToken).
//			Bytes(&tuple.QuoteToken).
//			Uint64(&tuple.Price).
//			Uint64(&tuple.ChainID).
//			Uint64(&tuple.Timestamp).
//			Encode()
//	}
//
// type ComponentDecoder struct{}
//
//	type EncoderResult struct {
//		Indirect bool
//		Data     []byte
//	}
//
// type Encoder func() (EncoderResult, error)
//
//	func EncodeTuple(encoders ...Encoder) ([]byte, error) {
//		return nil, errors.New("not implemented")
//	}
//
//	func EncodeBytes(v *[]byte) Encoder {
//		return nil
//	}
//
//	func EncodeUint64(*uint64) Encoder {
//		return nil
//	}
//
// type Decoder func(data []byte) ([]byte, error)
//
//	func ABIDecodeUint64(data []byte) (uint64, error) {
//		return 0, errors.New("not implemented")
//	}
//
//	func DecodeUint64(v *uint64) Decoder {
//		return func(data []byte) ([]byte, error) {
//			// check if it is long enough
//
//			head, tail := data[:32], data[:32]
//			vv, err := ABIDecodeUint64(head)
//			if err != nil {
//				return nil, fmt.Errorf("decoding: %w", err)
//			}
//			*v = vv
//
//			return tail, nil
//		}
//	}
//
//	func DecodeBytes(v *[]byte) Decoder {
//		return nil
//	}
//
//	func DecodeTuple(data []byte, decoders ...Decoder) error {
//		d := data
//		var err error
//		for i, decoder := range decoders {
//			d, err = decoder(d)
//			if err != nil {
//				return fmt.Errorf("deocding element %d", i)
//			}
//		}
//		return nil
//	}
//
//	func demoTWAP() ([]byte, error) {
//		input := TWAP{
//			BaseToken:  mustHexDecode("0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"),
//			QuoteToken: mustHexDecode("7ceb23fd6bc0f9ab16970adb4f2c47d08b8c7396"),
//			Price:      179131,
//			ChainID:    137,
//			Timestamp:  27919810352,
//		}
//
//		data, err := input.ABIEncode()
//		exitOnError("encoding TWAP", err)
//
//		got := &TWAP{}
//		err = got.ABIDecode(data)
//		exitOnError("decoding TWAP", err)
//
//		// then
//		if reflect.DeepEqual(input, got) {
//			exitOnError("", errors.New("input and got are not equal"))
//		}
//		return data, nil
//	}

func main() {
	sliceOfBytes, err := makeSliceOfByteSliceExamples()
	exitOnError("creating slice of bytes examples", err)

	allInts, err := makeAllIntsExamples()
	exitOnError("creating all ints examples", err)

	intAndBytes, err := makeIntAndBytesExamples()
	exitOnError("creating int and bytes examples", err)

	renderer := Renderer{
		packageName:  "abitestdata",
		sliceOfBytes: sliceOfBytes,
		allInts:      allInts,
		intAndBytes:  intAndBytes,
	}
	err = renderer.Render(os.Stdout)
	exitOnError("rendering code", err)
}
