package raydium

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type AmmInfo struct {
	Id                 solana.PublicKey
	BaseMint           solana.PublicKey
	QuoteMint          solana.PublicKey
	LpMint             solana.PublicKey
	BaseDecimals       uint64
	QuoteDecimals      uint64
	LpDecimals         uint8
	Version            uint64
	ProgramId          solana.PublicKey
	Authority          solana.PublicKey
	OpenOrders         solana.PublicKey
	TargetOrders       solana.PublicKey
	BaseVault          solana.PublicKey
	QuoteVault         solana.PublicKey
	WithdrawQueue      solana.PublicKey
	LpVault            solana.PublicKey
	MarketVersion      uint64
	MarketProgramId    solana.PublicKey
	MarketId           solana.PublicKey
	MarketAuthority    solana.PublicKey
	MarketBaseVault    solana.PublicKey
	MarketQuoteVault   solana.PublicKey
	MarketBids         solana.PublicKey
	MarketAsks         solana.PublicKey
	MarketEventQueue   solana.PublicKey
	LookupTableAccount solana.PublicKey
}

func (amm AmmInfo) Display() string {
	return fmt.Sprintf("Id: %v\n BaseMint: %v\n QuoteMint: %v\n LpMint: %v\n BaseDecimals: %v\n QuoteDecimals: %v\n LpDecimals: %v\n Version: %v\n ProgramId: %v\n Authority: %v\n OpenOrders: %v\n TargetOrders: %v\n BaseVault: %v\n QuoteVault: %v\n WithdrawQueue: %v\n LpVault: %v\n MarketVersion: %v\n MarketProgramId: %v\n MarketId: %v\n MarketAuthority: %v\n MarketBaseVault: %v\n MarketQuoteVault: %v\n MarketBids: %v\n MarketAsks: %v\n MarketEventQueue: %v\n LookupTableAccount: %v\n", amm.Id.String(), amm.BaseMint.String(), amm.QuoteMint.String(), amm.LpMint.String(), amm.BaseDecimals, amm.QuoteDecimals, amm.LpDecimals, amm.Version, amm.ProgramId.String(), amm.Authority.String(), amm.OpenOrders.String(), amm.TargetOrders.String(), amm.BaseVault.String(), amm.QuoteVault.String(), amm.WithdrawQueue.String(), amm.LpVault.String(), amm.MarketVersion, amm.MarketProgramId.String(), amm.MarketId.String(), amm.MarketAuthority.String(), amm.MarketBaseVault.String(), amm.MarketQuoteVault.String(), amm.MarketBids.String(), amm.MarketAsks.String(), amm.MarketEventQueue.String(), amm.LookupTableAccount.String())
}

func GetAmmInfo(client *rpc.Client, id string, userAccount solana.PublicKey) (*AmmInfo, error) {
	pubKey := solana.MustPublicKeyFromBase58(id)
	account, err := client.GetAccountInfo(context.TODO(), pubKey)
	if err != nil {
		return nil, err
	}
	owner := account.Value.Owner
	fmt.Println("size of liquidity state:", unsafe.Sizeof(LiquidityStateV4{}), ", size of data:", len(account.Value.Data.GetBinary()))
	liquidityState := (*(*LiquidityStateV4)(unsafe.Pointer(&account.Value.Data.GetBinary()[0])))

	fmt.Println("market id", solana.PublicKeyFromBytes(liquidityState.MarketId[:]).String())
	marketAccount, err := client.GetAccountInfo(context.TODO(), solana.PublicKeyFromBytes(liquidityState.MarketId[:]))
	if err != nil {
		return nil, err
	}
	fmt.Println("size of market state", unsafe.Sizeof(MarketStateV3{}), ", size of data:", len(marketAccount.Value.Data.GetBinary()[13:]))
	marketState := (*(*MarketStateV3)(unsafe.Pointer(&marketAccount.Value.Data.GetBinary()[13])))

	fmt.Println("market event queue", solana.PublicKeyFromBytes(marketState.EventQueue[:]).String())
	lpMintAccount, err := client.GetAccountInfo(context.TODO(), solana.PublicKeyFromBytes(liquidityState.LpMint[:]))
	if err != nil {
		return nil, err
	}
	fmt.Println("size of spl mint", unsafe.Sizeof(SplMint{}), ", size of data:", len(lpMintAccount.Value.Data.GetBinary()[4:]))
	splMint := (*(*SplMint)(unsafe.Pointer(&lpMintAccount.Value.Data.GetBinary()[4])))

	// "amm authority"
	authority, _, err := solana.FindProgramAddress([][]byte{{97, 109, 109, 32, 97, 117, 116, 104, 111, 114, 105, 116, 121}}, owner)
	if err != nil {
		return nil, err
	}

	marketAuthority, err := GetAssociatedAuthority(liquidityState.MarketProgramId[:], liquidityState.MarketId[:])
	if err != nil {
		return nil, err
	}

	return &AmmInfo{
		Id:               pubKey,
		BaseMint:         solana.PublicKeyFromBytes(liquidityState.BaseMint[:]),
		QuoteMint:        solana.PublicKeyFromBytes(liquidityState.QuoteMint[:]),
		LpMint:           solana.PublicKeyFromBytes(liquidityState.LpMint[:]),
		BaseDecimals:     liquidityState.BaseDecimal,
		QuoteDecimals:    liquidityState.QuoteDecimal,
		LpDecimals:       splMint.Decimals,
		Version:          4,
		ProgramId:        owner,
		Authority:        authority,
		OpenOrders:       solana.PublicKeyFromBytes(liquidityState.OpenOrders[:]),
		TargetOrders:     solana.PublicKeyFromBytes(liquidityState.TargetOrders[:]),
		BaseVault:        solana.PublicKeyFromBytes(liquidityState.BaseVault[:]),
		QuoteVault:       solana.PublicKeyFromBytes(liquidityState.QuoteVault[:]),
		WithdrawQueue:    solana.PublicKeyFromBytes(liquidityState.WithdrawQueue[:]),
		LpVault:          solana.PublicKeyFromBytes(liquidityState.LpVault[:]),
		MarketVersion:    3,
		MarketProgramId:  solana.PublicKeyFromBytes(liquidityState.MarketProgramId[:]),
		MarketId:         solana.PublicKeyFromBytes(liquidityState.MarketId[:]),
		MarketAuthority:  marketAuthority,
		MarketBaseVault:  solana.PublicKeyFromBytes(marketState.BaseVault[:]),
		MarketQuoteVault: solana.PublicKeyFromBytes(marketState.QuoteVault[:]),
		MarketBids:       solana.PublicKeyFromBytes(marketState.Bids[:]),
		MarketAsks:       solana.PublicKeyFromBytes(marketState.Asks[:]),
		MarketEventQueue: solana.PublicKeyFromBytes(marketState.EventQueue[:]),
		// MarketEventQueue:   solana.MustPublicKeyFromBase58("2CoBP2rr5HmjMdPC4nMwnYg1cdH9JPUuqbq2QGSMGfms"),
		LookupTableAccount: userAccount,
	}, nil
}
