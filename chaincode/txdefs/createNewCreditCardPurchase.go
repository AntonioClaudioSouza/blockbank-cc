package txdefs

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/blockbank-cc/chaincode/utils"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateNewCreditCardPurchase = tx.Transaction{
	Tag:         "createNewCreditCardPurchase",
	Label:       "Create new credit card purchase",
	Description: "Create a new credit card purchase",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:         "creditCard",
			Label:       "Credit Card",
			Description: "Credit Card",
			DataType:    "->creditCard",
			Required:    true,
		},
		{
			Tag:         "description",
			Label:       "Purchase description",
			Description: "Purchase description",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "value",
			Label:       "Purchase value",
			Description: "Purchase value",
			DataType:    "number",
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		purchaseDescription, _ := req["description"].(string)
		purchaseValue, _ := req["value"].(float64)
		creditCardKey, _ := req["creditCard"].(assets.Key)

		// ** Get creditCard and check limit
		creditCardMap, err := creditCardKey.GetMap(stub)
		if err != nil {
			return nil, errors.NewCCError("failed to get owner from the ledger: "+err.Error(), http.StatusBadRequest)
		}

		availableLimit := creditCardMap["limit"].(float64) - creditCardMap["limitUsed"].(float64)
		if availableLimit < purchaseValue {
			return nil, errors.NewCCError("credit card doesn't have enough limit to make this purchase", http.StatusBadRequest)
		}

		// ** REMOVE BALANCE FROM BUYER AND UPDATE
		creditCardMap, err = creditCardKey.Update(stub, map[string]interface{}{
			"limitUsed": creditCardMap["limitUsed"].(float64) + purchaseValue,
		})
		if err != nil {
			return nil, errors.NewCCError("failed to update creditCard asset: "+err.Error(), http.StatusBadRequest)
		}

		timeStamp, _ := stub.Stub.GetTxTimestamp()
		ccPurchaseMap := make(map[string]interface{})
		ccPurchaseMap["@assetType"] = "creditCardPurchase"
		ccPurchaseMap["txId"] = stub.Stub.GetTxID()
		ccPurchaseMap["description"] = purchaseDescription
		ccPurchaseMap["creditCard"] = creditCardMap
		ccPurchaseMap["value"] = purchaseValue
		ccPurchaseMap["date"] = utils.ReturnDate(timeStamp)

		ccPurchaseAsset, err := assets.NewAsset(ccPurchaseMap)
		if err != nil {
			return nil, errors.NewCCError("failed to create purchase asset: "+err.Error(), http.StatusBadRequest)
		}

		_, err = ccPurchaseAsset.PutNew(stub)
		if err != nil {
			return nil, errors.NewCCError("error saving asset on blockchain: "+err.Error(), http.StatusBadRequest)
		}

		ccPurchaseJSON, nerr := json.Marshal(ccPurchaseAsset)
		if nerr != nil {
			return nil, errors.NewCCError("failed to encode asset to JSON format: "+nerr.Error(), http.StatusBadRequest)
		}

		return ccPurchaseJSON, nil
	},
}
