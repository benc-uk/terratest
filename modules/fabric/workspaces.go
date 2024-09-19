package fabric

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/fabric-sdk-go/fabric/core"
	"github.com/stretchr/testify/require"
)

func coreClientFactory() (*core.ClientFactory, error) {
	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	cf, err := core.NewClientFactory(creds, nil, nil)
	if err != nil {
		return nil, err
	}

	return cf, nil
}

func WorkspaceExistsE(searchString string, st SearchType) (bool, error) {
	if st != SearchByID && st != SearchByDisplayName {
		return false, errors.New("invalid search type")
	}

	factory, err := coreClientFactory()
	if err != nil {
		return false, err
	}

	client := factory.NewWorkspacesClient()

	list, err := client.ListWorkspaces(context.Background(), nil)
	if err != nil {
		return false, err
	}

	var foundWS *core.Workspace
	for _, ws := range list {
		if st == SearchByID {
			if *ws.ID == searchString {
				foundWS = &ws
				break
			}
		} else {
			if *ws.DisplayName == searchString {
				foundWS = &ws
				break
			}
		}
	}

	if foundWS == nil {
		return false, errors.New("workspace could not be found with name " + searchString)
	}

	return true, nil
}

func WorkspaceExists(t *testing.T, searchString string, st SearchType) bool {
	exists, err := WorkspaceExistsE(searchString, st)
	require.NoError(t, err)

	return exists
}
