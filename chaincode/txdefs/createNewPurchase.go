package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateNewPurchase = tx.Transaction{
	Tag:         "createNewPurchase",
	Label:       "Create new purchase",
	Description: "Create a new purchase",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

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
		purchaseDescription := req["description"].(string)
		purchaseValue := req["value"].(float64)

		buyerKey, ok := req["buyer"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter owner must be an asset")
		}

		buyerAsset, err := buyerKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}
		buyerMap := (map[string]interface{})(*buyerAsset)

		if buyerMap["cash"].(float64) < purchaseValue {
			return nil, errors.WrapError(err, "Buyer doesn't have enough balance to make this purchase")
		}

		//REMOVE BALANCE FROM BUYER
		buyerMap["cash"] = buyerMap["cash"].(float64) - purchaseValue

		purchaseMap := make(map[string]interface{})
		purchaseMap["@assetType"] = "purchase"
		purchaseMap["tdId"] = stub.Stub.GetTxID()
		purchaseMap["description"] = purchaseDescription
		purchaseMap["buyer"] = buyerMap
		purchaseMap["value"] = purchaseValue

		buyerMap, err = buyerAsset.Update(stub, buyerMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update buyer asset")
		}

		purchaseAsset, err := assets.NewAsset(purchaseMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create creditCard asset")
		}

		_, err = purchaseAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		purchaseJSON, nerr := json.Marshal(purchaseAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		// // Marshall message to be logged
		// logMsg, ok := json.Marshal(fmt.Sprintf("New library name: %s", name))
		// if ok != nil {
		// 	return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		// }

		// // Call event to log the message
		// events.CallEvent(stub, "createLibraryLog", logMsg)

		return purchaseJSON, nil
	},
}
