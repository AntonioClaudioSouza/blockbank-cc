package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

func checkSent(sentBool bool, holder assets.Key) map[string]interface{} {
	if sentBool {
		return map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType":  "transferency",
				"sender.@key": holder.Key(),
			},
		}
	} else {
		return map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType":    "transferency",
				"receiver.@key": holder.Key(),
			},
		}
	}
}

var GetTransfersByHolderKey = tx.Transaction{
	Tag:         "getTransfersByHolderKey",
	Label:       "Get Transfers By Holder Key",
	Description: "Get Transfers By Holder Key",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Tag:         "sent",
			Label:       "Sent",
			Description: "Transfers Sent",
			DataType:    "boolean",
			Required:    true,
		},
		{
			Tag:         "holder",
			Label:       "Holder",
			Description: "Holder",
			DataType:    "->holder",
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		sent, _ := req["sent"].(bool)
		holderKey, _ := req["holder"].(assets.Key)

		query := checkSent(sent, holderKey)

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for transfers", 500)
		}

		transfersJSON, _ := json.Marshal(response.Result)
		return transfersJSON, nil
	},
}
