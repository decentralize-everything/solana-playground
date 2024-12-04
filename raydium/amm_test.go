package raydium

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

var (
	WSOL_MINT                   = solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	SYSVAR_RENT_PUBKEY          = solana.MustPublicKeyFromBase58("SysvarRent111111111111111111111111111111111")
	ASSOCIATED_TOKEN_PROGRAM_ID = solana.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
	SYSTEM_PROGRAM_ID           = solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
)

func TestGetAmmInfo(t *testing.T) {

	key := solana.MustPublicKeyFromBase58("2immgwYNHBbyVQKVGCEkgWpi53bLwWNRMB5G2nbgYV17")
	t.Log(key)

	data, err := base64.StdEncoding.DecodeString("AQAAAP//////////d49+DAAAAAAAAQZMWvw7GUNJdaccNBVnb57OKakxL2BHLYvhRwVILRsgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMGRm/lIRcy/+ytunLDm+e8jOW7xfcSayxDmzpAAAAABt324ddloZPZy+FGzut5rBy0he1fWzeROoz1hX7/AKkG3fbh7nWP3hhCXbzkbM3athr8TYO5DSf+vfko2KGL/AVKU1D4XciC1hSlVnJ4iilt3x6rq9CmBniISTL07vagBqfVFxksXFEhjMlMPUrxf1ja7gibof1E49vZigAAAAAGp9UXGMd0yShWY5hpHV62i164o5tLbVxzVVshAAAAAIyXJY9OJInxuz0QKRSODYMLWhOZ2v8QhASOe9jb6fhZC3BlsePRfEU4nVJ/awTDzVi4bHMaoP21SbbRvAP4KUbIScv+6Yw2LHF/6K0ZjUPibbSWXCirYPGuuVl7zT789IUPLW4CpHr4JNCatp3ELXDLKMv6JJ+37le50lbBJ2LvDQdRqCgtphMF/imcN7mY5YRx2xE1A3MQ+L4QRaYK9u4GRfZP3LsAd00a+IkCpA22UNQMKdq5BFbJuwuOLqc8zxCTDlqxBG8J0HcxtfogQHDK06ukzfaXiNDKAob1MqBHS9lJxDYCwz8gd5DtFqNSTKG5l1zxIaKpDP/sffi2is1H9aKveyXSu5StXElYRl9SD5As0DHE4N0GLnf84/siiKXVyp4Ez121kLcUui/jLLFZEz/BwZK3Ilf9B9OcsEAeDMKAy2vjGSxQODgBz0QwGA+eP4ZjIjrIAQaXENv31QfLlOdXSRCkaybRniDHF4C8YcwhcvsqrOVuTP4B2Na+9wLdtrB31uz2rtlFI5kahdsnp/d1SrASDInYCtTYtdoke4kX+hoKWcEWM4Tle8pTUkUVv4BxS6fje/EzKBE4Qu9N9LMnrw/JNO0hqMVB4rk/2ou4AB1loQ7FZoPwut2o4KZB+0p9xnbrQKw038qjpHar+PyDwvxBRcu5hpHw3dguezeWv+IwvgW5icu8EGkhGa9AkFPPJT7VMSFb8xowveU=")
	if err != nil {
		t.Error(err)
	}

	lookupTableState, err := addresslookuptable.DecodeAddressLookupTableState(data)
	if err != nil {
		t.Error(err)
	}

	client := rpc.New("https://aged-morning-glade.solana-mainnet.quiknode.pro/b57bbb1a4c8bdd409e1ac53aaedead26da057f59/")

	// ammInfo, err := GetAmmInfo(client, "EVzLJhqMtdC1nPmz8rNd6xGfVjDPxpLZgq7XJuNfMZ6", solana.PublicKey{})
	ammInfo, err := GetAmmInfo(client, "AVs9TA4nWDzfPJE9gGVNJMVhcQy3V9PGazuz33BfG2RA", solana.PublicKey{})
	if err != nil {
		t.Error(err)
	}
	t.Log(ammInfo.Display())

	start := time.Now()

	recentBlockhashResult, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentConfirmed)
	if err != nil {
		t.Error(err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			solana.NewInstruction(ammInfo.ProgramId, []*solana.AccountMeta{
				{PublicKey: ammInfo.Id},
				{PublicKey: ammInfo.Authority},
				{PublicKey: ammInfo.OpenOrders},
				{PublicKey: ammInfo.BaseVault},
				{PublicKey: ammInfo.QuoteVault},
				{PublicKey: ammInfo.LpMint},
				{PublicKey: ammInfo.MarketId},
				{PublicKey: ammInfo.MarketEventQueue},
			}, []byte{12, 0}),
		},
		recentBlockhashResult.Value.Blockhash,
		solana.TransactionPayer(solana.MustPublicKeyFromBase58("RaydiumSimuLateTransaction11111111111111111")),
	)
	tx.Signatures = append(tx.Signatures, solana.Signature{})
	if err != nil {
		t.Error(err)
	}

	simulateTransactionResponse, err := client.SimulateTransactionWithOpts(context.TODO(), tx, &rpc.SimulateTransactionOpts{ReplaceRecentBlockhash: true})
	if err != nil {
		t.Error(err)
	}

	var logs []string
	for _, log := range simulateTransactionResponse.Value.Logs {
		if !strings.Contains(log, "GetPoolData") {
			continue
		}
		logs = append(logs, log)
	}
	t.Log(logs)

	t.Log("Fetch info: ", time.Since(start))

	var poolInfos []*PoolInfo
	for _, log := range logs {
		jsonStr := log[strings.Index(log, "{") : strings.LastIndex(log, "}")+1]
		var poolInfo PoolInfo
		if err := json.Unmarshal([]byte(jsonStr), &poolInfo); err != nil {
			t.Error(err)
		}
		t.Log(poolInfo)
		poolInfos = append(poolInfos, &poolInfo)
	}
	t.Log(poolInfos)

	inputTokenMint := WSOL_MINT
	outputTokenMint := solana.MustPublicKeyFromBase58("4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R")
	userAccount := solana.MustPublicKeyFromBase58("4d6MQwQC21eXMWToBiTL3UbknXwN3xzZ5Af8EyED4554")
	amountIn := big.NewInt(1000)

	amountOut, minAmountOut := ComputeAmountOut(
		ammInfo,
		poolInfos[0],
		&Token{
			ProgramId: TOKEN_PROGRAM_ID,
			Mint:      inputTokenMint,
			Decimals:  9,
			Symbol:    "WSOL",
			Name:      "WSOL",
		},
		&Token{
			ProgramId: TOKEN_PROGRAM_ID,
			Mint:      outputTokenMint,
			Decimals:  6,
			Symbol:    "RAY",
			Name:      "RAY",
		},
		amountIn,
		1,
	)

	t.Log(amountOut, minAmountOut)

	tokenAccountsResult, err := client.GetTokenAccountsByOwner(
		context.TODO(),
		userAccount,
		&rpc.GetTokenAccountsConfig{
			ProgramId: &TOKEN_PROGRAM_ID,
		},
		nil,
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(tokenAccountsResult)

	var tokenAccounts []*TokenAccount
	for _, account := range tokenAccountsResult.Value {
		tokenAccounts = append(tokenAccounts, &TokenAccount{
			PublicKey:   account.Pubkey,
			ProgramId:   account.Account.Owner,
			AccountInfo: NewSplAccountFromBytes(account.Account.Data.GetBinary()),
		})
	}
	t.Log(tokenAccounts)

	tokenAccountIn := selectTokenAccount(tokenAccounts, inputTokenMint, userAccount, false)
	if tokenAccountIn != nil {
		t.Log(tokenAccountIn)
	}

	tokenAccountOut := selectTokenAccount(tokenAccounts, outputTokenMint, userAccount, true)
	if tokenAccountOut != nil {
		t.Log(tokenAccountOut)
	}

	tokenInPubKey, frontInstructionsIn, endInstructionsIn := handleTokenAccount(client, tokenAccountIn, "in", amountIn, inputTokenMint, userAccount, TOKEN_PROGRAM_ID)
	tokenOutPubKey, frontInstructionsOut, endInstructionsOut := handleTokenAccount(client, tokenAccountOut, "out", big.NewInt(0), outputTokenMint, userAccount, TOKEN_PROGRAM_ID)
	swapInstruction := makeSwapInstruction(ammInfo, tokenInPubKey, tokenOutPubKey, userAccount, amountIn, minAmountOut, "in")
	t.Log("tokenInPubKey", tokenInPubKey, "tokenOutPubKey", tokenOutPubKey)
	var instructions []solana.Instruction
	instructions = append(instructions, computebudget.NewSetComputeUnitLimitInstruction(10000000).Build())
	instructions = append(instructions, computebudget.NewSetComputeUnitPriceInstruction(10).Build())
	instructions = append(instructions, frontInstructionsIn...)
	instructions = append(instructions, frontInstructionsOut...)
	instructions = append(instructions, swapInstruction)
	instructions = append(instructions, endInstructionsIn...)
	instructions = append(instructions, endInstructionsOut...)
	tx, err = solana.NewTransaction(
		instructions,
		recentBlockhashResult.Value.Blockhash,
		solana.TransactionPayer(userAccount),
		solana.TransactionAddressTables(map[solana.PublicKey]solana.PublicKeySlice{
			key: lookupTableState.Addresses,
		}),
	)
	if err != nil {
		t.Error(err)
	}
	tx.Signatures = append(tx.Signatures, solana.Signature{})
	// "Create: address CwQcFzbLDdo2GLPDq1KVs6pFBdRv2wN8Bbxsn64mUwGP does not match derived address HiB5NeyoFbcvwsqyYiiQsdT5ZvdnwXdcwYfeQXH7Nxwp"
	simulateTransactionResponse, err = client.SimulateTransactionWithOpts(context.TODO(), tx, &rpc.SimulateTransactionOpts{ReplaceRecentBlockhash: true})
	if err != nil {
		t.Error(err)
	}
	t.Log(simulateTransactionResponse.Value.Logs)
}

