package controller

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/response"
)

// CreateRun handles `POST /runs/create` endpoint.
func (c Controller) CreateRun(ctx *fiber.Ctx) error {
	var req request.CreateRunRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}

	log.Debugf("create request: %#v", &req)
	run, err := c.runService.CreateRun(ctx.Context(), &req)
	if err != nil {
		return err
	}
	resp := response.NewCreateRunResponse(run)
	log.Debugf("create response: %#v", resp)

	return ctx.JSON(resp)
}

// UpdateRun handles `POST /runs/update` endpoint.
func (c Controller) UpdateRun(ctx *fiber.Ctx) error {
	var req request.UpdateRunRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}

	run, err := c.runService.UpdateRun(ctx.Context(), &req)
	if err != nil {
		return err
	}
	log.Debugf("updateRun request: %#v", req)
	resp := response.NewUpdateRunResponse(run)
	log.Debugf("updateRun response: %#v", resp)

	return ctx.JSON(resp)
}

// GetRun handles `GET /runs/get` endpoint.
func (c Controller) GetRun(ctx *fiber.Ctx) error {
	req := request.GetRunRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return api.NewBadRequestError(err.Error())
	}

	log.Debugf("getRun request: %#v", req)
	run, err := c.runService.GetRun(ctx.Context(), &req)
	if err != nil {
		return err
	}

	resp := response.NewGetRunResponse(run)
	log.Debugf("getRun response: %#v", resp)

	return ctx.JSON(resp)
}

// SearchRuns handles `POST /runs/search` endpoint.
func (c Controller) SearchRuns(ctx *fiber.Ctx) error {
	var req request.SearchRunsRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("searchRuns request: %#v", req)

	runs, limit, offset, err := c.runService.SearchRuns(ctx.Context(), &req)
	if err != nil {
		return err
	}

	resp, err := response.NewSearchRunsResponse(runs, limit, offset)
	if err != nil {
		return api.NewInternalError("Unable to build next_page_token: %s", err)
	}
	log.Debugf("searchRuns response: %#v", resp)

	return ctx.JSON(resp)
}

// DeleteRun handles `POST /runs/delete` endpoint.
func (c Controller) DeleteRun(ctx *fiber.Ctx) error {
	var req request.DeleteRunRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("deleteRun request: %#v", req)

	if err := c.runService.DeleteRun(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// RestoreRun handles `POST /runs/restore` endpoint.
func (c Controller) RestoreRun(ctx *fiber.Ctx) error {
	var req request.RestoreRunRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("restoreRun request: %#v", req)

	if err := c.runService.RestoreRun(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// LogMetric handles `POST /runs/log-metric` endpoint.
func (c Controller) LogMetric(ctx *fiber.Ctx) error {
	var req request.LogMetricRequest
	if err := ctx.BodyParser(&req); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return api.NewInvalidParameterValueError("Invalid value for parameter '%s' supplied. Hint: Value was of type '%s'. See the API docs for more information about request parameters.", err.Field, err.Value)
		}
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("logMetric request: %#v", req)

	if err := c.runService.LogMetric(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// LogParam handles `POST /runs/log-parameter` endpoint.
func (c Controller) LogParam(ctx *fiber.Ctx) error {
	var req request.LogParamRequest
	if err := ctx.BodyParser(&req); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return api.NewInvalidParameterValueError("Invalid value for parameter '%s' supplied. Hint: Value was of type '%s'. See the API docs for more information about request parameters.", err.Field, err.Value)
		}
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("logParam request: %#v", req)

	if err := c.runService.LogParam(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// SetRunTag handles `POST /runs/set-tag` endpoint.
func (c Controller) SetRunTag(ctx *fiber.Ctx) error {
	var req request.SetRunTagRequest
	if err := ctx.BodyParser(&req); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return api.NewInvalidParameterValueError("Invalid value for parameter '%s' supplied. Hint: Value was of type '%s'. See the API docs for more information about request parameters.", err.Field, err.Value)
		}
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("setRunTag request: %#v", req)

	if err := c.runService.SetRunTag(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// DeleteRunTag handles `POST /runs/delete-tag` endpoint.
func (c Controller) DeleteRunTag(ctx *fiber.Ctx) error {
	var req request.DeleteRunTagRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("deleteRunTag request: %#v", req)

	if err := c.runService.DeleteRunTag(ctx.Context(), &req); err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{})
}

// LogBatch handles `POST /runs/log-batch` endpoint.
func (c Controller) LogBatch(ctx *fiber.Ctx) error {
	var req request.LogBatchRequest
	if err := ctx.BodyParser(&req); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return api.NewInvalidParameterValueError("Invalid value for parameter '%s' supplied. Hint: Value was of type '%s'. See the API docs for more information about request parameters.", err.Field, err.Value)
		}
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("logBatch request: %#v", req)

	if err := c.runService.LogBatch(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}
