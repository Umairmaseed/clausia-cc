package main

import (
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/contract"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/document"
)

var txList = []tx.Transaction{

	tx.CreateAsset,
	tx.UpdateAsset,
	tx.DeleteAsset,

	document.CancelDocument,
	document.UploadDocument,
	document.PutSignature,
	document.CreateSigner,
	document.GetDoc,
	document.GetUserKey,
	document.GetSigner,
	document.GetExpiredDoc,
	document.UpdateDocument,
	document.UpdateSigner,
	document.ExpectedUserDoc,
	document.GetDocHistory,

	contract.CreateAutoExecutableContract,
	contract.AddClause,
	contract.RemoveClause,
	contract.AddClauses,
	contract.AddParticipants,
	contract.AddReferenceDateCDI,
	contract.AddEvalutedDateCDI,
	contract.ContractsWithExecutableClauses,
	contract.ExecuteAutoExecutableContract,
	contract.AddInputToCheckFineClause,
	contract.AddReviewToContract,
	contract.AddInputsToMakePaymentClause,
}
