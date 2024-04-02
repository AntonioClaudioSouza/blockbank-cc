package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var ListHolders = tx.Transaction{
	Tag:         "listHolders",
	Label:       "List holders",
	Description: "List holders",
	Method:      "GET",

	Args: []tx.Argument{},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "holder",
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for holders", 500)
		}

		holdersJSON, _ := json.Marshal(response.Result)
		return holdersJSON, nil
	},
}
