package repositories

import (
	"context"
	"errors"

	"github.com/rotisserie/eris"
	"gorm.io/gorm"

	"github.com/G-Research/fasttrackml/pkg/database"
	"github.com/G-Research/fasttrackml/pkg/models"
)

// ExperimentRepositoryProvider provides an interface to work with `experiment` entity.
type ExperimentRepositoryProvider interface {
	// List all the active Experiment
	List(ctx context.Context) (*[]models.Experiment, error)
	// Create creates new models.Experiment entity.
	Create(ctx context.Context, experiment *models.Experiment) error
	// GetByID returns experiment by its ID.
	GetByID(ctx context.Context, experimentID int32) (*models.Experiment, error)
	// GetByName returns experiment by its name.
	GetByName(ctx context.Context, name string) (*models.Experiment, error)
	// Update updates existing models.Experiment entity.
	Update(ctx context.Context, experiment *models.Experiment) error
}

// ExperimentRepository repository to work with `experiment` entity.
type ExperimentRepository struct {
	db *gorm.DB
}

// NewExperimentRepository creates repository to work with `experiment` entity.
func NewExperimentRepository(db *gorm.DB) *ExperimentRepository {
	return &ExperimentRepository{
		db: db,
	}
}

// List returns all the active Experiments
func (r ExperimentRepository) List(ctx context.Context) (*[]models.Experiment, error) {
	experiments := []models.Experiment{}
	if tx := r.db.Model(&models.Experiment{}).
		Select(
			"experiments.experiment_id",
			"experiments.name",
			"experiments.lifecycle_stage",
			"experiments.creation_time",
			"COUNT(runs.run_uuid) AS run_count",
		).
		Where("experiments.lifecycle_stage = ?", database.LifecycleStageActive).
		Joins("LEFT JOIN runs USING(experiment_id)").
		Group("experiments.experiment_id").
		Find(&experiments); tx.Error != nil {
		return nil, eris.Wrapf(tx.Error, "error fetching list of experiments")
	}
	return &experiments, nil
}

// Create creates new models.Experiment entity.
func (r ExperimentRepository) Create(ctx context.Context, experiment *models.Experiment) error {
	if err := r.db.Create(&experiment).Error; err != nil {
		return eris.Wrap(err, "error creating experiment entity")
	}
	if experiment.ArtifactLocation == "" {
		if err := database.DB.Model(
			&experiment,
		).Update(
			"ArtifactLocation", experiment.ArtifactLocation,
		).Error; err != nil {
			return eris.Wrapf(err, `error updating artifact_location: '%s'`, experiment.ArtifactLocation)
		}
	}
	return nil
}

// GetByID returns experiment by its ID.
func (r ExperimentRepository) GetByID(ctx context.Context, experimentID int32) (*models.Experiment, error) {
	var experiment models.Experiment
	if err := r.db.WithContext(ctx).Where(models.Experiment{ID: &experimentID}).Preload("Tags").First(&experiment).Error; err != nil {
		return nil, eris.Wrapf(err, "error getting experiment by id: %d", experimentID)
	}
	return &experiment, nil
}

// GetByName returns experiment by its name.
func (r ExperimentRepository) GetByName(ctx context.Context, name string) (*models.Experiment, error) {
	var experiment models.Experiment
	if err := r.db.WithContext(ctx).Preload(
		"Tags",
	).Where(
		models.Experiment{Name: name},
	).First(&experiment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, eris.Wrapf(err, "error getting experiment by id: %s", name)
	}
	return &experiment, nil
}

// Update updates existing models.Experiment entity.
func (r ExperimentRepository) Update(ctx context.Context, experiment *models.Experiment) error {
	if err := r.db.WithContext(ctx).Model(&experiment).Updates(experiment).Error; err != nil {
		return eris.Wrapf(err, "error updating experiment with id: %d", *experiment.ID)
	}
	return nil
}
