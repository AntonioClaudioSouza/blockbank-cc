package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetHolderByKey = tx.Transaction{
	Tag:         "getHolderByKey",
	Label:       "Get Holder By Key",
	Description: "Get Holder By Key",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Tag:         "holder",
			Label:       "Holder",
			Description: "Holder",
			DataType:    "->holder",
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		holderKey, _ := req["holder"].(assets.Key)
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "holder",
				"@key":       holderKey.Key(),
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for holder", 500)
		}

		if len(response.Result) > 0 {
			holderJSON, _ := json.Marshal(response.Result[0])
			return holderJSON, nil
		}
		return nil, errors.WrapErrorWithStatus(nil, "Holder not found", 500)
	},
}
