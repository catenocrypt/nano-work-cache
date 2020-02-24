// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"fmt"
	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func negDiff(diff uint64) uint64 {
	return 0xffffffffffffffff - diff
}

func incDiff(diff uint64) uint64 {
	return negDiff(negDiff(diff) / 4 * 3)
}

func rpcCallPrint(url string, reqJson string) {
	fmt.Println("Calling url", url, "with data", reqJson)
	resp, err := rpcclient.RpcCall(url, reqJson)
	if (err != nil) {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response: ")
	fmt.Println(resp)
}

func main() {
	var url = "http://localhost:7376"
	//var url = "https://nanovault.io/api/node-api"
	
	rpcCallPrint(url, `{"action":"block_count"}`)

	//hash2 := "718CC2121C3E641059BC1C2CFC45666C99E8AE922F7A807B7D07B62C995D79E2"
	hash1 := "DDDA8C4CB5825FF4F5D00C5F923BC6F632414F67D17039228325392671C50FA2"
	account1 := "nano_3rpb7ddcd6kux978gkwxh1i1s6cyn7pw3mzdb9aq7jbtsdfzceqdt3jureju"
	
	rpcCallPrint(url, `{"action":"block_account","hash":"` +  hash1 + `"}`)
	rpcCallPrint(url, `{"action":"accounts_frontiers","accounts":["` +  account1 + `"]}`)
	
	//rpcCallPrint(url, `{"action": "work_generate","hash": "718CC2121C3E641059BC1C2CFC45666C99E8AE922F7A807B7D07B62C995D79E2","difficulty": "ffffffd21c3933f3"}`)
	//var diff1 uint64 = 0xffffffc800000000
	rpcCallPrint(url, fmt.Sprintf("{\"action\": \"work_generate\",\"hash\": \"%v\"}", hash1))
	
	rpcCallPrint(url, fmt.Sprintf("{\"action\": \"work_generate\",\"hash\": \"%v\"}", hash1))

	var diff uint64 = 0xffffffc000000000
	for negDiff(diff) > 0x800000000 {
		fmt.Printf("diff %x %x\n", diff, negDiff(diff))
		rpcCallPrint(url, fmt.Sprintf("{\"action\": \"work_generate\",\"hash\": \"%v\", \"difficulty\": \"%x\"}", hash1, diff))
		diff = incDiff(diff)
	}
}
