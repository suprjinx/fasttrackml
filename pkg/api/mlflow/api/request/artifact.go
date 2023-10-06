package request

// ListArtifactsRequest is a request object for `GET /mlflow/artifacts/list` endpoint.
type ListArtifactsRequest struct {
	Path    string `query:"path"`
	RunID   string `query:"run_id"`
	RunUUID string `query:"run_uuid"`
}

// GetRunID returns Run ID.
func (r ListArtifactsRequest) GetRunID() string {
	if r.RunID != "" {
		return r.RunID
	}
	return r.RunUUID
}

// GetArtifactRequest is for retrieving an individual artifact item.
type GetArtifactRequest struct {
	Path    string `query:"path"`
	RunID   string `query:"run_id"`
	RunUUID string `query:"run_uuid"`
}

// GetRunID returns RunID if available, otherwise RunUUID.
func (r GetArtifactRequest) GetRunID() string {
	if r.RunID != "" {
		return r.RunID
	}
	return r.RunUUID
}
