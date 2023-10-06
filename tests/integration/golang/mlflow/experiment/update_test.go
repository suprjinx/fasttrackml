//go:build integration

package experiment

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/tests/integration/golang/fixtures"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type UpdateExperimentTestSuite struct {
	suite.Suite
	client             *helpers.HttpClient
	experimentFixtures *fixtures.ExperimentFixtures
}

func TestUpdateExperimentTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateExperimentTestSuite))
}

func (s *UpdateExperimentTestSuite) SetupTest() {
	s.client = helpers.NewMlflowApiClient(helpers.GetServiceUri())
	experimentFixtures, err := fixtures.NewExperimentFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.experimentFixtures = experimentFixtures
}

func (s *UpdateExperimentTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.experimentFixtures.UnloadFixtures())
	}()
	// 1. prepare database with test data.
	experiment, err := s.experimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name: "Test Experiment",
		CreationTime: sql.NullInt64{
			Int64: time.Now().UTC().UnixMilli(),
			Valid: true,
		},
		LastUpdateTime: sql.NullInt64{
			Int64: time.Now().UTC().UnixMilli(),
			Valid: true,
		},
		LifecycleStage:   models.LifecycleStageActive,
		ArtifactLocation: "/artifact/location",
	})
	assert.Nil(s.T(), err)

	req := request.UpdateExperimentRequest{
		ID:   fmt.Sprintf("%d", *experiment.ID),
		Name: "Test Updated Experiment",
	}
	err = s.client.DoPostRequest(
		fmt.Sprintf("%s%s", mlflow.ExperimentsRoutePrefix, mlflow.ExperimentsUpdateRoute),
		req,
		&struct{}{},
	)
	assert.Nil(s.T(), err)

	exp, err := s.experimentFixtures.GetExperimentByID(context.Background(), *experiment.ID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "Test Updated Experiment", exp.Name)
}

func (s *UpdateExperimentTestSuite) Test_Error() {
	testData := []struct {
		name    string
		error   *api.ErrorResponse
		request *request.UpdateExperimentRequest
	}{
		{
			name:  "EmptyIDProperty",
			error: api.NewInvalidParameterValueError("Missing value for required parameter 'experiment_id'"),
			request: &request.UpdateExperimentRequest{
				ID: "",
			},
		},
		{
			name:  "EmptyNameProperty",
			error: api.NewInvalidParameterValueError("Missing value for required parameter 'new_name'"),
			request: &request.UpdateExperimentRequest{
				ID:   "1",
				Name: "",
			},
		},
		{
			name: "InvalidIDFormat",
			error: api.NewBadRequestError(
				`unable to parse experiment id 'invalid_id': strconv.ParseInt: parsing "invalid_id": invalid syntax`,
			),
			request: &request.UpdateExperimentRequest{
				ID:   "invalid_id",
				Name: "New Name",
			},
		},
	}

	for _, tt := range testData {
		s.T().Run(tt.name, func(t *testing.T) {
			resp := api.ErrorResponse{}
			err := s.client.DoPostRequest(
				fmt.Sprintf("%s%s", mlflow.ExperimentsRoutePrefix, mlflow.ExperimentsUpdateRoute),
				tt.request,
				&resp,
			)
			assert.Nil(t, err)
			assert.Equal(s.T(), tt.error.Error(), resp.Error())
		})
	}
}
