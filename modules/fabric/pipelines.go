package fabric

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/fabric-sdk-go/fabric/datapipeline"
	"github.com/stretchr/testify/require"
)

func pipelineClientFactory() (*datapipeline.ClientFactory, error) {
	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	cf, err := datapipeline.NewClientFactory(creds, nil, nil)
	if err != nil {
		return nil, err
	}

	return cf, nil
}

// PipelineExistsE checks if a pipeline exists in a workspace and returns errors if necessary
func PipelineExistsE(wsID string, searchString string, st SearchType) (bool, error) {
	factory, err := pipelineClientFactory()
	if err != nil {
		return false, err
	}

	client := factory.NewItemsClient()

	list, err := client.ListDataPipelines(context.Background(), wsID, nil)
	if err != nil {
		return false, err
	}

	var foundPl *datapipeline.DataPipeline
	for _, pl := range list {
		if st == SearchByID {
			if *pl.ID == searchString {
				foundPl = &pl
				break
			}
		} else {
			if *pl.DisplayName == searchString {
				foundPl = &pl
				break
			}
		}
	}

	if foundPl == nil {
		return false, errors.New("pipeline could not be found with name " + searchString)
	}

	return true, nil
}

// PipelineExists checks if a pipeline exists in a workspace, failing the test if an error occurs
func PipelineExists(t *testing.T, wsID string, searchString string, st SearchType) bool {
	exists, err := PipelineExistsE(wsID, searchString, st)
	require.NoError(t, err)

	return exists
}
