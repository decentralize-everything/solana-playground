package raydium

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"
	"unsafe"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
)

type ApiClmmConfigItem struct {
	Id              string           `json:"id"`
	Index           uint16           `json:"index"`
	ProtocolFeeRate uint32           `json:"protocolFeeRate"`
	TradeFeeRate    uint32           `json:"tradeFeeRate"`
	TickSpacing     uint16           `json:"tickSpacing"`
	FundFeeRate     uint32           `json:"fundFeeRate"`
	FundOwner       solana.PublicKey `json:"fundOwner"`
	Description     string           `json:"description"`
}

type ApiClmmPoolsItemStatistics struct {
	Volume    uint64 `json:"volume"`
	VolumeFee uint64 `json:"volumeFee"`
	FeeA      uint64 `json:"feeA"`
	FeeB      uint64 `json:"feeB"`
	FeeApr    uint64 `json:"feeApr"`
	RewardApr struct {
		A uint64 `json:"a"`
		B uint64 `json:"b"`
		C uint64 `json:"c"`
	} `json:"rewardApr"`
	Apr      uint64 `json:"apr"`
	PriceMin uint64 `json:"priceMin"`
	PriceMax uint64 `json:"priceMax"`
}

type ApiClmmPoolsItem struct {
	Id                 solana.PublicKey
	MintProgramIdA     solana.PublicKey
	MintProgramIdB     solana.PublicKey
	MintA              solana.PublicKey
	MintB              solana.PublicKey
	VaultA             solana.PublicKey
	VaultB             solana.PublicKey
	MintDecimalsA      uint8
	MintDecimalsB      uint8
	AmmConfig          *ApiClmmConfigItem
	RewardInfos        map[solana.PublicKey]solana.PublicKey
	Tvl                uint64
	Day                *ApiClmmPoolsItemStatistics
	Week               *ApiClmmPoolsItemStatistics
	Month              *ApiClmmPoolsItemStatistics
	LookupTableAccount solana.PublicKey
}

type ClmmPoolRewardInfo struct {
	RewardState           uint8            `json:"rewardState"`
	OpenTime              uint64           `json:"openTime"`
	EndTime               uint64           `json:"endTime"`
	LastUpdateTime        uint64           `json:"lastUpdateTime"`
	EmissionsPerSecondX64 *big.Int         `json:"emissionsPerSecondX64"`
	RewardTotalEmissioned uint64           `json:"rewardTotalEmissioned"`
	RewardClaimed         uint64           `json:"rewardClaimed"`
	TokenProgramId        solana.PublicKey `json:"tokenProgramId"`
	TokenMint             solana.PublicKey `json:"tokenMint"`
	TokenVault            solana.PublicKey `json:"tokenVault"`
	Creator               solana.PublicKey `json:"creator"`
	RewardGrowthGlobalX64 *big.Int         `json:"rewardGrowthGlobalX64"`
	PerSecond             decimal.Decimal  `json:"perSecond"`
	RemainingRewards      *big.Int         `json:"remainingRewards"`
}

type Mint struct {
	ProgramId solana.PublicKey `json:"programId"`
	Mint      solana.PublicKey `json:"mint"`
	Vault     solana.PublicKey `json:"vault"`
	Decimals  uint8            `json:"decimals"`
}

