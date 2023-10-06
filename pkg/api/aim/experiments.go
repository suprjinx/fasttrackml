package aim

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/repositories"
	"github.com/G-Research/fasttrackml/pkg/database"
)

func GetExperiments(c *fiber.Ctx) error {
	var experiments []struct {
		database.Experiment
		RunCount int
	}
	if tx := database.DB.Model(&database.Experiment{}).
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
		return fmt.Errorf("error fetching experiments: %w", tx.Error)
	}

	resp := make([]fiber.Map, len(experiments))
	for i, e := range experiments {
		resp[i] = fiber.Map{
			"id":            strconv.Itoa(int(*e.ID)),
			"name":          e.Name,
			"description":   nil,
			"archived":      e.LifecycleStage == database.LifecycleStageDeleted,
			"run_count":     e.RunCount,
			"creation_time": float64(e.CreationTime.Int64) / 1000,
		}
	}

	return c.JSON(resp)
}

func GetExperiment(c *fiber.Ctx) error {
	p := struct {
		ID string `params:"id"`
	}{}

	if err := c.ParamsParser(&p); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	id, err := strconv.ParseInt(p.ID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, fmt.Sprintf("unable to parse experiment id %q: %s", p.ID, err))
	}
	id32 := int32(id)

	if tx := database.DB.Select("ID").First(&database.Experiment{
		ID: &id32,
	}); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		return fmt.Errorf("unable to find experiment %q: %w", p.ID, tx.Error)
	}

	var exp struct {
		database.Experiment
		RunCount int
	}
	if tx := database.DB.Model(&database.Experiment{}).
		Select(
			"experiments.experiment_id",
			"experiments.name",
			"experiments.lifecycle_stage",
			"experiments.creation_time",
			"COUNT(runs.run_uuid) AS run_count",
		).
		Joins("LEFT JOIN runs USING(experiment_id)").
		Group("experiments.experiment_id").
		Where("experiments.experiment_id = ?", id).
		First(&exp); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		return fmt.Errorf("error fetching experiment %q: %w", p.ID, tx.Error)
	}

	return c.JSON(fiber.Map{
		"id":            id,
		"name":          exp.Name,
		"description":   nil,
		"archived":      exp.LifecycleStage == database.LifecycleStageDeleted,
		"run_count":     exp.RunCount,
		"creation_time": float64(exp.CreationTime.Int64) / 1000,
	})
}

func GetExperimentRuns(c *fiber.Ctx) error {
	q := struct {
		Limit  int    `query:"limit"`
		Offset string `query:"offset"`
	}{}

	if err := c.QueryParser(&q); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	p := struct {
		ID string `params:"id"`
	}{}

	if err := c.ParamsParser(&p); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	id, err := strconv.ParseInt(p.ID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, fmt.Sprintf("unable to parse experiment id %q: %s", p.ID, err))
	}
	id32 := int32(id)

	if tx := database.DB.Select("ID").First(&database.Experiment{
		ID: &id32,
	}); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		return fmt.Errorf("unable to find experiment %q: %w", p.ID, tx.Error)
	}

	tx := database.DB.
		Where("experiment_id = ?", id).
		Order("row_num DESC")

	if q.Limit > 0 {
		tx.Limit(q.Limit)
	}

	if q.Offset != "" {
		run := &database.Run{
			ID: q.Offset,
		}
		if tx := database.DB.Select("row_num").First(&run); tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
			return fmt.Errorf("unable to find search runs offset %q: %w", q.Offset, tx.Error)
		}

		tx.Where("row_num < ?", run.RowNum)
	}

	var sqlRuns []database.Run
	tx.Find(&sqlRuns)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		return fmt.Errorf("error fetching runs of experiment %q: %w", p.ID, tx.Error)
	}

	runs := make([]fiber.Map, len(sqlRuns))
	for i, r := range sqlRuns {
		runs[i] = fiber.Map{
			"run_id":        r.ID,
			"name":          r.Name,
			"creation_time": float64(r.StartTime.Int64) / 1000,
			"end_time":      float64(r.EndTime.Int64) / 1000,
			"archived":      r.LifecycleStage == database.LifecycleStageDeleted,
		}
	}

	return c.JSON(fiber.Map{
		"id":   p.ID,
		"runs": runs,
	})
}

func GetExperimentActivity(c *fiber.Ctx) error {
	tzOffset, err := strconv.Atoi(c.Get("x-timezone-offset", "0"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "x-timezone-offset header is not a valid integer")
	}

	p := struct {
		ID string `params:"id"`
	}{}

	if err := c.ParamsParser(&p); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	id, err := strconv.ParseInt(p.ID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, fmt.Sprintf("unable to parse experiment id %q: %s", p.ID, err))
	}
	id32 := int32(id)

	if tx := database.DB.Select("ID").First(&database.Experiment{
		ID: &id32,
	}); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		return fmt.Errorf("unable to find experiment %q: %w", p.ID, tx.Error)
	}

	var runs []database.Run
	if tx := database.DB.
		Select("StartTime", "LifecycleStage", "Status").
		Where("experiment_id = ?", id).
		Find(&runs); tx.Error != nil {
		return fmt.Errorf("error retrieving runs for experiment %q: %w", p.ID, tx.Error)
	}

	numArchivedRuns := 0
	numActiveRuns := 0
	activity := map[string]int{}
	for _, r := range runs {
		key := time.UnixMilli(r.StartTime.Int64).Add(time.Duration(-tzOffset) * time.Minute).Format("2006-01-02T15:00:00")
		activity[key] += 1
		switch {
		case r.LifecycleStage == database.LifecycleStageDeleted:
			numArchivedRuns += 1
		case r.Status == database.StatusRunning:
			numActiveRuns += 1
		}
	}

	return c.JSON(fiber.Map{
		"num_runs":          len(runs),
		"num_archived_runs": numArchivedRuns,
		"num_active_runs":   numActiveRuns,
		"activity_map":      activity,
	})
}

func DeleteExperiment(c *fiber.Ctx) error {
	params := struct {
		ID string `params:"id"`
	}{}

	if err := c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	id, err := strconv.ParseInt(params.ID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, fmt.Sprintf("unable to parse experiment id %q: %s", params.ID, err))
	}
	id32 := int32(id)

	// TODO this code should move to service with injected repository
	experimentRepo := repositories.NewExperimentRepository(database.DB)
	experiment := models.Experiment{ID: &id32}
	err = experimentRepo.Delete(c.Context(), &experiment)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError,
			fmt.Sprintf("unable to delete experiment %q: %s", params.ID, err))
	}

	return c.JSON(fiber.Map{
		"id":     params.ID,
		"status": "OK",
	})
}
