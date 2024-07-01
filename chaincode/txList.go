package main

import (
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs"
)

var txList = []tx.Transaction{
	tx.CreateAsset,
	tx.UpdateAsset,
	tx.DeleteAsset,
	txdefs.UploadDocument,
	txdefs.PutSignature,
	txdefs.CancelDocument,
	txdefs.CreateSigner,
	txdefs.GetDoc,
	txdefs.GetSignerKey,
	txdefs.GetSigner,
	txdefs.GetExpiredDoc,
}
