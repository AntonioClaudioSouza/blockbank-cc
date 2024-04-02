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

var CreateNewPurchase = tx.Transaction{
	Tag:         "createNewPurchase",
	Label:       "Create new purchase",
	Description: "Create a new purchase",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:         "buyer",
			Label:       "Buyer",
			Description: "Buyer",
			DataType:    "->holder",
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
		purchaseValue := req["value"].(float64)
		buyerKey, _ := req["buyer"].(assets.Key)

		buyerMap, err := buyerKey.GetMap(stub)
		if err != nil {
			return nil, errors.NewCCError("failed to get owner from the ledger: "+err.Error(), http.StatusBadRequest)
		}

		if buyerMap["cash"].(float64) < purchaseValue {
			return nil, errors.NewCCError("buyer doesn't have enough balance to make this purchase", http.StatusBadRequest)
		}

		// ** REMOVE BALANCE FROM BUYER
		buyerMap, err = buyerKey.Update(stub, map[string]interface{}{
			"cash": buyerMap["cash"].(float64) - purchaseValue,
		})
		if err != nil {
			return nil, errors.NewCCError("failed to update buyer asset: "+err.Error(), http.StatusBadRequest)
		}

		timeStamp, _ := stub.Stub.GetTxTimestamp()
		purchaseMap := make(map[string]interface{})
		purchaseMap["@assetType"] = "purchase"
		purchaseMap["txId"] = stub.Stub.GetTxID()
		purchaseMap["description"] = purchaseDescription
		purchaseMap["buyer"] = buyerMap
		purchaseMap["value"] = purchaseValue
		purchaseMap["date"] = utils.ReturnDate(timeStamp)

		purchaseAsset, err := assets.NewAsset(purchaseMap)
		if err != nil {
			return nil, errors.NewCCError("failed to create creditCard asset: "+err.Error(), http.StatusBadRequest)
		}

		purchaseMap, err = purchaseAsset.PutNew(stub)
		if err != nil {
			return nil, errors.NewCCError("error saving asset on blockchain: "+err.Error(), http.StatusBadRequest)
		}

		purchaseJSON, nerr := json.Marshal(purchaseMap)
		if nerr != nil {
			return nil, errors.NewCCError("failed to encode asset to JSON format: "+nerr.Error(), http.StatusBadRequest)
		}

		return purchaseJSON, nil
	},
}