func makeSwapInstruction(ammInfo *AmmInfo, tokenInPubKey, tokenOutPubKey, owner solana.PublicKey, amountIn, amountOut *big.Int, fixedSide string) solana.Instruction {
	data := make([]byte, 1+8+8)
	data[0] = 9
	binary.LittleEndian.PutUint64(data[1:], amountIn.Uint64())
	binary.LittleEndian.PutUint64(data[9:], amountOut.Uint64())

	return solana.NewInstruction(
		ammInfo.ProgramId,
		[]*solana.AccountMeta{
			{PublicKey: TOKEN_PROGRAM_ID},
			{PublicKey: ammInfo.Id, IsWritable: true},
			{PublicKey: ammInfo.Authority},
			{PublicKey: ammInfo.OpenOrders, IsWritable: true},
			{PublicKey: ammInfo.TargetOrders, IsWritable: true},
			{PublicKey: ammInfo.BaseVault, IsWritable: true},
			{PublicKey: ammInfo.QuoteVault, IsWritable: true},
			{PublicKey: ammInfo.MarketProgramId},
			{PublicKey: ammInfo.MarketId, IsWritable: true},
			{PublicKey: ammInfo.MarketBids, IsWritable: true},
			{PublicKey: ammInfo.MarketAsks, IsWritable: true},
			{PublicKey: ammInfo.MarketEventQueue, IsWritable: true},
			{PublicKey: ammInfo.MarketBaseVault, IsWritable: true},
			{PublicKey: ammInfo.MarketQuoteVault, IsWritable: true},
			{PublicKey: ammInfo.MarketAuthority},
			{PublicKey: tokenInPubKey, IsWritable: true},
			{PublicKey: tokenOutPubKey, IsWritable: true},
			{PublicKey: owner, IsSigner: true, IsWritable: true},
		},
		data,
	)
}

