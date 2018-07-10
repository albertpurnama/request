package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/parnurzeal/gorequest"
)

const CREATE_NEW_ADDR_URL = "https://api.blockcypher.com/v1/beth/test/addrs?token=e7066438be7f48d1b08a325450d31694"
const CREATE_NEW_TX_URL = "https://api.blockcypher.com/v1/beth/test/txs/new?token=e7066438be7f48d1b08a325450d31694"
const FAUCET_URL = "https://api.blockcypher.com/v1/beth/test/faucet?token=e7066438be7f48d1b08a325450d31694"
const SEND_TX_URL = "https://api.blockcypher.com/v1/beth/test/txs/send?token=e7066438be7f48d1b08a325450d31694"

const (
	BETH_ADDR        = "3c689107595e512e5d0d01eaeb32ac345e928980"
	BETH_PRIVATE     = "721b7ede0647cf36761afb76f7ae2ac24009831ff535680c056ae284fb3ca0eb"
	BETH_PUBLIC      = "04c1e062adb0af8a0636ae0ea0a6db3932cb4e9ec28d6bb2fa706ecf6644c5e9b5b18f779a7d6e91ac605e21ea53bc2cee7e9e0a770f2a8127ebea88f3ee6a4042"
	BETH_ADDR_TWO    = "aa2397a02db808e10b043a848952ce9d4d3161f8"
	BETH_PRIVATE_TWO = "183b827d2707eaee6483fcd653ced915ada8fdd73f8d536c84b20476a62c6ed3"
	BETH_PUBLIC_TWO  = "04ac417ab0a3ad18006329701b0c8e82cc6eda4787d62dd23166650f0025bb512f7700250c94f05456aa7b1c89e524532f063bd0fd15f94d93b710c2dc3f8fad8e"
)

type EtherTX struct {
	BlockHeight int       `json:"block_height" binding:"required"`
	BlockIndex  uint      `json:"block_index" `
	Hash        string    `json:"hash" binding:"required"`
	Addresses   []string  `json:"addresses" binding:"required"`
	Total       uint      `json:"total" binding:"required"`
	Fees        uint      `json:"fees" binding:"required"`
	Size        uint      `json:"size" binding:"required"`
	GasUsed     uint      `json:"gas_used"`
	GasLimit    uint      `json:"gas_limit"`
	GasPrice    uint      `json:"gas_price" binding:"required"`
	Received    time.Time `json:"received" binding:"required"`
	Version     int       `json:"ver" binding:"required"`
	DoubleSpend bool      `json:"double_spend" binding:"required"`
	VinSz       uint      `json:"vin_sz" binding:"required"`
	VoutSz      uint      `json:"vout_sz" binding:"required"`
	Inputs      []struct {
		Sequence  uint     `json:"sequence"`
		Addresses []string `json:"addresses"`
	} `json:"inputs" binding:"required"`
	Outputs []struct {
		Value     uint     `json:"value"`
		Script    string   `json:"script"`
		Addresses []string `json:"addresses"`
	} `json:"outputs" binding:"required"`
	Confirmed time.Time `json:confirmed`
	BlockHash string    `json:block_hash`
	ExecError string    `json:execution_error`
}

type EtherTxSkeleton struct {
	Tx         EtherTX  `json:"tx" binding:"required"`
	Tosign     []string `json:"tosign" binding:"required"`
	Signatures []string `json:"signatures"`
	Errors     []struct {
		Error string `json:"error"`
	} `json:"errors"`
}

func main() {
	// var dat EtherTxSkeleton
	request := gorequest.New()

	// Making transaction
	data, err := MakeEtherTransaction(BETH_ADDR, BETH_ADDR_TWO, 15, request)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(data)
	// *************

	// resp, body, _ := request.Post(CREATE_NEW_TX_URL).
	// 	Send(CreatePartiallyFilledTX(BETH_ADDR, BETH_ADDR_TWO, 1)).
	// 	End()

	// if resp.StatusCode != 200 && resp.StatusCode != 201 {
	// 	log.Fatal(body)
	// 	return
	// }
	// fmt.Println(body)
	// json.Unmarshal([]byte(body), &dat)

	// str, err := Sign(BETH_PRIVATE, dat.Tosign[0])
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// fmt.Println(str)

	// dat.Signatures = append(dat.Signatures, str)

	// abc, _ := json.Marshal(dat)
	// fmt.Println(string(abc))

	// resp, body, _ = request.Post(SEND_TX_URL).
	// 	Send(string(abc)).
	// 	End()
	// if resp.StatusCode != 200 && resp.StatusCode != 201 {
	// 	log.Fatal(body)
	// 	return
	// }
	// fmt.Println(body)

	// // faucet
	// resp, body, _ := request.Post(FAUCET_URL).
	// 	Send(CreateFaucetString(BETH_ADDR, 1e18)).
	// 	End()
	// if resp.StatusCode != 200 {
	// 	fmt.Println(resp)
	// 	log.Fatal(CreateFaucetString(BETH_ADDR, 1e18))
	// 	return
	// }
	// fmt.Println(body)

	// // get balance
	// resp, body, _ := request.Get(GetETHBalanceURL(BETH_ADDR_TWO)).
	// 	End()
	// if resp.StatusCode != 200 {
	// 	fmt.Println(body)
	// 	log.Fatal("Status code not 200")
	// 	return
	// }
	// fmt.Println(body)
}

// MakeEtherTransaction makes Ethereum transaction
func MakeEtherTransaction(sender string, receiver string, amount float64, r *gorequest.SuperAgent) (EtherTxSkeleton, error) {
	// Create a new empty template for handling API calls payload
	var dat EtherTxSkeleton

	// Create initial TX structure using BlockCypher's TX endpoints
	resp, body, _ := r.Post(CREATE_NEW_TX_URL).
		Send(CreatePartiallyFilledTX(sender, receiver, amount)).
		End()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Fatal(body)
		return EtherTxSkeleton{}, errors.New("Error creating initial unsigned transaction")
	}

	json.Unmarshal([]byte(body), &dat)

	// Sign the initial transaction's tosign
	str, err := Sign(BETH_PRIVATE, dat.Tosign[0])
	if err != nil {
		log.Fatal(err)
		return EtherTxSkeleton{}, errors.New("Error signing transaction")
	}

	dat.Signatures = append(dat.Signatures, str)

	signedTxSkel, _ := json.Marshal(dat)

	// Send the initial transaction with signature
	resp, body, _ = r.Post(SEND_TX_URL).
		Send(string(signedTxSkel)).
		End()
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Fatal(body)
		return EtherTxSkeleton{}, errors.New("Error sending signed transaction")
	}

	json.Unmarshal([]byte(body), &dat)

	return dat, nil
}

// Sign functions signs the transaction with the main account's private
// key
func Sign(private string, data string) (str string, err error) {
	dat, err := hex.DecodeString(data)
	if err != nil {
		return str, err
	}

	priv, err := hex.DecodeString(private)
	if err != nil {
		return str, err
	}

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), priv)
	sig, err := privKey.Sign(dat)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sig.Serialize()), nil
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
