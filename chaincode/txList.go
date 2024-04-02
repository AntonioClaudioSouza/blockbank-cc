package main

import (
	txdefs "github.com/hyperledger-labs/blockbank-cc/chaincode/txdefs"

	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var txList = []tx.Transaction{
	txdefs.CreateNewManager,
	txdefs.CreateNewHolder,
	txdefs.CreateNewCreditCard,
	txdefs.CreateNewTransferency,
	txdefs.UpdateCreditCardLimit,
	txdefs.UpdateCreditCardName,
	txdefs.CreateNewPurchase,
	txdefs.CreateNewCreditCardPurchase,
	txdefs.PayCreditCardInvoice,
	txdefs.MakeDeposit,
	txdefs.MakeWithdrawal,
	txdefs.HealthCheck,
	txdefs.ActivateCreditCard,
	txdefs.ListManagers,
	txdefs.ListHolders,
	txdefs.GetHolderByKey,
	txdefs.GetDepositsByHolderKey,
	txdefs.GetWithdrawalsByHolderKey,
	txdefs.GetTransfersByHolderKey,
	txdefs.GetCreditCardByHolderKey,
	txdefs.GetPurchasesByHolderKey,
	txdefs.GetCreditCardPurchasesByCreditCardKey,
	txdefs.GetPaymentsByCreditCardKey,
}
