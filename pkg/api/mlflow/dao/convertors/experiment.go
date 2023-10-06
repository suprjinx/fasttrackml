package convertors

import (
	"database/sql"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/rotisserie/eris"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
)

// ConvertCreateExperimentToDBModel converts request.CreateExperimentRequest into actual models.Experiment model.
func ConvertCreateExperimentToDBModel(req *request.CreateExperimentRequest) (*models.Experiment, error) {
	ts := time.Now().UTC().UnixMilli()
	experiment := models.Experiment{
		Name:           req.Name,
		LifecycleStage: models.LifecycleStageActive,
		CreationTime: sql.NullInt64{
			Int64: ts,
			Valid: true,
		},
		LastUpdateTime: sql.NullInt64{
			Int64: ts,
			Valid: true,
		},
		Tags:             make([]models.ExperimentTag, len(req.Tags)),
		ArtifactLocation: strings.TrimRight(req.ArtifactLocation, "/"),
	}

	for n, tag := range req.Tags {
		experiment.Tags[n] = models.ExperimentTag{
			Key:   tag.Key,
			Value: tag.Value,
		}
	}

	if req.ArtifactLocation != "" {
		u, err := url.Parse(req.ArtifactLocation)
		if err != nil {
			return nil, eris.Wrap(err, "error parsing artifact location")
		}
		switch u.Scheme {
		case "s3":
			experiment.ArtifactLocation = strings.TrimRight(u.String(), "/")
		default:
			// TODO:DSuhinin - default case right now has to satisfy Python integration tests.
			p, err := filepath.Abs(u.Path)
			if err != nil {
				return nil, eris.Wrap(err, "error getting absolute path")
			}
			u.Path = p
			experiment.ArtifactLocation = u.String()
		}
	}

	return &experiment, nil
}

// ConvertUpdateExperimentToDBModel converts request.UpdateExperimentRequest into actual models.Experiment model.
func ConvertUpdateExperimentToDBModel(
	experiment *models.Experiment, req *request.UpdateExperimentRequest,
) *models.Experiment {
	experiment.Name = req.Name
	experiment.LastUpdateTime = sql.NullInt64{
		Int64: time.Now().UTC().UnixMilli(),
		Valid: true,
	}
	return experiment
}
