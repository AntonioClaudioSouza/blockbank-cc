package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var ListManagers = tx.Transaction{
	Tag:         "listManagers",
	Label:       "List managers",
	Description: "List managers",
	Method:      "POST",

	Args: []tx.Argument{},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "manager",
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for managers", 500)
		}

		managersJSON, _ := json.Marshal(response.Result)
		return managersJSON, nil
	},
}
