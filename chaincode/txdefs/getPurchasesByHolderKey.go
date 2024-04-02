package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetPurchasesByHolderKey = tx.Transaction{
	Tag:         "getPurchasesByHolderKey",
	Label:       "Get Purchases By Holder Key",
	Description: "Get Purchases By Holder Key",
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
				"@assetType": "purchase",
				"buyer.@key": holderKey.Key(),
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for purchases", 500)
		}

		purchasesJSON, _ := json.Marshal(response.Result)
		return purchasesJSON, nil
	},
}
