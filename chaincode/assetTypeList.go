package main

import (
	"github.com/hyperledger-labs/cc-tools-demo/chaincode/assettypes"
	"github.com/hyperledger-labs/cc-tools/assets"
)

var assetTypeList = []assets.AssetType{
	assettypes.Manager,
	assettypes.Holder,
	assettypes.CreditCard,
	assettypes.Transferency,
	assettypes.Purchase,
	assettypes.CreditCardPurchase,
}
