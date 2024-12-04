package raydium

import "unsafe"

var (
	SPL_ACCOUNT_SIZE = 165
)

type SplAccount struct {
	Mint                 [32]byte
	Owner                [32]byte
	Amount               uint64
	DelegateOption       uint32
	Delegate             [32]byte
	State                uint8
	IsNativeOption       uint32
	IsNative             uint64
	DelegatedAmount      uint64
	CloseAuthorityOption uint32
	CloseAuthority       [32]byte
}

func NewSplAccountFromBytes(b []byte) *SplAccount {
	return &SplAccount{
		Mint:                 *(*[32]byte)(b[0:32]),
		Owner:                *(*[32]byte)(b[32:64]),
		Amount:               *(*uint64)(unsafe.Pointer(&b[64])),
		DelegateOption:       *(*uint32)(unsafe.Pointer(&b[72])),
		Delegate:             *(*[32]byte)(b[76:108]),
		State:                *(*uint8)(unsafe.Pointer(&b[108])),
		IsNativeOption:       *(*uint32)(unsafe.Pointer(&b[109])),
		IsNative:             (*(*uint64)(unsafe.Pointer(&b[113]))),
		DelegatedAmount:      *(*uint64)(unsafe.Pointer(&b[121])),
		CloseAuthorityOption: *(*uint32)(unsafe.Pointer(&b[129])),
		CloseAuthority:       *(*[32]byte)(b[133:165]),
	}
}
