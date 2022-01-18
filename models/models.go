package models

import (
    //"bytes"
	"errors"
    //"encoding/json"
    //"fmt"
    "time"
    "io/ioutil"
    "net/http"
	"github.com/gastonstec/utils"
)

const HTTP_TIMEOUT = 15


func GetBlockInfo(networkCode string, blockHash string) (string, error) {
	var err error

	uri:= "https://sochain.com/api/v2/get_block/" + networkCode + "/" + blockHash

    client := &http.Client{Timeout: HTTP_TIMEOUT * time.Second,}
    response, err := client.Get(uri)
    if err != nil {
		return "", errors.New(utils.GetFunctionName() + ": The HTTP request failed with " + err.Error())
    } 
	
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New(utils.GetFunctionName() + ": Reading body failed with " + err.Error())
    }

    println(string(data))


    // jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
    // jsonValue, _ := json.Marshal(jsonData)
    // response, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
    // if err != nil {
    //     fmt.Printf("The HTTP request failed with error %s\n", err)
    // } else {
    //     data, _ := ioutil.ReadAll(response.Body)
    //     fmt.Println(string(data))
    // }
    // fmt.Println("Terminating the application...")
	
	return string(data), nil

}