package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type BlockNumberResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"` // Block number as a hex string (e.g., "0x5B9AC0")
}

func getBlockNumber() (int64, error) {
	// JSON-RPC payload
	currentTime := time.Now()
	id := rand.Intn(100)
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      id,
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return 0, err
	}

	rpc := os.Getenv("RPC")
	if rpc == "" {
		rpc = "http://127.0.0.1:8545"
	}

	// Send the request
	resp, err := http.Post(rpc, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return 0, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return 0, err
	}

	// Unmarshal the response into BlockNumberResponse struct
	var blockNumberResp BlockNumberResponse
	err = json.Unmarshal(body, &blockNumberResp)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return 0, err
	}

	// Print the block number in hex format
	fmt.Println("currentTime:", currentTime, "ID ", id, "Block Number (Hex):", blockNumberResp.Result)

	// Optionally, convert the hex block number to a decimal
	var blockNumber int64
	_, err = fmt.Sscanf(blockNumberResp.Result, "0x%x", &blockNumber)
	if err != nil {
		fmt.Println("Error converting hex to decimal:", err)
		return 0, err
	}

	// Print the block number in decimal format
	fmt.Println("currentTime:", currentTime, "ID ", id, "Block Number (Decimal):", blockNumber)
	return blockNumber, nil

}

func checkSync(w http.ResponseWriter, r *http.Request) {
	blockNumberFirst, err := getBlockNumber()
	if err != nil {
		w.WriteHeader(500)
	}

	time.Sleep(30 * time.Second)

	blockNumberSecond, err := getBlockNumber()
	if err != nil {
		w.WriteHeader(500)
	}

	blockNumber := blockNumberFirst - blockNumberSecond
	if blockNumber != 0 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
}

func main() {

	http.HandleFunc("/", checkSync)
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	fmt.Printf("Starting server on port %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
