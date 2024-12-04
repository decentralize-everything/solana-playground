package raydium

import (
	"errors"

	"github.com/gagliardetto/solana-go"
)

type MarketStateV3 struct {
	// Unknown1               [5]byte
	// Unknown2               [8]byte
	OwnAddress             [32]byte
	VaultSignerNonce       uint64
	BaseMint               [32]byte
	QuoteMint              [32]byte
	BaseVault              [32]byte
	BaseDepositsTotal      uint64
	BaseFeesAccrued        uint64
	QuoteVault             [32]byte
	QuoteDepositsTotal     uint64
	QuoteFeesAccrued       uint64
	QuoteDustThreshold     uint64
	RequestQueue           [32]byte
	EventQueue             [32]byte
	Bids                   [32]byte
	Asks                   [32]byte
	BaseLotSize            uint64
	QuoteLotSize           uint64
	FeeRateBps             uint64
	ReferrerRebatesAccrued uint64
	// Unknown3               [7]byte
}

func GetAssociatedAuthority(programId, marketId []byte) (solana.PublicKey, error) {
	seed := [][]byte{marketId}
	nonce := byte(0)

	for nonce < 100 {
		seedWithNonce := append(seed, []byte{nonce}, []byte{0, 0, 0, 0, 0, 0, 0})
		publicKey, err := solana.CreateProgramAddress(seedWithNonce, solana.PublicKeyFromBytes(programId))
		if err != nil {
			nonce++
			continue
		} else {
			return publicKey, nil
		}
	}
	return solana.PublicKey{}, errors.New("unable to find a viable program address nonce")
}
