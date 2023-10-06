//go:build integration

package experiment

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/hetiansu5/urlquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/response"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/tests/integration/golang/fixtures"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type SearchExperimentsTestSuite struct {
	suite.Suite
	client             *helpers.HttpClient
	experimentFixtures *fixtures.ExperimentFixtures
}

func TestSearchExperimentsTestSuite(t *testing.T) {
	suite.Run(t, new(SearchExperimentsTestSuite))
}

func (s *SearchExperimentsTestSuite) SetupTest() {
	s.client = helpers.NewMlflowApiClient(helpers.GetServiceUri())
	experimentFixtures, err := fixtures.NewExperimentFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.experimentFixtures = experimentFixtures
}

func (s *SearchExperimentsTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.experimentFixtures.UnloadFixtures())
	}()
	// 1. prepare database with test data.
	experiments := []models.Experiment{
		{
			Name: "Test Experiment 1",
			Tags: []models.ExperimentTag{
				{
					Key:   "key2",
					Value: "value2",
				},
			},

			LifecycleStage: models.LifecycleStageActive,
		},
		{
			Name: "Test Experiment 2",
			Tags: []models.ExperimentTag{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
			LifecycleStage: models.LifecycleStageActive,
		},
		{
			Name: "Test Experiment 3",
			Tags: []models.ExperimentTag{
				{
					Key:   "key3",
					Value: "value3",
				},
			},
			LifecycleStage: models.LifecycleStageActive,
		},
		{
			Name: "Test Experiment 4",
			Tags: []models.ExperimentTag{
				{
					Key:   "key4",
					Value: "value4",
				},
			},
			LifecycleStage: models.LifecycleStageActive,
		},
		{
			Name:           "Test Experiment 5",
			Tags:           nil,
			LifecycleStage: models.LifecycleStageActive,
		},
		{
			Name:           "Test Experiment 6",
			Tags:           nil,
			LifecycleStage: models.LifecycleStageDeleted,
		},
	}
	for _, ex := range experiments {
		_, err := s.experimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
			Name: ex.Name,
			Tags: ex.Tags,
			CreationTime: sql.NullInt64{
				Int64: time.Now().UTC().UnixMilli(),
				Valid: true,
			},
			LastUpdateTime: sql.NullInt64{
				Int64: time.Now().UTC().UnixMilli(),
				Valid: true,
			},
			LifecycleStage:   ex.LifecycleStage,
			ArtifactLocation: "/artifact/location",
		})
		assert.Nil(s.T(), err)
	}

	tests := []struct {
		name     string
		request  request.SearchExperimentsRequest
		expected []string
	}{
		{
			name: "TestFilter",
			request: request.SearchExperimentsRequest{
				Filter: "attribute.name != 'Test Experiment 5'",
			},
			expected: []string{"Test Experiment 1", "Test Experiment 2", "Test Experiment 3", "Test Experiment 4", "Default"},
		},
		{
			name: "TestViewType",
			request: request.SearchExperimentsRequest{
				ViewType: request.ViewTypeDeletedOnly,
			},
			expected: []string{"Test Experiment 6"},
		},
		{
			name: "TestOrderBy",
			request: request.SearchExperimentsRequest{
				OrderBy: []string{"name ASC"},
			},
			expected: []string{
				"Default",
				"Test Experiment 1",
				"Test Experiment 2",
				"Test Experiment 3",
				"Test Experiment 4",
				"Test Experiment 5",
			},
		},
		{
			name: "TestMaxResults",
			request: request.SearchExperimentsRequest{
				OrderBy:    []string{"name ASC"},
				MaxResults: 3,
			},
			expected: []string{"Default", "Test Experiment 1", "Test Experiment 2"},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			query, err := urlquery.Marshal(tt.request)
			assert.Nil(t, err)

			resp := response.SearchExperimentsResponse{}
			err = s.client.DoGetRequest(
				fmt.Sprintf("%s%s?%s", mlflow.ExperimentsRoutePrefix, mlflow.ExperimentsSearchRoute, query),
				&resp,
			)
			assert.Nil(t, err)

			names := make([]string, len(resp.Experiments))
			for i, exp := range resp.Experiments {
				names[i] = exp.Name
			}

			assert.ElementsMatch(t, tt.expected, names)
		})
	}
}

func (s *SearchExperimentsTestSuite) Test_Error() {
	testData := []struct {
		name    string
		error   *api.ErrorResponse
		request *request.SearchExperimentsRequest
	}{
		{
			name:  "InvalidViewType",
			error: api.NewInvalidParameterValueError("Invalid view_type 'invalid_ViewType'"),
			request: &request.SearchExperimentsRequest{
				ViewType: "invalid_ViewType",
			},
		},
		{
			name:  "InvalidMaxResult",
			error: api.NewInvalidParameterValueError("Invalid value for parameter 'max_results' supplied."),
			request: &request.SearchExperimentsRequest{
				MaxResults: 10000000,
			},
		},
		{
			name:  "InvalidFilterValue",
			error: api.NewInvalidParameterValueError("invalid numeric value 'cc'"),
			request: &request.SearchExperimentsRequest{
				Filter: "attribute.creation_time > cc",
			},
		},
		{
			name:  "MalformedFilter",
			error: api.NewInvalidParameterValueError("malformed filter 'invalid_filter'"),
			request: &request.SearchExperimentsRequest{
				Filter: "invalid_filter",
			},
		},
		{
			name:  "InvalidNumericValue",
			error: api.NewInvalidParameterValueError("invalid numeric value 'invalid_value'"),
			request: &request.SearchExperimentsRequest{
				Filter: "creation_time > invalid_value",
			},
		},
		{
			name:  "InvalidStringOperator",
			error: api.NewInvalidParameterValueError("invalid string attribute comparison operator '<'"),
			request: &request.SearchExperimentsRequest{
				Filter: "attribute.name < 'value'",
			},
		},
		{
			name:  "InvalidTagOperator",
			error: api.NewInvalidParameterValueError("invalid tag comparison operator '<'"),
			request: &request.SearchExperimentsRequest{
				Filter: "tag.value < 'value'",
			},
		},
		{
			name:  "InvalidEntity",
			error: api.NewInvalidParameterValueError("invalid entity type 'invalid_entity'. Valid values are ['tag', 'attribute']"),
			request: &request.SearchExperimentsRequest{
				Filter: "invalid_entity.name = value",
			},
		},
		{
			name:  "InvalidOrderByAttribute",
			error: api.NewInvalidParameterValueError("invalid attribute 'invalid_attribute'. Valid values are ['name', 'experiment_id', 'creation_time', 'last_update_time']"),
			request: &request.SearchExperimentsRequest{
				OrderBy: []string{"invalid_attribute"},
			},
		},
	}

	for _, tt := range testData {
		s.T().Run(tt.name, func(t *testing.T) {
			query, err := urlquery.Marshal(tt.request)
			assert.Nil(s.T(), err)
			resp := api.ErrorResponse{}
			err = s.client.DoGetRequest(
				fmt.Sprintf("%s%s?%s", mlflow.ExperimentsRoutePrefix, mlflow.ExperimentsSearchRoute, query),
				&resp,
			)
			assert.Nil(t, err)
			assert.Equal(s.T(), tt.error.Error(), resp.Error())
		})
	}
}
