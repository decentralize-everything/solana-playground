package raydium

import (
	"fmt"
	"math/big"
	"unsafe"

	"github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"
)

type PoolInfo struct {
	Status        *big.Int `json:"status"`
	BaseDecimals  int      `json:"coin_decimals"`
	QuoteDecimals int      `json:"pc_decimals"`
	LpDecimals    int      `json:"lp_decimals"`
	BaseReserve   *big.Int `json:"pool_coin_amount"`
	QuoteReserve  *big.Int `json:"pool_pc_amount"`
	LpSupply      *big.Int `json:"pool_lp_supply"`
	StartTime     int      `json:"pool_open_time"`
}

func ComputeAmountOut(ammInfo *AmmInfo, poolInfo *PoolInfo, inputToken *Token, outputToken *Token, inputAmount *big.Int, slippage int64) (*big.Int, *big.Int) {
	if !includesToken(ammInfo, inputToken) || !includesToken(ammInfo, outputToken) {
		fmt.Println("token not match with pool")
	}

	reserveIn, reserveOut := poolInfo.BaseReserve, poolInfo.QuoteReserve
	if inputToken.Mint.Equals(ammInfo.QuoteMint) {
		reserveIn, reserveOut = reserveOut, reserveIn
	}

	feeRaw := decimal.NewFromBigInt(inputAmount, 0).Mul(decimal.NewFromInt(25)).Div(decimal.NewFromInt(10000)).IntPart()
	amountInWithFee := new(big.Int).Sub(inputAmount, big.NewInt(feeRaw))
	denominator := new(big.Int).Add(reserveIn, amountInWithFee)
	amountOutRaw := new(big.Int).Div(new(big.Int).Mul(reserveOut, amountInWithFee), denominator)

	minAmountOutRaw := new(big.Int).Div(new(big.Int).Mul(amountOutRaw, big.NewInt(100)), big.NewInt(100+slippage))
	return amountOutRaw, minAmountOutRaw
}

func includesToken(ammInfo *AmmInfo, token *Token) bool {
	return ammInfo.BaseMint.Equals(token.Mint) || ammInfo.QuoteMint.Equals(token.Mint)
}

type RewardInfo struct {
	RewardState           uint8
	OpenTime              uint64
	EndTime               uint64
	LastUpdateTime        uint64
	EmissionsPerSecondX64 *big.Int
	RewardTotalEmissioned uint64
	RewardClaimed         uint64
	TokenMint             solana.PublicKey
	TokenVault            solana.PublicKey
	Creator               solana.PublicKey
	RewardGrowthGlobalX64 *big.Int
}

func NewRewardInfoFromBytes(data []byte) *RewardInfo {
	return &RewardInfo{
		RewardState:           data[0],
		OpenTime:              *(*uint64)(unsafe.Pointer(&data[1])),
		EndTime:               *(*uint64)(unsafe.Pointer(&data[9])),
		LastUpdateTime:        *(*uint64)(unsafe.Pointer(&data[17])),
		EmissionsPerSecondX64: new(big.Int).SetBytes(reverseByteSlice(data[25:41])),
		RewardTotalEmissioned: *(*uint64)(unsafe.Pointer(&data[41])),
		RewardClaimed:         *(*uint64)(unsafe.Pointer(&data[49])),
		TokenMint:             solana.PublicKeyFromBytes(data[57:89]),
		TokenVault:            solana.PublicKeyFromBytes(data[89:121]),
		Creator:               solana.PublicKeyFromBytes(data[121:153]),
		RewardGrowthGlobalX64: new(big.Int).SetBytes(reverseByteSlice(data[153:169])),
	}
}

type PoolInfoLayout struct {
	Bump                      uint8
	AmmConfig                 solana.PublicKey
	Creator                   solana.PublicKey
	MintA                     solana.PublicKey
	MintB                     solana.PublicKey
	VaultA                    solana.PublicKey
	VaultB                    solana.PublicKey
	ObservationId             solana.PublicKey
	MintDecimalsA             uint8
	MintDicimalsB             uint8
	TickSpacing               uint16
	Liquidity                 *big.Int
	SqrtPriceX64              *big.Int
	TickCurrent               int32
	ObservationIndex          uint16
	ObservationUpdateDuration uint16
	FeeGrowthGlobalX64A       *big.Int
	FeeGrowthGlobalX64B       *big.Int
	ProtocolFeesTokenA        uint64
	ProtocolFeesTokenB        uint64
	SwapInAmountTokenA        *big.Int
	SwapOutAmountTokenB       *big.Int
	SwapInAmountTokenB        *big.Int
	SwapOutAmountTokenA       *big.Int
	Status                    uint8
	RewardInfos               []*RewardInfo
	TickArrayBitmap           []uint64
	TotalFeesTokenA           uint64
	TotalFeesClaimedTokenA    uint64
	TotalFeesTokenB           uint64
	TotalFeesClaimedTokenB    uint64
	FundFeesTokenA            uint64
	FundFeesTokenB            uint64
	StartTime                 uint64
}

