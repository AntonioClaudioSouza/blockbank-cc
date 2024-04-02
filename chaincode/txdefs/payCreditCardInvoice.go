package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/blockbank-cc/chaincode/utils"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var PayCreditCardInvoice = tx.Transaction{
	Tag:         "payCreditCardInvoice",
	Label:       "Pay credit card invoice",
	Description: "Pay credit card invoice",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:         "creditCard",
			Label:       "Credit card",
			Description: "Credit card",
			DataType:    "->creditCard",
			Required:    true,
		},
		{
			Tag:         "valueToPay",
			Label:       "Value to pay",
			Description: "Value to pay",
			DataType:    "number",
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		valueToPay, _ := req["valueToPay"].(float64)
		creditCardKey, _ := req["creditCard"].(assets.Key)

		creditCardMap, err := creditCardKey.GetMap(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get creditCard from the ledger")
		}

		ownerKey, _ := assets.NewKey(creditCardMap["owner"].(map[string]interface{}))
		ownerMap, err := ownerKey.GetMap(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}

		if valueToPay == 0 {
			valueToPay = creditCardMap["limitUsed"].(float64)
		}

		if valueToPay > creditCardMap["limitUsed"].(float64) {
			return nil, errors.WrapError(err, "Value to pay can't be higher than the used limit")
		}

		if valueToPay > ownerMap["cash"].(float64) {
			return nil, errors.WrapError(err, "Owner doesn't have enough balance")
		}

		//UPDATING DATA
		creditCardMap, err = creditCardKey.Update(stub, map[string]interface{}{
			"limitUsed": creditCardMap["limitUsed"].(float64) - valueToPay,
			"owner":     ownerMap,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update creditCard asset")
		}

		ownerMap, err = ownerKey.Update(stub, map[string]interface{}{
			"cash": ownerMap["cash"].(float64) - valueToPay,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update holder asset")
		}

		//CREATING PAYMENT ASSET
		invoicePaymentMap := make(map[string]interface{})

		timeStamp, _ := stub.Stub.GetTxTimestamp()
		invoicePaymentMap["@assetType"] = "invoicePayment"
		invoicePaymentMap["txId"] = stub.Stub.GetTxID()
		invoicePaymentMap["value"] = valueToPay
		invoicePaymentMap["owner"] = ownerMap
		invoicePaymentMap["creditCard"] = creditCardMap
		invoicePaymentMap["date"] = utils.ReturnDate(timeStamp)

		invoicePaymentAsset, err := assets.NewAsset(invoicePaymentMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create payment asset")
		}

		_, err = invoicePaymentAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving payment asset on blockchain")
		}

		creditCardJSON, nerr := json.Marshal(creditCardMap)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return creditCardJSON, nil
	},
}
