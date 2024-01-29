package main

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"log"
)

func main() {
	createP2PKHAddress()

}

func createP2PKHAddress() {
	network := &chaincfg.TestNet3Params
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		log.Fatal(err)
	}

	pubKey := privKey.PubKey()
	addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), network)
	if err != nil {
		log.Fatal(err)
	}

	wif, err := btcutil.NewWIF(privKey, network, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Public Address is %s\n", addr)
	fmt.Printf(" Private Key is %s\n", wif.String())
}
