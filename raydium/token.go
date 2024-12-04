package raydium

import "github.com/gagliardetto/solana-go"

var (
	TOKEN_PROGRAM_ID = solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
)

type Token struct {
	ProgramId solana.PublicKey
	Mint      solana.PublicKey
	Decimals  int
	Symbol    string
	Name      string
}
