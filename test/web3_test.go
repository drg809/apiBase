package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nikola43/fibergormapitemplate/web3Manager"
	"github.com/stretchr/testify/assert"
)

func TestIsValidAddress(t *testing.T) {

	add1 := "0x0asdasdas8f21bE3dE0BA2ba6918E714dA6B45836"
	valid := web3Manager.IsValidAddress(add1)

	if valid {
		fmt.Println("La dirreción es valida, vamos a enviar")
	} else {
		fmt.Println("Error comprueba la direccion")
	}

	fmt.Println(valid)
	assert.True(t, valid)
}

func TestAddressIsContract(t *testing.T) {

	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	address := "0xHaD91ee08f21bE3dE0BA2ba6918E714dA6B45836"
	isContract := web3Manager.AddressIsContract(client, address)

	assert.True(t, isContract)
}