type ClmmPoolInfo struct {
	Id                        solana.PublicKey            `json:"id"`
	MintA                     Mint                        `json:"mintA"`
	MintB                     Mint                        `json:"mintB"`
	AmmConfig                 *ApiClmmConfigItem          `json:"ammConfig"`
	ObservationId             solana.PublicKey            `json:"observationId"`
	Creator                   solana.PublicKey            `json:"creator"`
	ProgramId                 solana.PublicKey            `json:"programId"`
	Version                   uint64                      `json:"version"`
	TickSpacing               uint16                      `json:"tickSpacing"`
	Liquidity                 string                      `json:"liquidity"`
	SqrtPriceX64              string                      `json:"sqrtPriceX64"`
	CurrentPrice              *decimal.Decimal            `json:"currentPrice"`
	TickCurrent               int32                       `json:"tickCurrent"`
	ObservationIndex          uint16                      `json:"observationIndex"`
	ObservationUpdateDuration uint16                      `json:"observationUpdateDuration"`
	FeeGrowthGlobalX64A       string                      `json:"feeGrowthGlobalX64A"`
	FeeGrowthGlobalX64B       string                      `json:"feeGrowthGlobalX64B"`
	ProtocolFeesTokenA        uint64                      `json:"protocolFeesTokenA"`
	ProtocolFeesTokenB        uint64                      `json:"protocolFeesTokenB"`
	SwapInAmountTokenA        string                      `json:"swapInAmountTokenA"`
	SwapOutAmountTokenB       string                      `json:"swapOutAmountTokenB"`
	SwapInAmountTokenB        string                      `json:"swapInAmountTokenB"`
	SwapOutAmountTokenA       string                      `json:"swapOutAmountTokenA"`
	TickArrayBitmap           []string                    `json:"tickArrayBitmap"`
	RewardInfos               []*ClmmPoolRewardInfo       `json:"rewardInfos"`
	Day                       *ApiClmmPoolsItemStatistics `json:"day"`
	Week                      *ApiClmmPoolsItemStatistics `json:"week"`
	Month                     *ApiClmmPoolsItemStatistics `json:"month"`
	Tvl                       uint64                      `json:"tvl"`
	LookupTableAccount        solana.PublicKey            `json:"lookupTableAccount"`
	StartTime                 uint64                      `json:"startTime"`
	ExBitmapInfo              *TickArrayBitmapEx          `json:"exBitmapInfo"`
}

type TickArrayBitmapEx struct {
	PoolId                  solana.PublicKey `json:"poolId"`
	PositiveTickArrayBitmap [][]string       `json:"positiveTickArrayBitmap"`
	NegativeTickArrayBitmap [][]string       `json:"negativeTickArrayBitmap"`
}

