package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"log"
)

func main() {

	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		log.Fatal(err)
	}

	pubKey := privKey.PubKey()
	addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.TestNet3Params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Public Address is %s\n", addr)
	fmt.Printf(" Private Key is %s\n", hex.EncodeToString(privKey.Serialize()))

}
