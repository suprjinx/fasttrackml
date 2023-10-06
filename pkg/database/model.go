package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Status string

const (
	StatusRunning   Status = "RUNNING"
	StatusScheduled Status = "SCHEDULED"
	StatusFinished  Status = "FINISHED"
	StatusFailed    Status = "FAILED"
	StatusKilled    Status = "KILLED"
)

type LifecycleStage string

const (
	LifecycleStageActive  LifecycleStage = "active"
	LifecycleStageDeleted LifecycleStage = "deleted"
)

type Experiment struct {
	ID               *int32          `gorm:"column:experiment_id;not null;primaryKey"`
	Name             string          `gorm:"type:varchar(256);not null;unique"`
	ArtifactLocation string          `gorm:"type:varchar(256)"`
	LifecycleStage   LifecycleStage  `gorm:"type:varchar(32);check:lifecycle_stage IN ('active', 'deleted')"`
	CreationTime     sql.NullInt64   `gorm:"type:bigint"`
	LastUpdateTime   sql.NullInt64   `gorm:"type:bigint"`
	Tags             []ExperimentTag `gorm:"constraint:OnDelete:CASCADE"`
	Runs             []Run           `gorm:"constraint:OnDelete:CASCADE"`
}

type ExperimentTag struct {
	Key          string `gorm:"type:varchar(250);not null;primaryKey"`
	Value        string `gorm:"type:varchar(5000)"`
	ExperimentID int32  `gorm:"not null;primaryKey"`
}

type Run struct {
	ID             string         `gorm:"<-:create;column:run_uuid;type:varchar(32);not null;primaryKey"`
	Name           string         `gorm:"type:varchar(250)"`
	SourceType     string         `gorm:"<-:create;type:varchar(20);check:source_type IN ('NOTEBOOK', 'JOB', 'LOCAL', 'UNKNOWN', 'PROJECT')"`
	SourceName     string         `gorm:"<-:create;type:varchar(500)"`
	EntryPointName string         `gorm:"<-:create;type:varchar(50)"`
	UserID         string         `gorm:"<-:create;type:varchar(256)"`
	Status         Status         `gorm:"type:varchar(9);check:status IN ('SCHEDULED', 'FAILED', 'FINISHED', 'RUNNING', 'KILLED')"`
	StartTime      sql.NullInt64  `gorm:"<-:create;type:bigint"`
	EndTime        sql.NullInt64  `gorm:"type:bigint"`
	SourceVersion  string         `gorm:"<-:create;type:varchar(50)"`
	LifecycleStage LifecycleStage `gorm:"type:varchar(20);check:lifecycle_stage IN ('active', 'deleted')"`
	ArtifactURI    string         `gorm:"<-:create;type:varchar(200)"`
	ExperimentID   int32
	Experiment     Experiment
	DeletedTime    sql.NullInt64  `gorm:"type:bigint"`
	RowNum         RowNum         `gorm:"<-:create;index"`
	Params         []Param        `gorm:"constraint:OnDelete:CASCADE"`
	Tags           []Tag          `gorm:"constraint:OnDelete:CASCADE"`
	Metrics        []Metric       `gorm:"constraint:OnDelete:CASCADE"`
	LatestMetrics  []LatestMetric `gorm:"constraint:OnDelete:CASCADE"`
}

type RowNum int64

func (rn *RowNum) Scan(v interface{}) error {
	nullInt := sql.NullInt64{}
	if err := nullInt.Scan(v); err != nil {
		return err
	}
	*rn = RowNum(nullInt.Int64)
	return nil
}

func (rn RowNum) GormDataType() string {
	return "bigint"
}

func (rn RowNum) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if rn == 0 {
		return clause.Expr{
			SQL: "(SELECT COALESCE(MAX(row_num), -1) FROM runs) + 1",
		}
	}
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{int64(rn)},
	}
}

type Param struct {
	Key   string `gorm:"type:varchar(250);not null;primaryKey"`
	Value string `gorm:"type:varchar(500);not null"`
	RunID string `gorm:"column:run_uuid;not null;primaryKey;index"`
}

type Tag struct {
	Key   string `gorm:"type:varchar(250);not null;primaryKey"`
	Value string `gorm:"type:varchar(5000)"`
	RunID string `gorm:"column:run_uuid;not null;primaryKey;index"`
}

type Metric struct {
	Key       string  `gorm:"type:varchar(250);not null;primaryKey"`
	Value     float64 `gorm:"type:double precision;not null;primaryKey"`
	Timestamp int64   `gorm:"not null;primaryKey"`
	RunID     string  `gorm:"column:run_uuid;not null;primaryKey;index"`
	Step      int64   `gorm:"default:0;not null;primaryKey"`
	IsNan     bool    `gorm:"default:false;not null;primaryKey"`
	Iter      int64   `gorm:"index"`
}

type LatestMetric struct {
	Key       string  `gorm:"type:varchar(250);not null;primaryKey"`
	Value     float64 `gorm:"type:double precision;not null"`
	Timestamp int64
	Step      int64  `gorm:"not null"`
	IsNan     bool   `gorm:"not null"`
	RunID     string `gorm:"column:run_uuid;not null;primaryKey;index"`
	LastIter  int64
}

type AlembicVersion struct {
	Version string `gorm:"column:version_num;type:varchar(32);not null;primaryKey"`
}

func (AlembicVersion) TableName() string {
	return "alembic_version"
}

type SchemaVersion struct {
	Version string `gorm:"not null;primaryKey"`
}

func (SchemaVersion) TableName() string {
	return "schema_version"
}

type Base struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsArchived bool      `json:"-"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}

type Dashboard struct {
	Base
	Name        string     `json:"name"`
	Description string     `json:"description"`
	AppID       *uuid.UUID `gorm:"type:uuid" json:"app_id"`
	App         App        `json:"-"`
}

func (d Dashboard) MarshalJSON() ([]byte, error) {
	type localDashboard Dashboard
	type jsonDashboard struct {
		localDashboard
		AppType *string `json:"app_type"`
	}
	jd := jsonDashboard{
		localDashboard: localDashboard(d),
	}
	if d.App.IsArchived {
		jd.AppID = nil
	} else {
		jd.AppType = &d.App.Type
	}
	return json.Marshal(jd)
}

type App struct {
	Base
	Type  string   `gorm:"not null" json:"type"`
	State AppState `json:"state"`
}

type AppState map[string]any

func (s AppState) Value() (driver.Value, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(v), nil
}

func (s *AppState) Scan(v interface{}) error {
	var nullS sql.NullString
	if err := nullS.Scan(v); err != nil {
		return err
	}
	if nullS.Valid {
		return json.Unmarshal([]byte(nullS.String), s)
	}
	s = nil
	return nil
}

func (s AppState) GormDataType() string {
	return "text"
}

func NewUUID() string {
	var r [32]byte
	u := uuid.New()
	hex.Encode(r[:], u[:])
	return string(r[:])
}
