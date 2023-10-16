package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateNewTransferency = tx.Transaction{
	Tag:         "createNewTransferency",
	Label:       "Create new transferency",
	Description: "Create a new transferency",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

	Args: []tx.Argument{

		{
			Tag:         "sender",
			Label:       "Transferency sender",
			Description: "Who's sending the money",
			DataType:    "->holder",
			Required:    true,
		},
		{
			Tag:         "receiver",
			Label:       "Transferency receiver",
			Description: "Who's receiving the money",
			DataType:    "->holder",
			Required:    true,
		},
		{
			Tag:         "value",
			Label:       "Transferency value",
			Description: "Transferency value",
			DataType:    "number",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		transferencyValue := req["value"].(float64)
		senderKey, ok := req["sender"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter sender must be an asset")
		}
		receiverKey, ok := req["receiver"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter receiver must be an asset")
		}

		if senderKey.Key() == receiverKey.Key() {
			return nil, errors.WrapError(nil, "Sender and receiver must be different holders")
		}

		senderAsset, err := senderKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}
		senderMap := (map[string]interface{})(*senderAsset)

		if senderMap["cash"].(float64) < transferencyValue {
			return nil, errors.WrapError(err, "Sender doesn't have enough balance to do this transferency")
		}

		receiverAsset, err := receiverKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}
		receiverMap := (map[string]interface{})(*receiverAsset)

		transferencyMap := make(map[string]interface{})
		transferencyMap["@assetType"] = "transferency"
		transferencyMap["sender"] = senderMap
		transferencyMap["receiver"] = receiverMap
		transferencyMap["value"] = transferencyValue

		//UPDATE BALANCES
		senderMap["cash"] = senderMap["cash"].(float64) - transferencyValue
		receiverMap["cash"] = receiverMap["cash"].(float64) + transferencyValue

		senderMap, err = senderAsset.Update(stub, senderMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update sender asset")
		}

		receiverMap, err = receiverAsset.Update(stub, receiverMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update receiver asset")
		}

		transferencyAsset, err := assets.NewAsset(transferencyMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create transferency asset")
		}

		_, err = transferencyAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving transferency asset on blockchain")
		}

		transferencyJSON, nerr := json.Marshal(transferencyAsset)
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

		return transferencyJSON, nil
	},
}
