package raydium

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
)

func BenchmarkFormatClmmKeys(b *testing.B) {
	client := rpc.New("https://aged-morning-glade.solana-mainnet.quiknode.pro/b57bbb1a4c8bdd409e1ac53aaedead26da057f59/")

	res, err := FormatClmmKeys(client)
	b.Log(res, err)
}

func TestGenerateClmmTransaction(t *testing.T) {
	client := rpc.New("https://aged-morning-glade.solana-mainnet.quiknode.pro/b57bbb1a4c8bdd409e1ac53aaedead26da057f59/")

	poolsInfo, err := FormatClmmKeys(client)
	if err != nil {
		t.Error(err)
		return
	}

	poolsInfoDict := make(map[string][]*ClmmPoolInfo)
	for _, poolInfo := range poolsInfo {
		k1 := poolInfo.MintA.Mint.String() + poolInfo.MintB.Mint.String()
		poolsInfoDict[k1] = append(poolsInfoDict[k1], poolInfo)
		k2 := poolInfo.MintB.Mint.String() + poolInfo.MintA.Mint.String()
		poolsInfoDict[k2] = append(poolsInfoDict[k2], poolInfo)
	}

	inputToken := "So11111111111111111111111111111111111111112"
	outputToken := "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"
	tokenProgramId := "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"

	poolsList := poolsInfoDict[inputToken+outputToken]
	poolsList = append(poolsList, poolsInfoDict[outputToken+inputToken]...)
	poolsList = append(poolsList, poolsInfoDict[inputToken+tokenProgramId]...)
	poolsList = append(poolsList, poolsInfoDict[outputToken+tokenProgramId]...)
	poolsList = append(poolsList, poolsInfoDict[tokenProgramId+inputToken]...)
	poolsList = append(poolsList, poolsInfoDict[tokenProgramId+outputToken]...)
	sort.Slice(poolsList, func(i, j int) bool {
		return poolsList[i].Id.String() < poolsList[j].Id.String()
	})
	httpClient := http.DefaultClient
	for {
		request := &RouteBuildRequest{
			InputToken:  inputToken,
			OutputToken: outputToken,
			Amount:      100,
			Slippage:    1,
			PublicKey:   "4d6MQwQC21eXMWToBiTL3UbknXwN3xzZ5Af8EyED4554",
			ClmmList:    poolsList,
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
