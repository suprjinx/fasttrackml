//go:build integration

package run

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
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

type SearchTestSuite struct {
	suite.Suite
	client             *helpers.HttpClient
	runFixtures        *fixtures.RunFixtures
	tagFixtures        *fixtures.TagFixtures
	paramFixtures      *fixtures.ParamFixtures
	metricFixtures     *fixtures.MetricFixtures
	experimentFixtures *fixtures.ExperimentFixtures
}

func TestSearchTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}

func (s *SearchTestSuite) SetupTest() {
	s.client = helpers.NewMlflowApiClient(helpers.GetServiceUri())
	runFixtures, err := fixtures.NewRunFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.runFixtures = runFixtures
	tagFixtures, err := fixtures.NewTagFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.tagFixtures = tagFixtures
	paramFixtures, err := fixtures.NewParamFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.paramFixtures = paramFixtures
	metricFixtures, err := fixtures.NewMetricFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.metricFixtures = metricFixtures
	expFixtures, err := fixtures.NewExperimentFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.experimentFixtures = expFixtures
}

func (s *SearchTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.experimentFixtures.UnloadFixtures())
	}()

	// create test experiment.
	experiment, err := s.experimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name:           uuid.New().String(),
		LifecycleStage: models.LifecycleStageActive,
	})
	assert.Nil(s.T(), err)

	// create 3 different test runs and attach tags, metrics, params, etc.
	run1, err := s.runFixtures.CreateRun(context.Background(), &models.Run{
		ID:         "id1",
		Name:       "TestRun1",
		UserID:     "1",
		Status:     models.StatusRunning,
		SourceType: "JOB",
		StartTime: sql.NullInt64{
			Int64: 123456789,
			Valid: true,
		},
		EndTime: sql.NullInt64{
			Int64: 123456789,
			Valid: true,
		},
		ExperimentID:   *experiment.ID,
		ArtifactURI:    "artifact_uri1",
		LifecycleStage: models.LifecycleStageActive,
	})
	assert.Nil(s.T(), err)
	_, err = s.tagFixtures.CreateTag(context.Background(), &models.Tag{
		Key:   "mlflow.runName",
		Value: "TestRunTag1",
		RunID: run1.ID,
	})
	assert.Nil(s.T(), err)
	_, err = s.metricFixtures.CreateLatestMetric(context.Background(), &models.LatestMetric{
		Key:       "run1",
		Value:     1.1,
		Timestamp: 1234567890,
		Step:      1,
		IsNan:     false,
		RunID:     run1.ID,
		LastIter:  1,
	})
	assert.Nil(s.T(), err)
	_, err = s.paramFixtures.CreateParam(context.Background(), &models.Param{
		Key:   "param1",
		Value: "value1",
		RunID: run1.ID,
	})
	assert.Nil(s.T(), err)

	run2, err := s.runFixtures.CreateRun(context.Background(), &models.Run{
		ID:         "id2",
		Name:       "TestRun2",
		UserID:     "2",
		Status:     models.StatusScheduled,
		SourceType: "JOB",
		StartTime: sql.NullInt64{
			Int64: 111111111,
			Valid: true,
		},
		EndTime: sql.NullInt64{
			Int64: 222222222,
			Valid: true,
		},
		ExperimentID:   *experiment.ID,
		ArtifactURI:    "artifact_uri2",
		LifecycleStage: models.LifecycleStageDeleted,
	})
	assert.Nil(s.T(), err)
	_, err = s.tagFixtures.CreateTag(context.Background(), &models.Tag{
		Key:   "mlflow.runName",
		Value: "TestRunTag2",
		RunID: run2.ID,
	})
	assert.Nil(s.T(), err)
	_, err = s.metricFixtures.CreateLatestMetric(context.Background(), &models.LatestMetric{
		Key:       "run2",
		Value:     2.1,
		Timestamp: 1234567890,
		Step:      1,
		IsNan:     false,
		RunID:     run2.ID,
		LastIter:  1,
	})
	assert.Nil(s.T(), err)
	_, err = s.paramFixtures.CreateParam(context.Background(), &models.Param{
		Key:   "param2",
		Value: "value2",
		RunID: run2.ID,
	})
	assert.Nil(s.T(), err)

	run3, err := s.runFixtures.CreateRun(context.Background(), &models.Run{
		ID:         "id3",
		Name:       "TestRun3",
		UserID:     "3",
		Status:     models.StatusRunning,
		SourceType: "JOB",
		StartTime: sql.NullInt64{
			Int64: 333444444,
			Valid: true,
		},
		EndTime: sql.NullInt64{
			Int64: 444555555,
			Valid: true,
		},
		ExperimentID:   *experiment.ID,
		ArtifactURI:    "artifact_uri3",
		LifecycleStage: models.LifecycleStageActive,
	})
	assert.Nil(s.T(), err)
	_, err = s.tagFixtures.CreateTag(context.Background(), &models.Tag{
		Key:   "mlflow.runName",
		Value: "TestRunTag3",
		RunID: run3.ID,
	})
	assert.Nil(s.T(), err)
	_, err = s.metricFixtures.CreateLatestMetric(context.Background(), &models.LatestMetric{
		Key:       "run3",
		Value:     3.1,
		Timestamp: 1234567890,
		Step:      1,
		IsNan:     false,
		RunID:     run3.ID,
		LastIter:  1,
	})
	assert.Nil(s.T(), err)
	_, err = s.paramFixtures.CreateParam(context.Background(), &models.Param{
		Key:   "param3",
		Value: "value3",
		RunID: run3.ID,
	})
	assert.Nil(s.T(), err)

	tests := []struct {
		name     string
		error    *api.ErrorResponse
		request  *request.SearchRunsRequest
		response *response.SearchRunsResponse
	}{
		{
			name: "SearchWithViewTypeAllParameter3RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				ViewType:      request.ViewTypeAll,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run2.ID,
							Name:           "TestRunTag2",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "2",
							Status:         string(models.StatusScheduled),
							StartTime:      111111111,
							EndTime:        222222222,
							ArtifactURI:    "artifact_uri2",
							LifecycleStage: string(models.LifecycleStageDeleted),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag2",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param2",
									Value: "value2",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run2",
									Value:     2.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithViewTypeActiveOnlyParameter2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				ViewType:      request.ViewTypeActiveOnly,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithViewTypeDeletedOnlyParameter1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				ViewType:      request.ViewTypeDeletedOnly,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run2.ID,
							Name:           "TestRunTag2",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "2",
							Status:         string(models.StatusScheduled),
							StartTime:      111111111,
							EndTime:        222222222,
							ArtifactURI:    "artifact_uri2",
							LifecycleStage: string(models.LifecycleStageDeleted),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag2",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param2",
									Value: "value2",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run2",
									Value:     2.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationGrater1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time > 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationGraterOrEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time >= 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time != 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time = 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationLess1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time < 333444444`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationLessOrEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time <= 333444444`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationGrater1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time > 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationGraterOrEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time >= 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time != 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time = 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationLess1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time < 444555555`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationLessOrEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time <= 444555555`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunNameOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.run_name != "TestRunTag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunNameOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.run_name = "TestRunTag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunNameOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.run_name LIKE "TestRunTag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunNameOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.run_name ILIKE "testruntag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStatusOperationNotEqualNoRunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.status != "RUNNING"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{},
		},
		{
			name: "SearchWithAttributeStatusOperationEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.status = "RUNNING"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStatusOperationLike2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.status LIKE "RUNNING"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStatusOperationILike2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.status ILIKE "running"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeUserIDOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.user_id != 1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeUserIDOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.user_id = 3`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeUserIDOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.user_id LIKE "3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeUserIDOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.user_id ILIKE "3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeArtifactURIOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.artifact_uri != "artifact_uri1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeArtifactURIOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.artifact_uri = "artifact_uri3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeArtifactURIOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.artifact_uri LIKE "artifact_uri3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeArtifactURIOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.artifact_uri ILIKE "ArTiFaCt_UrI3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id != "%s"`, run1.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id = "%s"`, run3.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id LIKE "%s"`, run3.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id ILIKE "%s"`, strings.ToUpper(run3.ID)),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationIN1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id IN ('%s')`, run3.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationIN1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id NOT IN ('%s')`, run1.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationGrater1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 > 1.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationGraterOrEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 >= 1.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 != 1.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 = 3.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationLess0RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 < 3.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationLessOrEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 <= 3.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeParamsOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `params.param3 != "value1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeParamsOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `params.param3 = "value3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeParamsOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `params.param3 LIKE "value3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeParamsOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `params.param3 ILIKE "VaLuE3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeTagsOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `tags.mlflow.runName != "TestRunTag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag3",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param3",
									Value: "value3",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run3",
									Value:     3.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeTagsOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `tags.mlflow.runName = "TestRunTag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeTagsOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `tags.mlflow.runName LIKE "TestRunTag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeTagsOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `tags.mlflow.runName ILIKE "TeStRuNTaG1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
						Data: response.RunDataPartialResponse{
							Tags: []response.RunTagPartialResponse{
								{
									Key:   "mlflow.runName",
									Value: "TestRunTag1",
								},
							},
							Params: []response.RunParamPartialResponse{
								{
									Key:   "param1",
									Value: "value1",
								},
							},
							Metrics: []response.RunMetricPartialResponse{
								{
									Key:       "run1",
									Value:     1.1,
									Timestamp: 1234567890,
									Step:      1,
								},
							},
						},
					},
				},
				NextPageToken: "",
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			query, err := urlquery.Marshal(tt.request)
			assert.Nil(s.T(), err)
			resp := &response.SearchRunsResponse{}
			err = s.client.DoPostRequest(
				fmt.Sprintf("%s%s?%s", mlflow.RunsRoutePrefix, mlflow.RunsSearchRoute, query),
				tt.request,
				&resp,
			)
			assert.Nil(s.T(), err)
			assert.Equal(s.T(), len(tt.response.Runs), len(resp.Runs))
			assert.Equal(s.T(), len(tt.response.NextPageToken), len(resp.NextPageToken))

			mappedExpectedResult := make(map[string]*response.RunPartialResponse, len(tt.response.Runs))
			for _, run := range tt.response.Runs {
				mappedExpectedResult[run.Info.ID] = run
			}

			if tt.response.Runs != nil && resp.Runs != nil {
				for _, actualRun := range resp.Runs {
					expectedRun, ok := mappedExpectedResult[actualRun.Info.ID]
					assert.True(s.T(), ok)
					assert.NotEmpty(s.T(), actualRun.Info.ID)
					assert.Equal(s.T(), expectedRun.Info.Name, actualRun.Info.Name)
					assert.Equal(s.T(), expectedRun.Info.Name, actualRun.Info.Name)
					assert.Equal(s.T(), expectedRun.Info.UserID, actualRun.Info.UserID)
					assert.Equal(s.T(), expectedRun.Info.Status, actualRun.Info.Status)
					assert.Equal(s.T(), expectedRun.Info.EndTime, actualRun.Info.EndTime)
					assert.Equal(s.T(), expectedRun.Info.StartTime, actualRun.Info.StartTime)
					assert.Equal(s.T(), expectedRun.Info.ArtifactURI, actualRun.Info.ArtifactURI)
					assert.Equal(s.T(), expectedRun.Info.ExperimentID, actualRun.Info.ExperimentID)
					assert.Equal(s.T(), expectedRun.Info.LifecycleStage, actualRun.Info.LifecycleStage)
					assert.Equal(s.T(), expectedRun.Data.Tags, actualRun.Data.Tags)
					assert.Equal(s.T(), expectedRun.Data.Params, actualRun.Data.Params)
					assert.Equal(s.T(), expectedRun.Data.Metrics, actualRun.Data.Metrics)
				}
			}
		})
	}
}
