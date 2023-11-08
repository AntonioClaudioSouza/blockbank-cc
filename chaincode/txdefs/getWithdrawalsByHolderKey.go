package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetWithdrawalsByHolderKey = tx.Transaction{
	Tag:         "getWithdrawalsByHolderKey",
	Label:       "Get Withdrawals By Holder Key",
	Description: "Get Withdrawals By Holder Key",
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
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "withdrawal",
				"holder":     holderKey,
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for deposits", 500)
		}

		withdrawalsJSON, err := json.Marshal(response.Result)

		return withdrawalsJSON, nil
	},
}
