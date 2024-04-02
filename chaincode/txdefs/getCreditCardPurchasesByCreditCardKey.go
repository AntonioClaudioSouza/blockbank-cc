package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetCreditCardPurchasesByCreditCardKey = tx.Transaction{
	Tag:         "getCreditCardPurchasesByCreditCardKey",
	Label:       "Get CreditCard Purchases By CreditCard Key",
	Description: "Get CreditCard Purchases By CreditCard Key",
	Method:      "GET",

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
		creditCardKey, _ := req["creditCard"].(assets.Key)

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType":      "creditCardPurchase",
				"creditCard.@Key": creditCardKey,
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for creditCard purchases", 500)
		}

		purchasesJSON, _ := json.Marshal(response.Result)
		return purchasesJSON, nil
	},
}
