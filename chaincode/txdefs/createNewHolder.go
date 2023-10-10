package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateNewHolder = tx.Transaction{
	Tag:         "createNewHolder",
	Label:       "Create new holder",
	Description: "Create a new holder",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

	Args: []tx.Argument{
		{
			Tag:         "name",
			Label:       "Name",
			Description: "Holder's name",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "document",
			Label:       "ID Document",
			Description: "Holder's id document",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:      "cash",
			Label:    "Cash",
			DataType: "number",
		},
		{
			Tag:      "ccAvailable",
			Label:    "Credit Card Available",
			DataType: "boolean",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		name, _ := req["name"].(string)
		document, _ := req["document"].(string)
		cash, _ := req["cash"].(float64)
		ccAvailable, _ := req["ccAvailable"].(bool)

		holderMap := make(map[string]interface{})
		holderMap["@assetType"] = "holder"
		holderMap["name"] = name
		holderMap["document"] = document
		holderMap["cash"] = cash
		holderMap["ccAvailable"] = ccAvailable

		holderAsset, err := assets.NewAsset(holderMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create holder asset")
		}

		_, err = holderAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		holderJSON, nerr := json.Marshal(holderAsset)
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

		return holderJSON, nil
	},
}
