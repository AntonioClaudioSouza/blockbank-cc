package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var MakeDeposit = tx.Transaction{
	Tag:         "makeDeposit",
	Label:       "Make deposit",
	Description: "Make a new deposit",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

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
		receiverKey, ok := req["receiver"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter receiver must be an asset")
		}

		if depositValue <= 0 {
			return nil, errors.WrapError(nil, "Deposit values must be higher than zero")
		}

		receiverAsset, err := receiverKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}
		receiverMap := (map[string]interface{})(*receiverAsset)

		//UPDATE BALANCE
		receiverMap["cash"] = receiverMap["cash"].(float64) + depositValue

		receiverMap, err = receiverAsset.Update(stub, receiverMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update receiver asset")
		}

		receiverJSON, nerr := json.Marshal(receiverAsset)
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

		return receiverJSON, nil
	},
}
