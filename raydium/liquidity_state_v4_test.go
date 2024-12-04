package raydium

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"
	"time"
	"unsafe"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type ApiPoolInfoV4 struct {
	Id                 string `json:"id"`
	BaseMint           string `json:"baseMint"`
	QuoteMint          string `json:"quoteMint"`
	LpMint             string `json:"lpMint"`
	BaseDecimals       uint64 `json:"baseDecimals"`
	QuoteDecimals      uint64 `json:"quoteDecimals"`
	LpDecimals         uint64 `json:"lpDecimals"`
	Version            uint64 `json:"version"`
	ProgramId          string `json:"programId"`
	Authority          string `json:"authority"`
	OpenOrders         string `json:"openOrders"`
	TargetOrders       string `json:"targetOrders"`
	BaseVault          string `json:"baseVault"`
	QuoteVault         string `json:"quoteVault"`
	WithdrawQueue      string `json:"withdrawQueue"`
	LpVault            string `json:"lpVault"`
	MarketVersion      uint64 `json:"marketVersion"`
	MarketId           string `json:"marketId"`
	MarketProgramId    string `json:"marketProgramId"`
	MarketAuthority    string `json:"marketAuthority"`
	MarketBaseVault    string `json:"marketBaseVault"`
	MarketQuoteVault   string `json:"marketQuoteVault"`
	MarketBids         string `json:"marketBids"`
	MarketAsks         string `json:"marketAsks"`
	MarketEventQueue   string `json:"marketEventQueue"`
	LookupTableAccount string `json:"lookupTableAccount"`
}

type RouteBuildRequest struct {
	InputToken  string           `json:"inputToken"`
	OutputToken string           `json:"outputToken"`
	Amount      uint64           `json:"amount"`
	Slippage    uint64           `json:"slippage"`
	PublicKey   string           `json:"publicKey"`
	PoolsList   []*ApiPoolInfoV4 `json:"poolsList"`
	ClmmList    []*ClmmPoolInfo  `json:"clmmList"`
}

