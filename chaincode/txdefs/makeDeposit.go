package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/blockbank-cc/chaincode/utils"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var MakeDeposit = tx.Transaction{
	Tag:         "makeDeposit",
	Label:       "Make deposit",
	Description: "Make a new deposit",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:         "receiver",
			Label:       "Deposit receiver",
			Description: "Deposit receiver",
			DataType:    "->holder",
			Required:    true,
		},
		{
			Tag:         "value",
			Label:       "Deposit value",
			Description: "Deposit value",
			DataType:    "number",
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		depositValue := req["value"].(float64)
		receiverKey, _ := req["receiver"].(assets.Key)

		if depositValue <= 0 {
			return nil, errors.WrapError(nil, "Deposit values must be higher than zero")
		}

		receiverMap, err := receiverKey.GetMap(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}

		//UPDATE BALANCE
		receiverMap, err = receiverKey.Update(stub, map[string]interface{}{
			"cash": receiverMap["cash"].(float64) + depositValue,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update receiver asset")
		}

		//DEPOSIT MAP
		timeStamp, _ := stub.Stub.GetTxTimestamp()
		depositMap := make(map[string]interface{})
		depositMap["@assetType"] = "deposit"
		depositMap["value"] = depositValue
		depositMap["holder"] = receiverMap
		depositMap["txId"] = stub.Stub.GetTxID()
		depositMap["date"] = utils.ReturnDate(timeStamp)

		depositAsset, err := assets.NewAsset(depositMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create deposit asset")
		}

		depositMap, err = depositAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		depositJSON, nerr := json.Marshal(depositMap)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}
		return depositJSON, nil
	},
}
