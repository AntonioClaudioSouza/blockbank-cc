package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/blockbank-cc/chaincode/utils"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var MakeWithdrawal = tx.Transaction{
	Tag:         "makeWithdrawal",
	Label:       "Make a withdrawal",
	Description: "Make a new withdrawal",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:         "holder",
			Label:       "Withdrawal holder",
			Description: "Withdrawal holder",
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
		withdrawValue, _ := req["value"].(float64)
		holderKey, _ := req["holder"].(assets.Key)

		if withdrawValue <= 0 {
			return nil, errors.WrapError(nil, "Withdraws values must be higher than zero")
		}

		holderMap, err := holderKey.GetMap(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}

		//UPDATE BALANCE
		holderMap, err = holderKey.Update(stub, map[string]interface{}{
			"cash": holderMap["cash"].(float64) - withdrawValue,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update receiver asset")
		}

		//WITHDRAWAL MAP
		timeStamp, _ := stub.Stub.GetTxTimestamp()
		withdrawalMap := make(map[string]interface{})
		withdrawalMap["@assetType"] = "withdrawal"
		withdrawalMap["value"] = withdrawValue
		withdrawalMap["holder"] = holderMap
		withdrawalMap["txId"] = stub.Stub.GetTxID()
		withdrawalMap["date"] = utils.ReturnDate(timeStamp)

		withdrawalAsset, err := assets.NewAsset(withdrawalMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create deposit asset")
		}

		withdrawalMap, err = withdrawalAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		withdrawalJSON, nerr := json.Marshal(withdrawalMap)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return withdrawalJSON, nil
	},
}