func selectTokenAccount(tokenAccounts []*TokenAccount, mint solana.PublicKey, owner solana.PublicKey, associatedOnly bool) *TokenAccount {
	var nonEmptyAccounts []*TokenAccount
	for _, account := range tokenAccounts {
		if account.AccountInfo.Amount > 0 {
			nonEmptyAccounts = append(nonEmptyAccounts, account)
			fmt.Println(account)
		}
	}

	ata, _, err := solana.FindProgramAddress([][]byte{owner.Bytes(), TOKEN_PROGRAM_ID[:], mint.Bytes()}, solana.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"))
	if err != nil {
		panic(err)
	}

	for _, account := range nonEmptyAccounts {
		if solana.PublicKeyFromBytes(account.AccountInfo.Mint[:]).Equals(mint) {
			if associatedOnly {
				if ata.Equals(account.PublicKey) {
					return account
				}
			} else {
				return account
			}
		}
	}
	return nil
}

func handleTokenAccount(client *rpc.Client, tokenAccount *TokenAccount, side string, amount *big.Int, mint solana.PublicKey, owner solana.PublicKey, programId solana.PublicKey) (solana.PublicKey, []solana.Instruction, []solana.Instruction) {
	ata, _, err := solana.FindProgramAddress([][]byte{owner.Bytes(), TOKEN_PROGRAM_ID[:], mint.Bytes()}, solana.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"))
	if err != nil {
		panic(err)
	}

	if WSOL_MINT.Equals(mint) {
		newTokenAccount, frontInstructions := insertCreateWrappedNativeAccount(client, amount, owner, mint, programId)
		endInstructions := []solana.Instruction{createCloseAccountInstruction(newTokenAccount, owner, owner, programId)}
		return newTokenAccount, frontInstructions, endInstructions
	} else if tokenAccount == nil || (side == "out" && !ata.Equals(tokenAccount.PublicKey)) {
		instruction := makeCreateAssociatedTokenAccountInstruction(owner, ata, mint, programId)
		return ata, []solana.Instruction{instruction}, nil
	}

	return tokenAccount.PublicKey, nil, nil
}

func makeCreateAssociatedTokenAccountInstruction(owner, associatedTokenAccount, mint, programId solana.PublicKey) solana.Instruction {
	return solana.NewInstruction(
		ASSOCIATED_TOKEN_PROGRAM_ID,
		[]*solana.AccountMeta{
			{PublicKey: owner, IsSigner: true, IsWritable: true},
			{PublicKey: associatedTokenAccount, IsWritable: true},
			{PublicKey: owner},
			{PublicKey: mint},
			{PublicKey: SYSTEM_PROGRAM_ID},
			{PublicKey: programId},
		},
		[]byte{},
	)
}

func createCloseAccountInstruction(account, destination, authrity, programId solana.PublicKey) solana.Instruction {
	return solana.NewInstruction(
		programId,
		[]*solana.AccountMeta{
			{PublicKey: account, IsWritable: true},
			{PublicKey: destination, IsWritable: true},
			{PublicKey: authrity, IsSigner: true, IsWritable: true},
		},
		[]byte{9},
	)
}

func insertCreateWrappedNativeAccount(client *rpc.Client, amount *big.Int, owner, mint, programId solana.PublicKey) (solana.PublicKey, []solana.Instruction) {
	return makeCreateWrappedNativeAccountInstructions(client, amount, owner, mint, programId)
}

func makeCreateWrappedNativeAccountInstructions(client *rpc.Client, amount *big.Int, owner, mint, programId solana.PublicKey) (solana.PublicKey, []solana.Instruction) {
	balanceNeeded, err := getMinimumBalanceForRentExemption(client)
	if err != nil {
		panic(err)
	}

	var instructions []solana.Instruction
	lamports := new(big.Int).Add(amount, big.NewInt(int64(balanceNeeded)))
	fmt.Println(owner, TOKEN_PROGRAM_ID)
	newAccount, seed := generatePubKey(owner, TOKEN_PROGRAM_ID)
	instructions = append(instructions, system.NewCreateAccountWithSeedInstruction(
		owner, // TODO
		seed,
		lamports.Uint64(),
		uint64(SPL_ACCOUNT_SIZE),
		TOKEN_PROGRAM_ID,
		owner,
		newAccount,
		owner, // TODO
	).Build())
	instructions = append(instructions, createInitializeAccountInstruction(newAccount, mint, owner, programId))

	return newAccount, instructions
}

func createInitializeAccountInstruction(account, mint, owner, programId solana.PublicKey) solana.Instruction {
	return solana.NewInstruction(
		TOKEN_PROGRAM_ID,
		[]*solana.AccountMeta{
			{PublicKey: account, IsWritable: true},
			{PublicKey: mint},
			{PublicKey: owner},
			{PublicKey: SYSVAR_RENT_PUBKEY},
		},
		[]byte{1},
	)
}

func generatePubKey(fromPublicKey solana.PublicKey, programId solana.PublicKey) (solana.PublicKey, string) {
	seed := solana.NewWallet().PublicKey().String()[0:32]
	var buffer []byte
	buffer = append(buffer, fromPublicKey.Bytes()...)
	buffer = append(buffer, []byte(seed)...)
	buffer = append(buffer, programId.Bytes()...)
	publicKeyBytes := sha256.Sum256(buffer)
	return solana.PublicKeyFromBytes(publicKeyBytes[:]), seed
}

func getMinimumBalanceForRentExemption(client *rpc.Client) (uint64, error) {
	return client.GetMinimumBalanceForRentExemption(context.TODO(), uint64(SPL_ACCOUNT_SIZE), rpc.CommitmentConfirmed)
}

// market id 56ZNe9c73XrizrXXqzPd9xdjNPGpxsWbiV4RszFKfBL8
// market event queue 7Yb9UPpS6ykpFrWjrNRi37JUzynCbs4BtqHQTw8g4gfk
// 2CoBP2rr5HmjMdPC4nMwnYg1cdH9JPUuqbq2QGSMGfms
func TestDecodeTransaction(t *testing.T) {
	encoded := "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAkKBkxbUEfMBSgVDsnH6oknAyyoxjw99OknHjB1lSbiAAAR4UphPbxYtPYvPfqX/EzZd+3TfQqtA0vVv+4duOPf+h3mL8NBPWAS/89boM7VX/KOlVSTp5sKW7W88GIajL5oI25L98mLi31X7oX8JRQ6mpkYXGipnu2NI9cRqT0b83Y820GSuWpFiG3hKdfJho5uqmPZzHr6KvcNbzTv1myyzUFXsFgPMcX85EpiWC28+deO51lDoISjk7NQNo0iiZMIS9lJxDYCwz8gd5DtFqNSTKG5l1zxIaKpDP/sffi2is0DdVPyVZQlfze4p0MwF83mZGBavX22USFPChIaKs6+FdP++SJgq9ALC2Sj+3kMDXg5yIDAWBlwCAEOcTDdw5at1fIgZ9U9Kg9/8FPJpiK24qa3bsmOzrAqLIi7vGb3IRoKVLMpU4xtpxwTz75c/hdhLsUs3QvSDp2kHWGQYyysJwEGCAcFAwkCCAQBAgwA"
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Error(err)
	}

	// parse transaction:
	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(data))
	if err != nil {
		t.Error(err)
	}

	client := rpc.New("https://aged-morning-glade.solana-mainnet.quiknode.pro/b57bbb1a4c8bdd409e1ac53aaedead26da057f59/")
	recentBlockhashResult, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentConfirmed)
	if err != nil {
		t.Error(err)
	}
	tx.Message.RecentBlockhash = recentBlockhashResult.Value.Blockhash

	simulateTransactionResponse, err := client.SimulateTransactionWithOpts(context.TODO(), tx, &rpc.SimulateTransactionOpts{ReplaceRecentBlockhash: true})
	if err != nil {
		t.Error(err)
	}

	t.Log(simulateTransactionResponse)

	t.Log(tx.Message.AccountKeys[1].String())
	t.Log(tx)
}

// AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAQAGFDXSBXhYPmoBRu5x0tSE4Dpjl6e3wotK2XljAhVXLkdTVxlKRrCVxrqjSnTGQIuTxkmIz8Oqab8y3adCv0fHzjmNHiUZn+aC9/yLVjGEXv4IxkgQrcoQ3mZByijavkoiFVDs2n3OBDb+pbbsf+iadatphbqjzsmfRgoRiu8VJW2LRbis01bXHbBS0QLGzc9LmMkjVqnNnY6QYOvk+cz51L3MdBmbRbrPmC02Xka/72oRoXqdx1MeFypn4krmb3GtOCkL67h9c54CiRN7OsEfwdc5APCimsW2ZmBbZR0fP6iRpPMK+KyZYVVWsVJclPiChTagQVyUjjWUvlWIklITutGjo7PoiVhbaK5w4bw3aNV0TY9kN7VfA5UshJKSVT71OC+9+4oe4q57oouTNmX7547mrejUmeE6cRg1YRpZqVMKML5cHBZTQqjoKYcuFdBu2dLzAe337bdHAcF1VrzZhsxRO0w5kdF8SZcPU3XNV5mmnryNnZnmg/QEdqo8ticZbjST1mrakRIDDpi7woZhsZK3ZihaIdnwcVmOQHgWFa/1c+13V/0oVbYwfpAIS69/AL1CEs2cLpEkwFQ9N18RsNEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAbd9uHXZaGT2cvhRs7reawctIXtX1s3kTqM9YV+/wCpBpuIV/6rgYT7aH9jRhjANdrEOdwa6ztVmKDwAAAAAAFL2UnENgLDPyB3kO0Wo1JMobmXXPEhoqkM/+x9+LaKzUFXsFgPMcX85EpiWC28+deO51lDoISjk7NQNo0iiZMIX7cor9zcBHh47oQ6b+tEryjNU6NIKNmQO1rV5pcZVSoU/aWdG4DckMQPqfcAOm8TRoNxPqtn2SaMF4IoXCW5gAQOAgABfAMAAAA10gV4WD5qAUbucdLUhOA6Y5ent8KLStl5YwIVVy5HUyAAAAAAAAAAajFNMW1Dc1RKZW1kRkVaaFpvNVlkU3BldFk4bnVqWXHYIR8AAAAAAKUAAAAAAAAABt324ddloZPZy+FGzut5rBy0he1fWzeROoz1hX7/AKkPBAEQABQBARESDwISAwQFBhUHCAkKCwwTAQ0AEQnoAwAAAAAAAFQAAAAAAAAADwMBAAABCQEZjx9MOkUiY9QTss0X68vBoOWIc2TmJhoSqBeS6hZaPgACBQo=

