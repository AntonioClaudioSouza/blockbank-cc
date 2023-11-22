package txdefs

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetCreditCardByHolderKey = tx.Transaction{
	Tag:         "getCreditCardByHolderKey",
	Label:       "Get Credit Card By Holder Key",
	Description: "Get Credit Card By Holder Key",
	Method:      "GET",
	Callers:     []string{"$orgMSP"},

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

		fmt.Println("Cheguei aqui")
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "creditCard",
				"owner":      holderKey,
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for credit card key", 500)
		}

		if len(response.Result) > 0 {
			creditCardJSON, _ := json.Marshal(response.Result[0])
			return creditCardJSON, nil
		}
		return nil, errors.WrapErrorWithStatus(nil, "Credit card not found", 500)
	},
}
