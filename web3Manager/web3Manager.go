package web3Manager

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"regexp"
)

type Web3Manager struct {
}

func AddressIsContract(client *ethclient.Client, address string) bool {
	bytecode, err := client.CodeAt(context.Background(), common.HexToAddress(address), nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}
	return len(bytecode) > 0
}

func IsValidAddress(address string) bool {
	return regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(address)
}
