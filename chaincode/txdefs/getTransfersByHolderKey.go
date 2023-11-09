package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

func checkSent(sentBool bool, holder assets.Key) map[string]interface{} {
	if sentBool {
		return map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "transferency",
				"sender":     holder,
			},
		}
	} else {
		return map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "transferency",
				"receiver":   holder,
			},
		}
	}
}

var GetTransfersByHolderKey = tx.Transaction{
	Tag:         "getTransfersByHolderKey",
	Label:       "Get Transfers By Holder Key",
	Description: "Get Transfers By Holder Key",
	Method:      "GET",
	Callers:     []string{"$orgMSP"},

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
		sent := req["sent"].(bool)
		holderKey, ok := req["holder"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter holder must be an asset")
		}

		query := checkSent(sent, holderKey)

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for transfers", 500)
		}

		transfersJSON, err := json.Marshal(response.Result)

		return transfersJSON, nil
	},
}
