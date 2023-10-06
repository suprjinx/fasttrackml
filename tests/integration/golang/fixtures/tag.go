package fixtures

import (
	"context"

	"github.com/rotisserie/eris"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
)

// TagFixtures represents data fixtures object.
type TagFixtures struct {
	baseFixtures
}

// NewTagFixtures creates new instance of TagFixtures.
func NewTagFixtures(databaseDSN string) (*TagFixtures, error) {
	db, err := CreateDB(databaseDSN)
	if err != nil {
		return nil, err
	}
	return &TagFixtures{
		baseFixtures: baseFixtures{db: db.GormDB()},
	}, nil
}

// CreateTag creates new test Tag.
func (f TagFixtures) CreateTag(ctx context.Context, tag *models.Tag) (*models.Tag, error) {
	if err := f.baseFixtures.db.WithContext(ctx).Create(tag).Error; err != nil {
		return nil, eris.Wrap(err, "error creating test tag")
	}
	return tag, nil
}

// GetByRunID returns tag list by requested Run ID.
func (f TagFixtures) GetByRunID(ctx context.Context, runID string) ([]models.Tag, error) {
	var tags []models.Tag
	if err := f.db.WithContext(ctx).Where(
		models.Tag{RunID: runID},
	).Find(&tags).Error; err != nil {
		return nil, eris.Wrapf(err, "error getting tags by run id: %s", runID)
	}
	return tags, nil
}