func TestDecodeTransaction2(t *testing.T) {
	encoded := "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAQAGFDXSBXhYPmoBRu5x0tSE4Dpjl6e3wotK2XljAhVXLkdTVxlKRrCVxrqjSnTGQIuTxkmIz8Oqab8y3adCv0fHzjmNHiUZn+aC9/yLVjGEXv4IxkgQrcoQ3mZByijavkoiFVDs2n3OBDb+pbbsf+iadatphbqjzsmfRgoRiu8VJW2LRbis01bXHbBS0QLGzc9LmMkjVqnNnY6QYOvk+cz51L3MdBmbRbrPmC02Xka/72oRoXqdx1MeFypn4krmb3GtOCkL67h9c54CiRN7OsEfwdc5APCimsW2ZmBbZR0fP6iRpPMK+KyZYVVWsVJclPiChTagQVyUjjWUvlWIklITutGjo7PoiVhbaK5w4bw3aNV0TY9kN7VfA5UshJKSVT71OC+9+4oe4q57oouTNmX7547mrejUmeE6cRg1YRpZqVMKML5cHBZTQqjoKYcuFdBu2dLzAe337bdHAcF1VrzZhsxRO0w5kdF8SZcPU3XNV5mmnryNnZnmg/QEdqo8ticZbjST1mrakRIDDpi7woZhsZK3ZihaIdnwcVmOQHgWFa/1c+13V/0oVbYwfpAIS69/AL1CEs2cLpEkwFQ9N18RsNEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAbd9uHXZaGT2cvhRs7reawctIXtX1s3kTqM9YV+/wCpBpuIV/6rgYT7aH9jRhjANdrEOdwa6ztVmKDwAAAAAAFL2UnENgLDPyB3kO0Wo1JMobmXXPEhoqkM/+x9+LaKzUFXsFgPMcX85EpiWC28+deO51lDoISjk7NQNo0iiZMIX7cor9zcBHh47oQ6b+tEryjNU6NIKNmQO1rV5pcZVSoU/aWdG4DckMQPqfcAOm8TRoNxPqtn2SaMF4IoXCW5gAQOAgABfAMAAAA10gV4WD5qAUbucdLUhOA6Y5ent8KLStl5YwIVVy5HUyAAAAAAAAAAajFNMW1Dc1RKZW1kRkVaaFpvNVlkU3BldFk4bnVqWXHYIR8AAAAAAKUAAAAAAAAABt324ddloZPZy+FGzut5rBy0he1fWzeROoz1hX7/AKkPBAEQABQBARESDwISAwQFBhUHCAkKCwwTAQ0AEQnoAwAAAAAAAFQAAAAAAAAADwMBAAABCQEZjx9MOkUiY9QTss0X68vBoOWIc2TmJhoSqBeS6hZaPgACBQo="
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Error(err)
	}

	// parse transaction:
	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(data))
	if err != nil {
		t.Error(err)
	}

	client := rpc.New("https://aged-morning-glade.solana-mainnet.quiknode.pro/b57bbb1a4c8bdd409e1ac53aaedead26da057f59/")
	recentBlockhashResult, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentConfirmed)
	if err != nil {
		t.Error(err)
	}
	tx.Message.RecentBlockhash = recentBlockhashResult.Value.Blockhash

	simulateTransactionResponse, err := client.SimulateTransactionWithOpts(context.TODO(), tx, &rpc.SimulateTransactionOpts{ReplaceRecentBlockhash: true})
	if err != nil {
		t.Error(err)
	}

	t.Log(simulateTransactionResponse.Value.Logs)

	t.Log(tx.Message.AccountKeys[1].String())
	t.Log(tx)
}

