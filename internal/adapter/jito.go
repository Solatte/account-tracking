package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RpcResult struct {
	JsonRpc string   `json:"jsonrpc"`
	Result  []string `json:"result"`
	Id      int      `json:"id"`
}

func GetJitoTipAccounts() []string {
	url := "https://mainnet.block-engine.jito.wtf/api/v1/bundles"

	// Define the request body
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTipAccounts",
		"params":  []interface{}{},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}
	}

	var res RpcResult
	err = json.Unmarshal(body, &res)

	if err != nil {
		fmt.Println("unmarshal result error:", err)
		return []string{}
	}

	fmt.Println("Response:", res)

	return res.Result

}
