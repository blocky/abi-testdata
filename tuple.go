package main

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type AllInts struct {
	Val1 uint64
	Val2 uint64
	Val3 uint64
}

func (t AllInts) ABIEncode() ([]byte, error) {
	abiType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "val1", Type: "uint64"},
		{Name: "val2", Type: "uint64"},
		{Name: "val3", Type: "uint64"},
	})
	if err != nil {
		return nil, fmt.Errorf("creating type: %w", err)
	}

	args := abi.Arguments{{Name: "allInts", Type: abiType}}
	encoded, err := args.Pack(struct {
		Val1 uint64 `json:"val1"`
		Val2 uint64 `json:"val2"`
		Val3 uint64 `json:"val3"`
	}{
		Val1: t.Val1,
		Val2: t.Val2,
		Val3: t.Val3,
	})
	if err != nil {
		return nil, fmt.Errorf("packing: %w", err)
	}

	// make sure that we can decode it back
	unpacked, err := args.Unpack(encoded)
	if err != nil {
		return nil, fmt.Errorf("unpacking: %w", err)
	}

	tuple, ok := unpacked[0].(struct {
		Val1 uint64 `json:"val1"`
		Val2 uint64 `json:"val2"`
		Val3 uint64 `json:"val3"`
	})
	if !ok {
		return nil, errors.New("casting unpacked data")
	}

	if t.Val1 != tuple.Val1 || t.Val2 != tuple.Val2 || t.Val3 != tuple.Val3 {
		return nil, errors.New("enc/dec produced different data")
	}

	return encoded, nil
}

type IntAndBytes struct {
	Int1   uint64
	Bytes1 []byte
	Bytes2 []byte
}

func (t *IntAndBytes) ABIEncode() ([]byte, error) {
	abiType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "int1", Type: "uint64"},
		{Name: "bytes1", Type: "bytes"},
		{Name: "bytes2", Type: "bytes"},
	})
	if err != nil {
		return nil, fmt.Errorf("creating type: %w", err)
	}

	args := abi.Arguments{
		{
			Name: "heterogeneous",
			Type: abiType,
		},
	}

	encoded, err := args.Pack(struct {
		Int1   uint64 `json:"int1"`
		Bytes1 []byte `json:"bytes1"`
		Bytes2 []byte `json:"bytes2"`
	}{
		Int1:   t.Int1,
		Bytes1: t.Bytes1,
		Bytes2: t.Bytes2,
	})
	if err != nil {
		return nil, fmt.Errorf("could not encode TWAP tuple: %w", err)
	}

	unpacked, err := args.Unpack(encoded)
	if err != nil {
		return nil, fmt.Errorf("unpacking: %w", err)
	}

	got, ok := unpacked[0].(struct {
		Int1   uint64 `json:"int1"`
		Bytes1 []byte `json:"bytes1"`
		Bytes2 []byte `json:"bytes2"`
	})
	if !ok {
		return nil, errors.New("could not cast")
	}

	if t.Int1 != got.Int1 ||
		!bytes.Equal(t.Bytes1, got.Bytes1) ||
		!bytes.Equal(t.Bytes2, got.Bytes2) {
		return nil, errors.New("input doesn't match decoded result")
	}

	// Note that we trim off the first 32 bytes of the
	// encoding because when we use the standard tooling for
	// abi encoding a tuple, if the tuples contains a dynamic type,
	// such as bytes, the go-ethereum pack method
	// adds a dynamic header.  Unfortunately,
	// when you try to decode with solidity
	// using abi.decode(data, (int, bytes, bytes))
	// the "dynamic" portion is expected to not be there.
	return encoded[32:], nil
}
