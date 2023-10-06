//go:build integration

package run

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/tests/integration/golang/fixtures"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type DeleteRunTestSuite struct {
	suite.Suite
	client             *helpers.HttpClient
	runFixtures        *fixtures.RunFixtures
	experimentFixtures *fixtures.ExperimentFixtures
	runs               []*models.Run
}

func TestDeleteRunTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteRunTestSuite))
}

func (s *DeleteRunTestSuite) SetupTest() {
	s.client = helpers.NewAimApiClient(helpers.GetServiceUri())

	runFixtures, err := fixtures.NewRunFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.runFixtures = runFixtures

	expFixtures, err := fixtures.NewExperimentFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.experimentFixtures = expFixtures

	exp := &models.Experiment{
		Name:           uuid.New().String(),
		LifecycleStage: models.LifecycleStageActive,
	}
	_, err = s.experimentFixtures.CreateExperiment(context.Background(), exp)
	assert.Nil(s.T(), err)

	s.runs, err = s.runFixtures.CreateExampleRuns(context.Background(), exp, 10)
	assert.Nil(s.T(), err)
}

func (s *DeleteRunTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.runFixtures.UnloadFixtures())
	}()
	tests := []struct {
		name             string
		request          request.DeleteRunRequest
		expectedRunCount int
	}{
		{
			name:             "DeleteOneRun",
			request:          request.DeleteRunRequest{RunID: s.runs[4].ID},
			expectedRunCount: 9,
		},
		{
			name:             "RowNumbersAreRecalculated",
			request:          request.DeleteRunRequest{RunID: s.runs[1].ID},
			expectedRunCount: 8,
		},
		{
			name:             "DeleteFirstRun",
			request:          request.DeleteRunRequest{RunID: s.runs[0].ID},
			expectedRunCount: 7,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			originalMinRowNum, originalMaxRowNum, err := s.runFixtures.FindMinMaxRowNums(
				context.Background(),
				s.runs[0].ExperimentID,
			)
			assert.Nil(s.T(), err)

			var resp fiber.Map
			err = s.client.DoDeleteRequest(
				fmt.Sprintf("/runs/%s", tt.request.RunID),
				&resp,
			)
			assert.Nil(s.T(), err)

			runs, err := s.runFixtures.GetRuns(context.Background(), s.runs[0].ExperimentID)
			assert.Nil(s.T(), err)
			assert.Equal(s.T(), tt.expectedRunCount, len(runs))

			newMinRowNum, newMaxRowNum, err := s.runFixtures.FindMinMaxRowNums(
				context.Background(), s.runs[0].ExperimentID,
			)
			assert.Nil(s.T(), err)
			assert.Equal(s.T(), originalMinRowNum, newMinRowNum)
			assert.Greater(s.T(), originalMaxRowNum, newMaxRowNum)
		})
	}
}

func (s *DeleteRunTestSuite) Test_Error() {
	defer func() {
		assert.Nil(s.T(), s.runFixtures.UnloadFixtures())
	}()
	tests := []struct {
		name    string
		request request.DeleteRunRequest
	}{
		{
			name:    "DeleteWithUnknownID",
			request: request.DeleteRunRequest{RunID: "some-other-id"},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			originalMinRowNum, originalMaxRowNum, err := s.runFixtures.FindMinMaxRowNums(
				context.Background(), s.runs[0].ExperimentID,
			)
			assert.Nil(s.T(), err)

			var resp api.ErrorResponse
			err = s.client.DoDeleteRequest(
				fmt.Sprintf("/runs/%s", tt.request.RunID),
				&resp,
			)
			assert.Nil(s.T(), err)
			assert.Contains(s.T(), resp.Error(), "count of deleted runs does not match length of ids input")

			newMinRowNum, newMaxRowNum, err := s.runFixtures.FindMinMaxRowNums(
				context.Background(), s.runs[0].ExperimentID,
			)
			assert.Nil(s.T(), err)
			assert.Equal(s.T(), originalMinRowNum, newMinRowNum)
			assert.Equal(s.T(), originalMaxRowNum, newMaxRowNum)
		})
	}
}
