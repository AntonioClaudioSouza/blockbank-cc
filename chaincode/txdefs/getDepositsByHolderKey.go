package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetDepositsByHolderKey = tx.Transaction{
	Tag:         "getDepositsByHolderKey",
	Label:       "Get Deposits By Holder Key",
	Description: "Get Deposits By Holder Key",
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
		holderKey, ok := req["holder"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter holder must be an asset")
		}
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType":  "deposit",
				"holder.@key": holderKey.Key(),
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for deposits", 500)
		}

		depositsJSON, _ := json.Marshal(response.Result)
		return depositsJSON, nil
	},
}
