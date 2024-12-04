package raydium

import (
	"unsafe"

	"github.com/gagliardetto/solana-go"
)

type TickArrayBitmap struct {
	PoolId                  solana.PublicKey `json:"poolId"`
	PositiveTickArrayBitmap [][]uint64       `json:"positiveTickArrayBitmap"`
	NegativeTickArrayBitmap [][]uint64       `json:"negativeTickArrayBitmap"`
}

func NewTickArrayBitmapFromBytes(data []byte) *TickArrayBitmap {
	positiveTickArrayBitmap := make([][]uint64, 0, 14)
	for i := 0; i < 14; i++ {
		row := make([]uint64, 0, 8)
		for j := 0; j < 8; j++ {
			row = append(row, *(*uint64)(unsafe.Pointer(&data[40+i*64+j*8])))
		}
		positiveTickArrayBitmap = append(positiveTickArrayBitmap, row)
	}
	negativeTickArrayBitmap := make([][]uint64, 0, 14)
	for i := 0; i < 14; i++ {
		row := make([]uint64, 0, 8)
		for j := 0; j < 8; j++ {
			row = append(row, *(*uint64)(unsafe.Pointer(&data[936+i*64+j*8])))
		}
		negativeTickArrayBitmap = append(negativeTickArrayBitmap, row)
	}

	return &TickArrayBitmap{
		PoolId:                  solana.PublicKeyFromBytes(data[8:40]),
		PositiveTickArrayBitmap: positiveTickArrayBitmap,
		NegativeTickArrayBitmap: negativeTickArrayBitmap,
	}
}
