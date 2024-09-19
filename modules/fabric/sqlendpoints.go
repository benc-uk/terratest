package fabric

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/fabric-sdk-go/fabric/sqlendpoint"
	"github.com/stretchr/testify/require"
)

func sqlendpointClientFactory() (*sqlendpoint.ClientFactory, error) {
	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	cf, err := sqlendpoint.NewClientFactory(creds, nil, nil)
	if err != nil {
		return nil, err
	}

	return cf, nil
}

// SQLEndpointExistsE checks if a SQL endpoint exists in a workspace and returns errors if necessary
func SQLEndpointExistsE(wsID string, searchString string, st SearchType) (bool, error) {
	factory, err := sqlendpointClientFactory()
	if err != nil {
		return false, err
	}

	client := factory.NewItemsClient()

	list, err := client.ListSQLEndpoints(context.Background(), wsID, nil)
	if err != nil {
		return false, err
	}

	var foundItem *sqlendpoint.SQLEndpoint
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
		return false, errors.New("sql endpoint could not be found with name " + searchString)
	}

	return true, nil
}

// SQLEndpointExists checks if a SQL endpoint exists in a workspace, failing the test if an error occurs
func SQLEndpointExists(t *testing.T, wsID string, searchString string, st SearchType) bool {
	exists, err := SQLEndpointExistsE(wsID, searchString, st)
	require.NoError(t, err)

	return exists
}
