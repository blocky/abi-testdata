package main

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type SliceOfBytes [][]byte

func (v SliceOfBytes) ABIEncode() ([]byte, error) {
	args := abi.Arguments{
		{
			Type: abi.Type{
				T:    abi.SliceTy,               // Slice (dynamic array)
				Elem: &abi.Type{T: abi.BytesTy}, // Element type is bytes (dynamic)
			},
		},
	}

	encoded, err := args.Pack(v)
	if err != nil {
		return nil, fmt.Errorf("packing: %w", err)
	}

	// and make sure that we can decode it back
	x, err := args.Unpack(encoded)
	if err != nil {
		return nil, fmt.Errorf("unpacking: %w", err)
	}
	got := x[0].([][]byte)

	// and check if they are the same
	if len(v) != len(got) {
		return nil, fmt.Errorf("check failed, not equal length %d %d", len(v), len(got))
	}

	for i := range v {
		if !bytes.Equal(v[i], got[i]) {
			return nil, fmt.Errorf("element %d not equal %v %v", i, v[i], got[i])
		}
	}

	return encoded, nil
}