// export const LOOKUP_TABLE_CACHE: CacheLTA = {
// 	'2immgwYNHBbyVQKVGCEkgWpi53bLwWNRMB5G2nbgYV17': new AddressLookupTableAccount({
// 	  key: new PublicKey('2immgwYNHBbyVQKVGCEkgWpi53bLwWNRMB5G2nbgYV17'),
// 	  state: AddressLookupTableAccount.deserialize(
// 		Buffer.from(
// 		  'AQAAAP//////////d49+DAAAAAAAAQZMWvw7GUNJdaccNBVnb57OKakxL2BHLYvhRwVILRsgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMGRm/lIRcy/+ytunLDm+e8jOW7xfcSayxDmzpAAAAABt324ddloZPZy+FGzut5rBy0he1fWzeROoz1hX7/AKkG3fbh7nWP3hhCXbzkbM3athr8TYO5DSf+vfko2KGL/AVKU1D4XciC1hSlVnJ4iilt3x6rq9CmBniISTL07vagBqfVFxksXFEhjMlMPUrxf1ja7gibof1E49vZigAAAAAGp9UXGMd0yShWY5hpHV62i164o5tLbVxzVVshAAAAAIyXJY9OJInxuz0QKRSODYMLWhOZ2v8QhASOe9jb6fhZC3BlsePRfEU4nVJ/awTDzVi4bHMaoP21SbbRvAP4KUbIScv+6Yw2LHF/6K0ZjUPibbSWXCirYPGuuVl7zT789IUPLW4CpHr4JNCatp3ELXDLKMv6JJ+37le50lbBJ2LvDQdRqCgtphMF/imcN7mY5YRx2xE1A3MQ+L4QRaYK9u4GRfZP3LsAd00a+IkCpA22UNQMKdq5BFbJuwuOLqc8zxCTDlqxBG8J0HcxtfogQHDK06ukzfaXiNDKAob1MqBHS9lJxDYCwz8gd5DtFqNSTKG5l1zxIaKpDP/sffi2is1H9aKveyXSu5StXElYRl9SD5As0DHE4N0GLnf84/siiKXVyp4Ez121kLcUui/jLLFZEz/BwZK3Ilf9B9OcsEAeDMKAy2vjGSxQODgBz0QwGA+eP4ZjIjrIAQaXENv31QfLlOdXSRCkaybRniDHF4C8YcwhcvsqrOVuTP4B2Na+9wLdtrB31uz2rtlFI5kahdsnp/d1SrASDInYCtTYtdoke4kX+hoKWcEWM4Tle8pTUkUVv4BxS6fje/EzKBE4Qu9N9LMnrw/JNO0hqMVB4rk/2ou4AB1loQ7FZoPwut2o4KZB+0p9xnbrQKw038qjpHar+PyDwvxBRcu5hpHw3dguezeWv+IwvgW5icu8EGkhGa9AkFPPJT7VMSFb8xowveU=',
// 		  'base64',
// 		),
// 	  ),
// 	}),
//   }

