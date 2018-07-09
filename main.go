package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/parnurzeal/gorequest"
)

const CREATE_NEW_ADDR_URL = "https://api.blockcypher.com/v1/beth/test/addrs?token=e7066438be7f48d1b08a325450d31694"
const CREATE_NEW_TX_URL = "https://api.blockcypher.com/v1/beth/test/txs/new?token=e7066438be7f48d1b08a325450d31694"
const FAUCET_URL = "https://api.blockcypher.com/v1/beth/test/faucet?token=e7066438be7f48d1b08a325450d31694"

const (
	BETH_ADDR        = "3c689107595e512e5d0d01eaeb32ac345e928980"
	BETH_PRIVATE     = "721b7ede0647cf36761afb76f7ae2ac24009831ff535680c056ae284fb3ca0eb"
	BETH_PUBLIC      = "04c1e062adb0af8a0636ae0ea0a6db3932cb4e9ec28d6bb2fa706ecf6644c5e9b5b18f779a7d6e91ac605e21ea53bc2cee7e9e0a770f2a8127ebea88f3ee6a4042"
	BETH_ADDR_TWO    = "aa2397a02db808e10b043a848952ce9d4d3161f8"
	BETH_PRIVATE_TWO = "183b827d2707eaee6483fcd653ced915ada8fdd73f8d536c84b20476a62c6ed3"
	BETH_PUBLIC_TWO  = "04ac417ab0a3ad18006329701b0c8e82cc6eda4787d62dd23166650f0025bb512f7700250c94f05456aa7b1c89e524532f063bd0fd15f94d93b710c2dc3f8fad8e"
)

func main() {
	var dat map[string]interface{}
	request := gorequest.New()
	resp, body, _ := request.Post(CREATE_NEW_TX_URL).
		Send(CreatePartiallyFilledTX(BETH_ADDR, BETH_ADDR_TWO, 1)).
		End()

	if resp.StatusCode != 200 || resp.StatusCode != 201 {
		log.Fatal(body)
		return
	}

	fmt.Println(body)
	json.Unmarshal([]byte(body), dat)
	fmt.Println(dat)

	// // faucet
	// resp, body, _ := request.Post(FAUCET_URL).
	// 	Send(str).
	// 	End()
	// if resp.StatusCode != 200 {
	// 	log.Fatal(CreateFaucetString(BETH_ADDR, 1e18))
	// 	return
	// }
	// fmt.Println(body)

	// // get balance
	// resp, body, _ := request.Get(GetETHBalanceURL(BETH_ADDR)).
	// 	End()
	// if resp.StatusCode != 200 {
	// 	fmt.Println(body)
	// 	log.Fatal("Status code not 200")
	// 	return
	// }
	// fmt.Println(body)
}

// GetETHBalanceURL create the URL for getting Ethereum balance
func GetETHBalanceURL(addr string) (str string) {
	str = fmt.Sprintf("https://api.blockcypher.com/v1/beth/test/addrs/%s/balance", addr)
	return str
}

// CreateFaucetString creates POST payload for Faucet
func CreateFaucetString(addr string, amount uint) (str string) {
	str = fmt.Sprintf("{\"address\": \"%s\", \"amount\": %d}", addr, amount)
	fmt.Println(str)
	return str
}

// CreatePartiallyFilledTX creates a partially filled tx ( no shit bruh )
func CreatePartiallyFilledTX(sender string, recipient string, amount float64) (str string) {
	amountInWei := uint(amount * 1e18)
	str = fmt.Sprintf("{\"inputs\":[{\"addresses\": [\"%s\"]}],\"outputs\":[{\"addresses\": [\"%s\"], \"value\": %d}]}", sender, recipient, amountInWei)
	fmt.Println(str)
	return str
}
