package txdefs

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	Args: []tx.Argument{
		{
			Tag:         "owner",
			Label:       "Credit card owner",
			Description: "Credit card owner",
			DataType:    "->holder",
			Required:    true,
		},
		{
			Tag:         "creditCardName",
			Label:       "Credit card name",
			Description: "Credit card name",
			DataType:    "string",
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		ownerKey, _ := req["owner"].(assets.Key)
		creditCardName := req["creditCardName"].(string)

		// ** Prepare couchdb query for checking if the user already has a credit card
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "creditCard",
				"owner":      ownerKey,
			},
		}

		if response, err := assets.Search(stub, query, "", true); err != nil {
			return nil, errors.NewCCError("error searching for credit card", http.StatusBadRequest)
		} else {

			// ** If exist credcard with this owner - return error
			if len(response.Result) >= 1 {
				return nil, errors.NewCCError("holder already have a credit card", http.StatusBadRequest)
			}
		}

		ownerMap, err := ownerKey.GetMap(stub)
		if err != nil {
			return nil, errors.NewCCError("failed to get owner from the ledger", http.StatusBadRequest)
		}

		if !ownerMap["ccAvailable"].(bool) {
			return nil, errors.NewCCError("holder doesn't have a credit card available", http.StatusBadRequest)
		}

		timeStamp, _ := stub.Stub.GetTxTimestamp()
		ccAsset, err := assets.NewAsset(map[string]interface{}{
			"@assetType":     "creditCard",
			"number":         generateValidCreditCardNumber(timeStamp),
			"owner":          ownerMap,
			"creditCardName": creditCardName,
		})
		if err != nil {
			return nil, errors.NewCCError("failed to create creditCard asset :"+err.Error(), http.StatusBadRequest)
		}

		ccAssetMap, err := ccAsset.PutNew(stub)
		if err != nil {
			return nil, errors.NewCCError("error saving asset on blockchain :"+err.Error(), http.StatusBadRequest)
		}

		ccJSON, nerr := json.Marshal(ccAssetMap)
		if nerr != nil {
			return nil, errors.NewCCError("failed to encode asset to JSON format :"+nerr.Error(), http.StatusBadRequest)
		}

		return ccJSON, nil
	},
}
