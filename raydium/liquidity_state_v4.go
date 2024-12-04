package raydium

import "github.com/gagliardetto/solana-go"

type LiquidityStateV4 struct {
	Status                 uint64
	Nonce                  uint64
	MaxOrder               uint64
	Depth                  uint64
	BaseDecimal            uint64
	QuoteDecimal           uint64
	State                  uint64
	ResetFlag              uint64
	MinSize                uint64
	VolMaxCutRatio         uint64
	AmountWaveRatio        uint64
	BaseLotSize            uint64
	QuoteLotSize           uint64
	MinPriceMultiplier     uint64
	MaxPriceMultiplier     uint64
	SystemDecimalValue     uint64
	MinSeparateNumerator   uint64
	MinSeparateDenominator uint64
	TradeFeeNumerator      uint64
	TradeFeeDenominator    uint64
	PnlNumerator           uint64
	PnlDenominator         uint64
	SwapFeeNumerator       uint64
	SwapFeeDenominator     uint64
	BaseNeedTakePnl        uint64
	QuoteNeedTakePnl       uint64
	QuoteTotalPnl          uint64
	BaseTotalPnl           uint64
	PoolOpenTime           uint64
	PunishPcAmount         uint64
	PunishCoinAmount       uint64
	OrderbookToInitTime    uint64
	SwapBaseInAmount       [16]byte
	SwapQuoteOutAmount     [16]byte
	SwapBase2QuoteFee      uint64
	SwapQuoteInAmount      [16]byte
	SwapBaseOutAmount      [16]byte
	SwapQuote2BaseFee      uint64
	BaseVault              [32]byte
	QuoteVault             [32]byte
	BaseMint               [32]byte
	QuoteMint              [32]byte
	LpMint                 [32]byte
	OpenOrders             [32]byte
	MarketId               [32]byte
	MarketProgramId        [32]byte
	TargetOrders           [32]byte
	WithdrawQueue          [32]byte
	LpVault                [32]byte
	Owner                  [32]byte
	LpReserve              uint64
	Padding                [24]byte
}

func (state *LiquidityStateV4) GetMarketId() solana.PublicKey {
	return solana.PublicKeyFromBytes(state.MarketId[:])
}

func (state *LiquidityStateV4) GetLpMint() solana.PublicKey {
	return solana.PublicKeyFromBytes(state.LpMint[:])
}

type AmmAccount struct {
	Id             solana.PublicKey
	ProgramId      solana.PublicKey
	LiquidityState *LiquidityStateV4
}