func FormatClmmKeys(client *rpc.Client) (map[string]*ClmmPoolInfo, error) {
	filterDefKey := solana.MustPublicKeyFromBase58("11111111111111111111111111111111")

	poolAccountInfo, err := client.GetProgramAccountsWithOpts(
		context.TODO(),
		solana.MustPublicKeyFromBase58("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK"),
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					DataSize: 1544,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	poolAccountFormat := make(map[string]*PoolInfoLayout)
	for _, acc := range poolAccountInfo {
		poolAccountFormat[acc.Pubkey.String()] = NewPoolInfoLayoutFromBytes(acc.Account.Data.GetBinary())
	}

	allMint := make(map[solana.PublicKey]struct{})
	for _, pool := range poolAccountFormat {
		if !pool.MintA.Equals(filterDefKey) {
			allMint[pool.MintA] = struct{}{}
		}
		if !pool.MintB.Equals(filterDefKey) {
			allMint[pool.MintB] = struct{}{}
		}
		for _, rewardInfo := range pool.RewardInfos {
			if !rewardInfo.TokenMint.Equals(filterDefKey) {
				allMint[rewardInfo.TokenMint] = struct{}{}
			}
		}
	}
	var allMintAccount []solana.PublicKey
	for mint := range allMint {
		allMintAccount = append(allMintAccount, mint)
	}

	mintInfoDict := make(map[solana.PublicKey]*rpc.Account)
	accounts, err := getMultipleAccountsInfo(client, allMintAccount)
	if err != nil {
		return nil, err
	}
	for i, acc := range accounts {
		mintInfoDict[allMintAccount[i]] = acc
	}

	configAccountInfo, err := client.GetProgramAccountsWithOpts(
		context.TODO(),
		solana.MustPublicKeyFromBase58("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK"),
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					DataSize: 117,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	configIdToData := make(map[string]*ApiClmmConfigItem)
	for _, acc := range configAccountInfo {
		data := acc.Account.Data.GetBinary()
		configIdToData[acc.Pubkey.String()] = &ApiClmmConfigItem{
			Id:              acc.Pubkey.String(),
			Index:           *(*uint16)(unsafe.Pointer(&data[9])),
			ProtocolFeeRate: *(*uint32)(unsafe.Pointer(&data[43])),
			TradeFeeRate:    *(*uint32)(unsafe.Pointer(&data[47])),
			TickSpacing:     *(*uint16)(unsafe.Pointer(&data[51])),
			FundFeeRate:     *(*uint32)(unsafe.Pointer(&data[53])),
			FundOwner:       solana.PublicKeyFromBytes(data[61:93]),
			Description:     "",
		}
	}

	poolInfoDict := map[string]*ApiClmmPoolsItem{}
	for id, pool := range poolAccountFormat {
		mintProgramIdA := mintInfoDict[pool.MintA].Owner
		mintProgramIdB := mintInfoDict[pool.MintB].Owner
		rewardInfos := make(map[solana.PublicKey]solana.PublicKey)
		for _, rewardInfo := range pool.RewardInfos {
			if rewardInfo.TokenMint.Equals(filterDefKey) {
				continue
			}
			rewardInfos[rewardInfo.TokenMint] = mintInfoDict[rewardInfo.TokenMint].Owner
		}

		poolInfoDict[id] = &ApiClmmPoolsItem{
			Id:                 solana.MustPublicKeyFromBase58(id),
			MintProgramIdA:     mintProgramIdA,
			MintProgramIdB:     mintProgramIdB,
			MintA:              pool.MintA,
			MintB:              pool.MintB,
			VaultA:             pool.VaultA,
			VaultB:             pool.VaultB,
			MintDecimalsA:      pool.MintDecimalsA,
			MintDecimalsB:      pool.MintDicimalsB,
			AmmConfig:          configIdToData[pool.AmmConfig.String()],
			RewardInfos:        rewardInfos,
			Tvl:                0,
			Day:                &ApiClmmPoolsItemStatistics{},
			Week:               &ApiClmmPoolsItemStatistics{},
			Month:              &ApiClmmPoolsItemStatistics{},
			LookupTableAccount: filterDefKey,
		}
	}

	poolKeys := make([]*ApiClmmPoolsItem, 0, len(poolInfoDict))
	for _, pool := range poolInfoDict {
		poolKeys = append(poolKeys, pool)
	}

	publicKeys := make([]solana.PublicKey, 0, len(poolKeys))
	for _, pool := range poolKeys {
		publicKeys = append(publicKeys, pool.Id)
	}
	poolAccountInfos, err := getMultipleAccountsInfo(client, publicKeys)
	if err != nil {
		return nil, err
	}

	exBitmapAddress := make(map[solana.PublicKey]solana.PublicKey)
	for i, apiPoolInfo := range poolKeys {
		accountInfo := poolAccountInfos[i]
		if accountInfo == nil {
			continue
		}

		exBitmapAddress[apiPoolInfo.Id] = getPdaExBitmapAccount(accountInfo.Owner, apiPoolInfo.Id)
	}

	exBitmapAddressValues := make([]solana.PublicKey, 0, len(exBitmapAddress))
	for _, value := range exBitmapAddress {
		exBitmapAddressValues = append(exBitmapAddressValues, value)
	}

	fetchedBitmapAccount, err := getMultipleAccountsInfo(client, exBitmapAddressValues)
	if err != nil {
		return nil, err
	}

	exBitmapAccountInfos := make(map[solana.PublicKey]*TickArrayBitmap)
	for i, acc := range fetchedBitmapAccount {
		if acc == nil {
			continue
		}
		exBitmapAccountInfos[exBitmapAddressValues[i]] = NewTickArrayBitmapFromBytes(acc.Data.GetBinary())
	}

	// var programIds []solana.PublicKey
	poolsInfo := make(map[string]*ClmmPoolInfo)
	// var updateRewardInfos []ClmmPoolRewardInfo
	for i, apiPoolInfo := range poolKeys {
		accountInfo := poolAccountInfos[i]
		exBitmapInfo := exBitmapAccountInfos[exBitmapAddress[apiPoolInfo.Id]]
		if accountInfo == nil {
			continue
		}

		layoutAccountInfo := NewPoolInfoLayoutFromBytes(accountInfo.Data.GetBinary())
		poolsInfo[apiPoolInfo.Id.String()] = &ClmmPoolInfo{
			Id: apiPoolInfo.Id,
			MintA: Mint{
				ProgramId: apiPoolInfo.MintProgramIdA,
				Mint:      layoutAccountInfo.MintA,
				Vault:     layoutAccountInfo.VaultA,
				Decimals:  layoutAccountInfo.MintDecimalsA,
			},
			MintB: Mint{
				ProgramId: apiPoolInfo.MintProgramIdB,
				Mint:      layoutAccountInfo.MintB,
				Vault:     layoutAccountInfo.VaultB,
				Decimals:  layoutAccountInfo.MintDicimalsB,
			},
			ObservationId: layoutAccountInfo.ObservationId,
			AmmConfig:     apiPoolInfo.AmmConfig,
			Creator:       layoutAccountInfo.Creator,
			ProgramId:     accountInfo.Owner,
			Version:       6,
			TickSpacing:   layoutAccountInfo.TickSpacing,
			Liquidity:     layoutAccountInfo.Liquidity.String(),
			SqrtPriceX64:  layoutAccountInfo.SqrtPriceX64.String(),
			CurrentPrice: sqrtPriceX64ToPrice(
				layoutAccountInfo.SqrtPriceX64,
				int64(layoutAccountInfo.MintDecimalsA),
				int64(layoutAccountInfo.MintDicimalsB),
			),
			TickCurrent:               layoutAccountInfo.TickCurrent,
			ObservationIndex:          layoutAccountInfo.ObservationIndex,
			ObservationUpdateDuration: layoutAccountInfo.ObservationUpdateDuration,
			FeeGrowthGlobalX64A:       layoutAccountInfo.FeeGrowthGlobalX64A.String(),
			FeeGrowthGlobalX64B:       layoutAccountInfo.FeeGrowthGlobalX64B.String(),
			ProtocolFeesTokenA:        layoutAccountInfo.ProtocolFeesTokenA,
			ProtocolFeesTokenB:        layoutAccountInfo.ProtocolFeesTokenB,
			SwapInAmountTokenA:        layoutAccountInfo.SwapInAmountTokenA.String(),
			SwapOutAmountTokenB:       layoutAccountInfo.SwapOutAmountTokenB.String(),
			SwapInAmountTokenB:        layoutAccountInfo.SwapInAmountTokenB.String(),
			SwapOutAmountTokenA:       layoutAccountInfo.SwapOutAmountTokenA.String(),
			TickArrayBitmap:           toStringArray(layoutAccountInfo.TickArrayBitmap),
			RewardInfos: updatePoolRewardInfos(
				client,
				apiPoolInfo,
				uint64(time.Now().Unix()),
				layoutAccountInfo.Liquidity,
				layoutAccountInfo.RewardInfos,
				filterDefKey,
			),
			Day:                apiPoolInfo.Day,
			Week:               apiPoolInfo.Week,
			Month:              apiPoolInfo.Month,
			Tvl:                apiPoolInfo.Tvl,
			LookupTableAccount: apiPoolInfo.LookupTableAccount,
			StartTime:          layoutAccountInfo.StartTime,
			ExBitmapInfo:       toStringMatrix(exBitmapInfo),
		}
	}

	return poolsInfo, nil
}

func getMultipleAccountsInfo(client *rpc.Client, publicKeys []solana.PublicKey) ([]*rpc.Account, error) {
	chunkedKeys := make([][]solana.PublicKey, 0)
	chunkSize := 100
	for i := 0; i < len(publicKeys); i += chunkSize {
		end := i + chunkSize
		if end > len(publicKeys) {
			end = len(publicKeys)
		}
		chunkedKeys = append(chunkedKeys, publicKeys[i:end])
	}

	accounts := make([]*rpc.Account, 0, len(publicKeys))
	for _, chunk := range chunkedKeys {
		result, err := client.GetMultipleAccountsWithOpts(context.TODO(), chunk, &rpc.GetMultipleAccountsOpts{
			Encoding: solana.EncodingBase64,
		})
		if err != nil {
			return nil, err
		}

		for _, acc := range result.Value {
			if acc == nil {
				return nil, errors.New("account not found")
			}
			accounts = append(accounts, acc)
		}
	}
	return accounts, nil
}

func getPdaExBitmapAccount(programId, poolId solana.PublicKey) solana.PublicKey {
	publicKey, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("pool_tick_array_bitmap_extension"),
			poolId.Bytes(),
		},
		programId,
	)
	return publicKey
}

