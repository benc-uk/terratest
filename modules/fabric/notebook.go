package fabric

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/fabric-sdk-go/fabric/notebook"
	"github.com/stretchr/testify/require"
)

func notebookClientFactory() (*notebook.ClientFactory, error) {
	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	cf, err := notebook.NewClientFactory(creds, nil, nil)
	if err != nil {
		return nil, err
	}

	return cf, nil
}

// NotebookExistsE checks if a notebook exists in a workspace.
func NotebookExistsE(wsID string, searchString string, st SearchType) (bool, error) {
	factory, err := notebookClientFactory()
	if err != nil {
		return false, err
	}

	client := factory.NewItemsClient()

	list, err := client.ListNotebooks(context.Background(), wsID, nil)
	if err != nil {
		return false, err
	}

	var foundItem *notebook.Notebook
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
		return false, errors.New("notebook could not be found with name " + searchString)
	}

	return true, nil
}

func NotebookExists(t *testing.T, wsID string, searchString string, st SearchType) bool {
	exists, err := NotebookExistsE(wsID, searchString, st)
	require.NoError(t, err)

	return exists
}
