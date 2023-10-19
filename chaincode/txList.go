package main

import (
	txdefs "github.com/hyperledger-labs/cc-tools-demo/chaincode/txdefs"

	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var txList = []tx.Transaction{
	tx.CreateAsset,
	tx.UpdateAsset,
	tx.DeleteAsset,
	txdefs.CreateNewManager,
	txdefs.CreateNewHolder,
	txdefs.CreateNewCreditCard,
	txdefs.CreateNewTransferency,
	txdefs.UpdateCreditCardLimit,
	txdefs.UpdateCreditCardName,
	txdefs.CreateNewPurchase,
	txdefs.CreateNewCreditCardPurchase,
	txdefs.PayCreditCardInvoice,
}