func sqrtPriceX64ToPrice(sqrtPriceX64 *big.Int, decimalsA int64, decimalsB int64) *decimal.Decimal {
	d := decimal.NewFromBigInt(sqrtPriceX64, 0).Pow(decimal.NewFromInt(2)).Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(decimalsA - decimalsB)))
	return &d
}

func updatePoolRewardInfos(
	client *rpc.Client,
	apiPoolInfo *ApiClmmPoolsItem,
	chainTime uint64,
	poolLiquidity *big.Int,
	rewardInfos []*RewardInfo,
	filterDefKey solana.PublicKey,
) []*ClmmPoolRewardInfo {
	nRewardInfo := []*ClmmPoolRewardInfo{}
	for _, rewardInfo := range rewardInfos {
		if rewardInfo.TokenMint.Equals(filterDefKey) {
			continue
		}

		accInfo, err := client.GetAccountInfo(context.TODO(), rewardInfo.TokenMint)
		if err != nil {
			// TODO
			continue
		}
		apiRewardProgram := accInfo.Value.Owner

		itemReward := &ClmmPoolRewardInfo{
			RewardState:           rewardInfo.RewardState,
			OpenTime:              rewardInfo.OpenTime,
			EndTime:               rewardInfo.EndTime,
			LastUpdateTime:        rewardInfo.LastUpdateTime,
			EmissionsPerSecondX64: rewardInfo.EmissionsPerSecondX64,
			RewardTotalEmissioned: rewardInfo.RewardTotalEmissioned,
			RewardClaimed:         rewardInfo.RewardClaimed,
			TokenProgramId:        apiRewardProgram,
			TokenMint:             rewardInfo.TokenMint,
			TokenVault:            rewardInfo.TokenVault,
			Creator:               rewardInfo.Creator,
			RewardGrowthGlobalX64: rewardInfo.RewardGrowthGlobalX64,
			PerSecond:             decimal.NewFromBigInt(rewardInfo.EmissionsPerSecondX64, 0),
			RemainingRewards:      nil,
		}
		if chainTime <= itemReward.OpenTime || poolLiquidity == nil || poolLiquidity.Int64() == 0 {
			nRewardInfo = append(nRewardInfo, itemReward)
			continue
		}

		latestUpdateTime := itemReward.EndTime
		if chainTime < latestUpdateTime {
			latestUpdateTime = chainTime
		}
		timeDelta := int64(latestUpdateTime) - int64(itemReward.LastUpdateTime)
		rewardGrowthDeltaX64 := new(big.Int).Div(new(big.Int).Mul(big.NewInt(int64(timeDelta)), itemReward.EmissionsPerSecondX64), poolLiquidity)
		rewardEmissionedDelta := new(big.Int).Div(new(big.Int).Mul(big.NewInt(int64(timeDelta)), itemReward.EmissionsPerSecondX64), new(big.Int).Lsh(big.NewInt(1), 64))
		rewardTotalEmissioned := itemReward.RewardTotalEmissioned + rewardEmissionedDelta.Uint64()
		itemReward.RewardGrowthGlobalX64 = rewardGrowthDeltaX64
		itemReward.RewardTotalEmissioned = rewardTotalEmissioned
		itemReward.LastUpdateTime = latestUpdateTime
		nRewardInfo = append(nRewardInfo, itemReward)
	}
	return nRewardInfo
}

func toStringArray(array []uint64) []string {
	var strArray []string
	for _, value := range array {
		if int64(value) < 0 {
			fmt.Println("value is negative")
		}
		strArray = append(strArray, strconv.FormatUint(value, 10))
	}
	return strArray
}

func toStringMatrix(bitmap *TickArrayBitmap) *TickArrayBitmapEx {
	if bitmap == nil {
		return nil
	}

	p := make([][]string, 0, len(bitmap.PositiveTickArrayBitmap))
	for _, value := range bitmap.PositiveTickArrayBitmap {
		p = append(p, toStringArray(value))
	}
	n := make([][]string, 0, len(bitmap.NegativeTickArrayBitmap))
	for _, value := range bitmap.NegativeTickArrayBitmap {
		n = append(n, toStringArray(value))
	}
	return &TickArrayBitmapEx{
		PoolId:                  bitmap.PoolId,
		PositiveTickArrayBitmap: p,
		NegativeTickArrayBitmap: n,
	}
}
