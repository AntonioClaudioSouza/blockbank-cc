package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateNewManager = tx.Transaction{
	Tag:         "createNewManager",
	Label:       "Create new manager",
	Description: "Create a new manager",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

	Args: []tx.Argument{
		{
			Tag:         "name",
			Label:       "Name",
			Description: "Manager's name",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "document",
			Label:       "ID Document",
			Description: "Manager's id document",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		name, _ := req["name"].(string)
		document, _ := req["document"].(string)

		managerMap := make(map[string]interface{})
		managerMap["@assetType"] = "manager"
		managerMap["name"] = name
		managerMap["document"] = document

		managerAsset, err := assets.NewAsset(managerMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create manager asset")
		}

		_, err = managerAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		managerJSON, nerr := json.Marshal(managerAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return managerJSON, nil
	},
}