func NewPoolInfoLayoutFromBytes(data []byte) *PoolInfoLayout {
	var rewardInfos []*RewardInfo
	for i := 0; i < 3; i++ {
		rewardInfos = append(rewardInfos, NewRewardInfoFromBytes(data[397+(i*169):397+((i+1)*169)]))
	}

	var tickArrayBitmap []uint64
	for i := 0; i < 16; i++ {
		tickArrayBitmap = append(tickArrayBitmap, *(*uint64)(unsafe.Pointer(&data[904+(i*8)])))
	}

	return &PoolInfoLayout{
		Bump:                      data[8],
		AmmConfig:                 solana.PublicKeyFromBytes(data[9:41]),
		Creator:                   solana.PublicKeyFromBytes(data[41:73]),
		MintA:                     solana.PublicKeyFromBytes(data[73:105]),
		MintB:                     solana.PublicKeyFromBytes(data[105:137]),
		VaultA:                    solana.PublicKeyFromBytes(data[137:169]),
		VaultB:                    solana.PublicKeyFromBytes(data[169:201]),
		ObservationId:             solana.PublicKeyFromBytes(data[201:233]),
		MintDecimalsA:             data[233],
		MintDicimalsB:             data[234],
		TickSpacing:               *(*uint16)(unsafe.Pointer(&data[235])),
		Liquidity:                 new(big.Int).SetBytes(reverseByteSlice(data[237:253])),
		SqrtPriceX64:              new(big.Int).SetBytes(reverseByteSlice(data[253:269])),
		TickCurrent:               *(*int32)(unsafe.Pointer(&data[269])),
		ObservationIndex:          *(*uint16)(unsafe.Pointer(&data[273])),
		ObservationUpdateDuration: *(*uint16)(unsafe.Pointer(&data[275])),
		FeeGrowthGlobalX64A:       new(big.Int).SetBytes(reverseByteSlice(data[277:293])),
		FeeGrowthGlobalX64B:       new(big.Int).SetBytes(reverseByteSlice(data[293:309])),
		ProtocolFeesTokenA:        *(*uint64)(unsafe.Pointer(&data[309])),
		ProtocolFeesTokenB:        *(*uint64)(unsafe.Pointer(&data[317])),
		SwapInAmountTokenA:        new(big.Int).SetBytes(reverseByteSlice(data[325:341])),
		SwapOutAmountTokenB:       new(big.Int).SetBytes(reverseByteSlice(data[341:357])),
		SwapInAmountTokenB:        new(big.Int).SetBytes(reverseByteSlice(data[357:373])),
		SwapOutAmountTokenA:       new(big.Int).SetBytes(reverseByteSlice(data[373:389])),
		Status:                    data[389],
		RewardInfos:               rewardInfos,
		TickArrayBitmap:           tickArrayBitmap,
		TotalFeesTokenA:           *(*uint64)(unsafe.Pointer(&data[1032])),
		TotalFeesClaimedTokenA:    *(*uint64)(unsafe.Pointer(&data[1040])),
		TotalFeesTokenB:           *(*uint64)(unsafe.Pointer(&data[1048])),
		TotalFeesClaimedTokenB:    *(*uint64)(unsafe.Pointer(&data[1056])),
		FundFeesTokenA:            *(*uint64)(unsafe.Pointer(&data[1064])),
		FundFeesTokenB:            *(*uint64)(unsafe.Pointer(&data[1072])),
		StartTime:                 *(*uint64)(unsafe.Pointer(&data[1080])),
	}
}

func reverseByteSlice(data []byte) []byte {
	for i := 0; i < len(data)/2; i++ {
		data[i], data[len(data)-1-i] = data[len(data)-1-i], data[i]
	}
	return data
}
