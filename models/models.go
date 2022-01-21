// Package models implements the data layer interface
package models

import (
	"errors"
    "time"
	"fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "encoding/json"
	"github.com/gastonstec/utils"
)

const( 
    HTTP_TIMEOUT = 3
    HTTP_STATUSOK = 200
)

// Block information struct
type BlockInfo struct {
	Status string `json:"status"`
	Data   struct {
	  Network           string   `json:"network"`
	  Blockhash         string   `json:"blockhash"`
	  BlockNo           int      `json:"block_no"`
	  MiningDifficulty  string   `json:"mining_difficulty"`
	  Time              int      `json:"time"`
	  Confirmations     int      `json:"confirmations"`
	  IsOrphan          bool     `json:"is_orphan"`
	  Txs               []string `json:"txs"`
	  Merkleroot        string   `json:"merkleroot"`
	  PreviousBlockhash string   `json:"previous_blockhash"`
	  NextBlockhash     string   `json:"next_blockhash"`
	  Size              int      `json:"size"`
	} `json:"data"`
}

// Transaction informatio struct
type TxInfo struct {
	Status string `json:"status"`
	Data   struct {
	  Network       string `json:"network"`
	  Txid          string `json:"txid"`
	  Blockhash     string `json:"blockhash"`
	  BlockNo       int    `json:"block_no"`
	  Confirmations int    `json:"confirmations"`
	  Time          int64  `json:"time"`
	  Size          int    `json:"size"`
	  Vsize         int    `json:"vsize"`
	  Version       int    `json:"version"`
	  Locktime      int    `json:"locktime"`
	  SentValue     string `json:"sent_value"`
	  Fee           string `json:"fee"`
	  Inputs        []struct {
		InputNo      int    `json:"input_no"`
		Address      string `json:"address"`
		Value        string `json:"value"`
		ReceivedFrom struct {
		  Txid     string `json:"txid"`
		  OutputNo int    `json:"output_no"`
		} `json:"received_from"`
		ScriptAsm string      `json:"script_asm"`
		ScriptHex string      `json:"script_hex"`
		Witness   interface{} `json:"witness"`
	  } `json:"inputs"`
	  Outputs []struct {
		OutputNo int         `json:"output_no"`
		Address  string      `json:"address"`
		Value    string      `json:"value"`
		Type     string      `json:"type"`
		ReqSigs  interface{} `json:"req_sigs"`
		Spent    struct {
		  Txid    string `json:"txid"`
		  InputNo int    `json:"input_no"`
		} `json:"spent"`
		ScriptAsm string `json:"script_asm"`
		ScriptHex string `json:"script_hex"`
	  } `json:"outputs"`
	  TxHex string `json:"tx_hex"`
	} `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}


func GetTxInfo(networkCode string, txId string) (*TxInfo, error) {
	var err error

    // Validate parameters
	if networkCode == "" || txId == "" {
		return nil, errors.New(utils.GetFunctionName() + ": invalid parameters")
	}

    // Call get transaction api
	uri:= "https://sochain.com/api/v2/tx/" + networkCode + "/" + txId
    client := &http.Client{Timeout: HTTP_TIMEOUT * time.Second,}
    response, err := client.Get(uri)
    if err != nil {
		return nil, errors.New(utils.GetFunctionName() + ": The http request failed with " + err.Error())
    }
	if response.StatusCode != HTTP_STATUSOK {
		return nil, errors.New(utils.GetFunctionName() + ": The http request failed with " + fmt.Sprint(response.Status))
	}

    // Send error if response is not OK
    if response.StatusCode != HTTP_STATUSOK {
        return nil, errors.New("invalid block hash")
    }

	 // Read response body
     body, err := ioutil.ReadAll(response.Body)
     if err != nil {
         return nil, errors.New(utils.GetFunctionName() + ": Reading body failed with " + err.Error())
     }

    jsonDataReader := strings.NewReader(string(body))
    decoder := json.NewDecoder(jsonDataReader)
    var tx TxInfo
    err = decoder.Decode(&tx)
	if err != nil {
		return nil, errors.New(utils.GetFunctionName() + ": Body decoding failed with " + err.Error())
    }


    return &tx, nil

}


func GetBlockInfo(networkCode string, blockHash string) (*BlockInfo, error) {
	var err error

    // Call get_block api
	uri:= "https://sochain.com/api/v2/get_block/" + networkCode + "/" + blockHash
    client := &http.Client{Timeout: HTTP_TIMEOUT * time.Second,}
    response, err := client.Get(uri)
    if err != nil {
		return nil, errors.New(utils.GetFunctionName() + ": The HTTP request failed with " + err.Error())
    }
    
    // Send error if response is not OK
    if response.StatusCode != HTTP_STATUSOK {
        return nil, errors.New("invalid block hash")
    }

    // Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New(utils.GetFunctionName() + ": Reading body failed with " + err.Error())
    }

    jsonDataReader := strings.NewReader(string(body))
    decoder := json.NewDecoder(jsonDataReader)
    var block BlockInfo
    err = decoder.Decode(&block)
    if err != nil {
		return nil, errors.New(utils.GetFunctionName() + ": Body decoding failed with " + err.Error())
    }

	
	return &block, nil

}