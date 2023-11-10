package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetPaymentsByCreditCardKey = tx.Transaction{
	Tag:         "getPaymentsByCreditCardKey",
	Label:       "Get Payments By CreditCard Key",
	Description: "Get Payments By CreditCard Key",
	Method:      "GET",
	Callers:     []string{"$orgMSP"},

	Args: []tx.Argument{
		{
			Tag:         "creditCard",
			Label:       "Credit Card",
			Description: "Credit Card",
			DataType:    "->creditCard",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		creditCardKey, ok := req["creditCard"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter creditCard must be an asset")
		}

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "invoicePayment",
				"creditCard": creditCardKey,
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for creditCard payments", 500)
		}

		purchasesJSON, err := json.Marshal(response.Result)

		return purchasesJSON, nil
	},
}
