package raydium

type SplMint struct {
	// MintAuthorityOption   uint32
	MintAuthority         [32]byte
	Supply                uint64
	Decimals              uint8
	IsInitialized         uint8
	FreezeAuthorityOption uint32
	// FreezeAuthority       [32]byte
}
