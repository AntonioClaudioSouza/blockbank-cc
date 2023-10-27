package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var ActivateCreditCard = tx.Transaction{
	Tag:         "activateCreditCard",
	Label:       "Activate credit card",
	Description: "Activate cretit card to a holder",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

	Args: []tx.Argument{

		{
			Tag:         "owner",
			Label:       "Credit card owner",
			Description: "Credit card owner",
			DataType:    "->holder",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		ownerKey, ok := req["owner"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter owner must be an asset")
		}

		ownerAsset, err := ownerKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}
		ownerMap := (map[string]interface{})(*ownerAsset)

		if ownerMap["ccAvailable"].(bool) {
			return nil, errors.WrapError(err, "Holder credit card is already available")
		}

		ownerMap["ccAvailable"] = true

		ownerMap, err = ownerAsset.Update(stub, ownerMap)
		if err != nil {
			return nil, errors.WrapError(err, "Error updating holder asset")
		}

		ownerJSON, nerr := json.Marshal(ownerAsset)
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

		return ownerJSON, nil
	},
}