func BenchmarkFormatAmmKeys(t *testing.B) {
	client := rpc.New("https://aged-morning-glade.solana-mainnet.quiknode.pro/b57bbb1a4c8bdd409e1ac53aaedead26da057f59/")

	start := time.Now()
	allAmmAccount, err := client.GetProgramAccountsWithOpts(
		context.TODO(),
		solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"),
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					DataSize: 752,
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("allAmmAccount:", time.Since(start), "length:", len(allAmmAccount))
	start = time.Now()

	filterDefKey := solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
	amAccountData := make([]*AmmAccount, 0, len(allAmmAccount))
	allMarketProgram := make(map[string]struct{})
	for _, acc := range allAmmAccount {
		liquidityState := (*LiquidityStateV4)(unsafe.Pointer(&acc.Account.Data.GetBinary()[0]))
		if bytes.Equal(liquidityState.MarketProgramId[:], filterDefKey[:]) {
			continue
		}

		amAccountData = append(amAccountData, &AmmAccount{
			Id:             acc.Pubkey,
			ProgramId:      acc.Account.Owner,
			LiquidityState: liquidityState,
		})
		allMarketProgram[solana.PublicKeyFromBytes(liquidityState.MarketProgramId[:]).String()] = struct{}{}
	}
	t.Log("amAccountData:", time.Since(start), "length:", len(amAccountData))
	start = time.Now()

	marketInfo := make(map[string]*MarketInfo)
	for marketProgram := range allMarketProgram {
		allMarketInfo, err := client.GetProgramAccountsWithOpts(
			context.TODO(),
			solana.MustPublicKeyFromBase58(marketProgram),
			&rpc.GetProgramAccountsOpts{
				Filters: []rpc.RPCFilter{
					{
						DataSize: 388,
					},
				},
			},
		)
		if err != nil {
			t.Fatal(err)
			return
		}

		for _, market := range allMarketInfo {
			itemMarketInfo := (*MarketStateV3)(unsafe.Pointer(&market.Account.Data.GetBinary()[13]))
			marketAuthority, err := GetAssociatedAuthority(market.Account.Owner.Bytes(), market.Pubkey.Bytes())
			if err != nil {
				t.Fatal(err)
				return
			}

			marketInfo[market.Pubkey.String()] = &MarketInfo{
				MarketProgramId:  market.Account.Owner.String(),
				MarketAuthority:  marketAuthority.String(),
				MarketBaseVault:  solana.PublicKeyFromBytes(itemMarketInfo.BaseVault[:]).String(),
				MarketQuoteVault: solana.PublicKeyFromBytes(itemMarketInfo.QuoteVault[:]).String(),
				MarketBids:       solana.PublicKeyFromBytes(itemMarketInfo.Bids[:]).String(),
				MarketAsks:       solana.PublicKeyFromBytes(itemMarketInfo.Asks[:]).String(),
				MarketEventQueue: solana.PublicKeyFromBytes(itemMarketInfo.EventQueue[:]).String(),
			}
		}
	}
	t.Log("marketInfo:", time.Since(start), "length:", len(marketInfo))
	start = time.Now()

	ammFormatData := make(map[string]*ApiPoolInfoV4)
	for _, itemAmm := range amAccountData {
		itemMarket, ok := marketInfo[solana.PublicKeyFromBytes(itemAmm.LiquidityState.MarketId[:]).String()]
		if !ok {
			continue
		}

		// "amm authority"
		authority, _, err := solana.FindProgramAddress([][]byte{{97, 109, 109, 32, 97, 117, 116, 104, 111, 114, 105, 116, 121}}, itemAmm.ProgramId)
		if err != nil {
			t.Fatal(err)
			return
		}
		ammFormatData[itemAmm.Id.String()] = &ApiPoolInfoV4{
			Id:                 itemAmm.Id.String(),
			BaseMint:           solana.PublicKeyFromBytes(itemAmm.LiquidityState.BaseMint[:]).String(),
			QuoteMint:          solana.PublicKeyFromBytes(itemAmm.LiquidityState.QuoteMint[:]).String(),
			LpMint:             solana.PublicKeyFromBytes(itemAmm.LiquidityState.LpMint[:]).String(),
			BaseDecimals:       itemAmm.LiquidityState.BaseDecimal,
			QuoteDecimals:      itemAmm.LiquidityState.QuoteDecimal,
			LpDecimals:         itemAmm.LiquidityState.BaseDecimal,
			Version:            4,
			ProgramId:          itemAmm.ProgramId.String(),
			Authority:          authority.String(),
			OpenOrders:         solana.PublicKeyFromBytes(itemAmm.LiquidityState.OpenOrders[:]).String(),
			TargetOrders:       solana.PublicKeyFromBytes(itemAmm.LiquidityState.TargetOrders[:]).String(),
			BaseVault:          solana.PublicKeyFromBytes(itemAmm.LiquidityState.BaseVault[:]).String(),
			QuoteVault:         solana.PublicKeyFromBytes(itemAmm.LiquidityState.QuoteVault[:]).String(),
			WithdrawQueue:      solana.PublicKeyFromBytes(itemAmm.LiquidityState.WithdrawQueue[:]).String(),
			LpVault:            solana.PublicKeyFromBytes(itemAmm.LiquidityState.LpVault[:]).String(),
			MarketVersion:      3,
			MarketId:           solana.PublicKeyFromBytes(itemAmm.LiquidityState.MarketId[:]).String(),
			MarketProgramId:    itemMarket.MarketProgramId,
			MarketAuthority:    itemMarket.MarketAuthority,
			MarketBaseVault:    itemMarket.MarketBaseVault,
			MarketQuoteVault:   itemMarket.MarketQuoteVault,
			MarketBids:         itemMarket.MarketBids,
			MarketAsks:         itemMarket.MarketAsks,
			MarketEventQueue:   itemMarket.MarketEventQueue,
			LookupTableAccount: filterDefKey.String(),
		}
	}
	t.Log("ammFormatData:", time.Since(start), "length:", len(ammFormatData))
	start = time.Now()

	ltas, err := client.GetProgramAccountsWithOpts(
		context.TODO(),
		solana.MustPublicKeyFromBase58("AddressLookupTab1e1111111111111111111111111"),
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 22,
						Bytes:  solana.MustPublicKeyFromBase58("RayZuc5vEK174xfgNFdD9YADqbbwbFjVjY4NM8itSF9").Bytes(),
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, itemLTA := range ltas {
		keyStr := itemLTA.Pubkey.String()
		ltaFormat := AddressLookupTableAccount{
			Key:   itemLTA.Pubkey,
			State: NewAddressLookupTableStateFromBytes(itemLTA.Account.Data.GetBinary()),
		}

		for _, itemKey := range ltaFormat.State.Addresses {
			itemKeyStr := itemKey.String()
			if _, ok := ammFormatData[itemKeyStr]; !ok {
				continue
			}
			ammFormatData[itemKeyStr].LookupTableAccount = keyStr
		}
	}
	t.Log("ltas:", time.Since(start), "length:", len(ltas))
	start = time.Now()

	result := make(map[string][]*ApiPoolInfoV4)
	for _, value := range ammFormatData {
		result[value.BaseMint+value.QuoteMint] = append(result[value.BaseMint+value.QuoteMint], value)
		result[value.QuoteMint+value.BaseMint] = append(result[value.QuoteMint+value.BaseMint], value)
	}
	t.Log("result:", time.Since(start), "length:", len(result))

	inputToken := "So11111111111111111111111111111111111111112"
	outputToken := "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"

	poolsList := result[inputToken+outputToken]
	poolsList = append(poolsList, result[outputToken+inputToken]...)
	httpClient := http.DefaultClient
	for {
		request := &RouteBuildRequest{
			InputToken:  inputToken,
			OutputToken: outputToken,
			Amount:      100,
			Slippage:    1,
			PublicKey:   "4d6MQwQC21eXMWToBiTL3UbknXwN3xzZ5Af8EyED4554",
			PoolsList:   poolsList,
		}
		jsonStr, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
			return
		}

		r, err := http.NewRequest("POST", "http://192.168.216.128:8888/route_build", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
			return
		}
		r.Header.Add("Content-Type", "application/json")

		response, err := httpClient.Do(r)
		if err != nil {
			t.Log(err)
			continue
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Log(err)
			continue
		}
		t.Log(body)
	}
}

type AddressLookupTableState struct {
	DeactivationSlot           *big.Int
	LastExtendedSlot           uint64
	LastExtendedSlotStartIndex uint64
	Authority                  solana.PublicKey
	Addresses                  []solana.PublicKey
}

const (
	LOOKUP_TABLE_META_SIZE = 56
)

func NewAddressLookupTableStateFromBytes(data []byte) *AddressLookupTableState {
	if (len(data)-LOOKUP_TABLE_META_SIZE)%32 != 0 {
		panic("invalid data length")
	}

	addresses := make([]solana.PublicKey, 0, (len(data)-LOOKUP_TABLE_META_SIZE)/32)
	for i := LOOKUP_TABLE_META_SIZE; i < len(data); i += 32 {
		addresses = append(addresses, solana.PublicKeyFromBytes(data[i:i+32]))
	}
	return &AddressLookupTableState{
		Addresses: addresses,
	}
}

type AddressLookupTableAccount struct {
	Key   solana.PublicKey
	State *AddressLookupTableState
}

func TestGetAddressLookupTableAccount(t *testing.T) {
	client := rpc.New("https://aged-morning-glade.solana-mainnet.quiknode.pro/b57bbb1a4c8bdd409e1ac53aaedead26da057f59/")

	ltas, err := client.GetProgramAccountsWithOpts(
		context.TODO(),
		solana.MustPublicKeyFromBase58("AddressLookupTab1e1111111111111111111111111"),
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 22,
						Bytes:  solana.MustPublicKeyFromBase58("RayZuc5vEK174xfgNFdD9YADqbbwbFjVjY4NM8itSF9").Bytes(),
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
		return
	}

	lookupTable := make(map[string]string)
	for _, itemLTA := range ltas {
		keyStr := itemLTA.Pubkey.String()
		state := NewAddressLookupTableStateFromBytes(itemLTA.Account.Data.GetBinary())
		for _, itemKey := range state.Addresses {
			itemKeyStr := itemKey.String()
			lookupTable[itemKeyStr] = keyStr
		}
	}
	t.Log(len(lookupTable))
}
