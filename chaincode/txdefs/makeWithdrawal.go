package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools-demo/chaincode/utils"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var MakeWithdrawal = tx.Transaction{
	Tag:         "makeWithdrawal",
	Label:       "Make a withdrawal",
	Description: "Make a new withdrawal",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

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

		timeStamp, err := stub.Stub.GetTxTimestamp()
		withdrawValue := req["value"].(float64)
		holderKey, ok := req["holder"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter receiver must be an asset")
		}

		if withdrawValue <= 0 {
			return nil, errors.WrapError(nil, "Withdraws values must be higher than zero")
		}

		holderAsset, err := holderKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}
		holderMap := (map[string]interface{})(*holderAsset)

		//UPDATE BALANCE
		holderMap["cash"] = holderMap["cash"].(float64) - withdrawValue

		holderMap, err = holderAsset.Update(stub, holderMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update receiver asset")
		}

		//WITHDRAWAL MAP
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

		withdrawalJSON, nerr := json.Marshal(withdrawalAsset)
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

		return withdrawalJSON, nil
	},
}
