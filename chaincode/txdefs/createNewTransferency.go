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

var CreateNewTransferency = tx.Transaction{
	Tag:         "createNewTransferency",
	Label:       "Create new transferency",
	Description: "Create a new transferency",
	Method:      "POST",

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

		transferencyValue, _ := req["value"].(float64)
		senderKey, _ := req["sender"].(assets.Key)
		receiverKey, _ := req["receiver"].(assets.Key)

		if senderKey.Key() == receiverKey.Key() {
			return nil, errors.WrapError(nil, "Sender and receiver must be different holders")
		}

		senderMap, err := senderKey.GetMap(stub)
		if err != nil {
			return nil, errors.NewCCError("failed to get owner from the ledger", http.StatusBadRequest)
		}

		if senderMap["cash"].(float64) < transferencyValue {
			return nil, errors.NewCCError("sender doesn't have enough balance to do this transferency", http.StatusBadRequest)
		}

		receiverMap, err := receiverKey.GetMap(stub)
		if err != nil {
			return nil, errors.NewCCError("failed to get owner from the ledger", http.StatusBadRequest)
		}

		//UPDATE BALANCES
		senderMap["cash"] = senderMap["cash"].(float64) - transferencyValue
		receiverMap["cash"] = receiverMap["cash"].(float64) + transferencyValue

		senderMap, err = senderKey.Update(stub, map[string]interface{}{
			"cash": senderMap["cash"].(float64) - transferencyValue,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update sender asset")
		}

		receiverMap, err = receiverKey.Update(stub, map[string]interface{}{
			"cash": receiverMap["cash"].(float64) + transferencyValue,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update receiver asset")
		}

		timeStamp, _ := stub.Stub.GetTxTimestamp()
		transferencyMap := make(map[string]interface{})
		transferencyMap["@assetType"] = "transferency"
		transferencyMap["txId"] = stub.Stub.GetTxID()
		transferencyMap["sender"] = senderMap
		transferencyMap["receiver"] = receiverMap
		transferencyMap["value"] = transferencyValue
		transferencyMap["date"] = utils.ReturnDate(timeStamp)

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

		return transferencyJSON, nil
	},
}