func TestLookupTableCache(t *testing.T) {
	key := solana.MustPublicKeyFromBase58("2immgwYNHBbyVQKVGCEkgWpi53bLwWNRMB5G2nbgYV17")
	t.Log(key)

	data, err := base64.StdEncoding.DecodeString("AQAAAP//////////d49+DAAAAAAAAQZMWvw7GUNJdaccNBVnb57OKakxL2BHLYvhRwVILRsgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMGRm/lIRcy/+ytunLDm+e8jOW7xfcSayxDmzpAAAAABt324ddloZPZy+FGzut5rBy0he1fWzeROoz1hX7/AKkG3fbh7nWP3hhCXbzkbM3athr8TYO5DSf+vfko2KGL/AVKU1D4XciC1hSlVnJ4iilt3x6rq9CmBniISTL07vagBqfVFxksXFEhjMlMPUrxf1ja7gibof1E49vZigAAAAAGp9UXGMd0yShWY5hpHV62i164o5tLbVxzVVshAAAAAIyXJY9OJInxuz0QKRSODYMLWhOZ2v8QhASOe9jb6fhZC3BlsePRfEU4nVJ/awTDzVi4bHMaoP21SbbRvAP4KUbIScv+6Yw2LHF/6K0ZjUPibbSWXCirYPGuuVl7zT789IUPLW4CpHr4JNCatp3ELXDLKMv6JJ+37le50lbBJ2LvDQdRqCgtphMF/imcN7mY5YRx2xE1A3MQ+L4QRaYK9u4GRfZP3LsAd00a+IkCpA22UNQMKdq5BFbJuwuOLqc8zxCTDlqxBG8J0HcxtfogQHDK06ukzfaXiNDKAob1MqBHS9lJxDYCwz8gd5DtFqNSTKG5l1zxIaKpDP/sffi2is1H9aKveyXSu5StXElYRl9SD5As0DHE4N0GLnf84/siiKXVyp4Ez121kLcUui/jLLFZEz/BwZK3Ilf9B9OcsEAeDMKAy2vjGSxQODgBz0QwGA+eP4ZjIjrIAQaXENv31QfLlOdXSRCkaybRniDHF4C8YcwhcvsqrOVuTP4B2Na+9wLdtrB31uz2rtlFI5kahdsnp/d1SrASDInYCtTYtdoke4kX+hoKWcEWM4Tle8pTUkUVv4BxS6fje/EzKBE4Qu9N9LMnrw/JNO0hqMVB4rk/2ou4AB1loQ7FZoPwut2o4KZB+0p9xnbrQKw038qjpHar+PyDwvxBRcu5hpHw3dguezeWv+IwvgW5icu8EGkhGa9AkFPPJT7VMSFb8xowveU=")
	if err != nil {
		t.Error(err)
	}

	lookupTableState, err := addresslookuptable.DecodeAddressLookupTableState(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(lookupTableState)
}

func TestGeneratePubKey(t *testing.T) {
	// seed := solana.NewWallet().PublicKey().String()[0:32]
	fromPublicKey := solana.MustPublicKeyFromBase58("4d6MQwQC21eXMWToBiTL3UbknXwN3xzZ5Af8EyED4554")
	programId := solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	seed := "CPbW7ehsqd6agh8SZsmKipUPYXHtH4ZN"
	var buffer []byte
	buffer = append(buffer, fromPublicKey.Bytes()...)
	buffer = append(buffer, []byte(seed)...)
	buffer = append(buffer, programId.Bytes()...)
	publicKeyBytes := sha256.Sum256(buffer)
	pubKey := solana.PublicKeyFromBytes(publicKeyBytes[:])
	t.Log(pubKey)
}

// Program 11111111111111111111111111111111 invoke [1]
// Program 11111111111111111111111111111111 success
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]
// Program log: Instruction: InitializeAccount
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 3443 of 799850 compute units
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success
// Program 675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8 invoke [1]
// Program log: ray_log: A+gDAAAAAAAAUgAAAAAAAAABAAAAAAAAAOgDAAAAAAAAHodcPlgDAAC9rZJn7ycAAFMAAAAAAAAA
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]
// Program log: Instruction: Transfer
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 4736 of 777451 compute units
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]
// Program log: Instruction: Transfer
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 4645 of 769734 compute units
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success
// Program 675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8 consumed 32121 of 796407 compute units
// Program 675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8 success
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]
// Program log: Instruction: CloseAccount
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 2915 of 764286 compute units
// Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success

func BenchmarkGetProgramAccounts(b *testing.B) {
	client := rpc.New("https://aged-morning-glade.solana-mainnet.quiknode.pro/b57bbb1a4c8bdd409e1ac53aaedead26da057f59/")

	start := time.Now()
	// offset := uint64(0)
	// length := uint64(8)
	programAccountsResult, err := client.GetProgramAccountsWithOpts(
		context.TODO(),
		solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"),
		&rpc.GetProgramAccountsOpts{
			// DataSlice: &rpc.DataSlice{
			// 	Offset: &offset,
			// 	Length: &length,
			// },
			Filters: []rpc.RPCFilter{
				{
					DataSize: 1544,
				},
			},
		},
	)
	if err != nil {
		b.Error(err)
	}
	b.Log(time.Since(start), len(programAccountsResult))
}
