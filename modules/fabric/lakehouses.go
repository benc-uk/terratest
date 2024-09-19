package fabric

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/fabric-sdk-go/fabric/lakehouse"
	"github.com/stretchr/testify/require"
)

func lakehouseClientFactory() (*lakehouse.ClientFactory, error) {
	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	cf, err := lakehouse.NewClientFactory(creds, nil, nil)
	if err != nil {
		return nil, err
	}

	return cf, nil
}

func LakehouseExistsE(wsID string, searchString string, st SearchType) (bool, error) {
	if st != SearchByID && st != SearchByDisplayName {
		return false, errors.New("invalid search type")
	}

	factory, err := lakehouseClientFactory()
	if err != nil {
		return false, err
	}

	client := factory.NewItemsClient()

	list, err := client.ListLakehouses(context.Background(), wsID, nil)
	if err != nil {
		return false, err
	}

	var foundItem *lakehouse.Lakehouse
	for _, item := range list {
		if st == SearchByID {
			if *item.ID == searchString {
				foundItem = &item
				break
			}
		} else {
			if *item.DisplayName == searchString {
				foundItem = &item
				break
			}
		}
	}

	if foundItem == nil {
		return false, errors.New("lakehouse could not be found with name " + searchString)
	}

	return true, nil
}

func LakehouseExists(t *testing.T, wsID string, searchString string, st SearchType) bool {
	exists, err := LakehouseExistsE(wsID, searchString, st)
	require.NoError(t, err)

	return exists
}
