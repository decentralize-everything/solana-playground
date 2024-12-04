package raydium

import "github.com/gagliardetto/solana-go"

type TokenAccount struct {
	PublicKey   solana.PublicKey
	ProgramId   solana.PublicKey
	AccountInfo *SplAccount
}
