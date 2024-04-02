package txdefs

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var ActivateCreditCard = tx.Transaction{
	Tag:         "activateCreditCard",
	Label:       "Activate credit card",
	Description: "Activate cretit card to a holder",
	Method:      "POST",

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
		var err error
		var ownerMap map[string]interface{}

		// ** Get ownerMap by asset keys
		ownerKey, _ := req["owner"].(assets.Key)

		if ownerMap, err = ownerKey.GetMap(stub); err != nil {
			return nil, errors.NewCCError(err.Error(), http.StatusBadRequest)
		}

		// ** Check if credCard is activate for owner
		if ownerMap["ccAvailable"].(bool) {
			return nil, errors.NewCCError("Holder credit card is already available", http.StatusBadRequest)
		}

		// ** Active credCard
		if ownerMap, err = ownerKey.Update(stub, map[string]interface{}{"ccAvailable": true}); err != nil {
			return nil, errors.WrapError(err, "Error updating holder asset")
		}

		ownerJSON, nerr := json.Marshal(ownerMap)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return ownerJSON, nil
	},
}
