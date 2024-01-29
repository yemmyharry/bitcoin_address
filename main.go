package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"log"
	"os"
)

func main() {
	//createP2PKHAddress()

	rawTx, err := createTx(os.Getenv("PRIVATE_KEY"), os.Getenv("PUBLIC_KEY"), 6000)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(" Raw signed transaction is ", *rawTx)

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

func NewTx() (*wire.MsgTx, error) {
	return wire.NewMsgTx(wire.TxVersion), nil
}

func GetUTXO(address string) (string, int64, string, error) {
	var prevTxID string = os.Getenv("PREV_TX_ID")
	var balance int64 = 62000
	var pubKeyScript string = os.Getenv("PUB_KEY_SCRIPT")

	return prevTxID, balance, pubKeyScript, nil

}

func createTx(privKey string, destination string, amount int64) (*string, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return nil, err
	}

	addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		return nil, err
	}

	txId, balance, pkScript, err := GetUTXO(addrPubKey.EncodeAddress())
	if err != nil {
		return nil, err
	}
	if balance < amount {
		return nil, errors.New("insufficient bal")
	}

	destAddr, err := btcutil.DecodeAddress(destination, &chaincfg.TestNet3Params)
	if err != nil {
		return nil, err
	}

	destAddrByte, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		return nil, err
	}

	redeemTx, err := NewTx()
	if err != nil {
		return nil, err
	}

	utxoHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		return nil, err
	}

	outPoint := wire.NewOutPoint(utxoHash, 1)
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	redeemTxOut := wire.NewTxOut(amount, destAddrByte)
	redeemTx.AddTxOut(redeemTxOut)

	finalRawTx, err := signTx(privKey, pkScript, redeemTx)

	return finalRawTx, nil
}

func signTx(privKey, pkScript string, redeemTx *wire.MsgTx) (*string, error) {

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return nil, err
	}
	sourcePKScript, err := hex.DecodeString(pkScript)
	if err != nil {
		return nil, err
	}

	signature, err := txscript.SignatureScript(redeemTx, 0, sourcePKScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return nil, err
	}

	redeemTx.TxIn[0].SignatureScript = signature

	var signedTx bytes.Buffer
	err = redeemTx.Serialize(&signedTx)
	if err != nil {
		return nil, err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return &hexSignedTx, nil

}
