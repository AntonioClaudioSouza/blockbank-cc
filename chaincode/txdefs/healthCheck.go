package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var HealthCheck = tx.Transaction{
	Tag:         "healthCheck",
	Label:       "Health Check",
	Description: "Check if CCAPI is working",
	Method:      "GET",
	Callers:     []string{"$orgMSP"},

	Args: []tx.Argument{},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		ok := map[string]interface{}{
			"status": "ok",
		}

		okJSON, _ := json.Marshal(ok)

		return okJSON, nil
	},
}
