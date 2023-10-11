package txdefs

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

func generateValidCreditCardNumber(timestamp *timestamppb.Timestamp) string {
	rand.Seed(timestamp.Seconds)
	prefix := "4"

	var numberDigits []int
	for i := 0; i < 14; i++ {
		numberDigits = append(numberDigits, rand.Intn(10))
	}

	sum := 0
	double := false
	for i := len(numberDigits) - 1; i >= 0; i-- {
		digit := numberDigits[i]
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	checksum := (sum * 9) % 10

	numberDigits = append(numberDigits, checksum)
	cardNumber := prefix + digitsToString(numberDigits)
	return cardNumber
}

func digitsToString(digits []int) string {
	result := ""
	for _, digit := range digits {
		result += fmt.Sprint(digit)
	}
	return result
}

var CreateNewCreditCard = tx.Transaction{
	Tag:         "createNewCreditCard",
	Label:       "Create new credit card",
	Description: "Create a new credit card",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

	Args: []tx.Argument{

		{
			Tag:         "owner",
			Label:       "Credit card owner",
			Description: "Credit card owner",
			DataType:    "->holder",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		timeStamp, _ := stub.Stub.GetTxTimestamp()
		number := generateValidCreditCardNumber(timeStamp)

		ownerKey, ok := req["owner"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter owner must be an asset")
		}

		// Prepare couchdb query
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "creditCard",
				"owner":      ownerKey,
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for credit card", 500)
		}
		if len(response.Result) >= 1 {
			return nil, errors.WrapErrorWithStatus(err, "Holder already have a credit card", 500)
		}

		ownerAsset, err := ownerKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}
		ownerMap := (map[string]interface{})(*ownerAsset)

		if !ownerMap["ccAvailable"].(bool) {
			return nil, errors.WrapError(err, "Holder doesn't have a credit card available")
		}

		ccMap := make(map[string]interface{})
		ccMap["@assetType"] = "creditCard"
		ccMap["number"] = number
		ccMap["owner"] = ownerMap

		ccAsset, err := assets.NewAsset(ccMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create creditCard asset")
		}

		_, err = ccAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		ccJSON, nerr := json.Marshal(ccAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		// // Marshall message to be logged
		// logMsg, ok := json.Marshal(fmt.Sprintf("New library name: %s", name))
		// if ok != nil {
		// 	return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		// }

		// // Call event to log the message
		// events.CallEvent(stub, "createLibraryLog", logMsg)

		return ccJSON, nil
	},
}
