package database

import (
	"github.com/rotisserie/eris"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type experimentInfo struct {
	sourceID int64
	destID   int64
}

// Importer will handle transport of data from source to destination db.
type Importer struct {
	sourceDB        *gorm.DB
	destDB          *gorm.DB
	experimentInfos []experimentInfo
}

// NewImporter initializes an Importer.
func NewImporter(input, output DBProvider) *Importer {
	return &Importer{
		sourceDB:        input.GormDB(),
		destDB:          output.GormDB(),
		experimentInfos: []experimentInfo{},
	}
}

// Import will copy the contents of input db to output db.
func (s *Importer) Import() error {
	tables := []string{
		"experiment_tags",
		"runs",
		"tags",
		"params",
		"metrics",
		"latest_metrics",
		// "apps",
		// "dashboards",
	}
	// experiments needs special handling
	if err := s.importExperiments(); err != nil {
		return eris.Wrapf(err, "error importing table %s", "experiements")
	}
	// all other tables
	for _, table := range tables {
		if err := s.importTable(table); err != nil {
			return eris.Wrapf(err, "error importing table %s", table)
		}
	}
	return nil
}

// importExperiments will copy the contents of the experiments table from sourceDB to destDB,
// while recording the new ID.
func (s *Importer) importExperiments() error {
	// Start transaction in the destDB
	err := s.destDB.Transaction(func(destTX *gorm.DB) error {
		// Query data from the source database
		rows, err := s.sourceDB.Model(Experiment{}).Rows()
		if err != nil {
			return eris.Wrap(err, "error creating Rows instance from source")
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var scannedItem Experiment
			if err := s.sourceDB.ScanRows(rows, &scannedItem); err != nil {
				return eris.Wrap(err, "error creating Rows instance from source")
			}
			newItem := Experiment{
				Name:             scannedItem.Name,
				ArtifactLocation: scannedItem.ArtifactLocation,
				LifecycleStage:   scannedItem.LifecycleStage,
				CreationTime:     scannedItem.CreationTime,
				LastUpdateTime:   scannedItem.LastUpdateTime,
			}
			if err := destTX.
				Where(Experiment{Name: scannedItem.Name}).
				FirstOrCreate(&newItem).Error; err != nil {
				return eris.Wrap(err, "error creating destination row")
			}
			s.saveExperimentInfo(scannedItem, newItem)
			count++
		}
		log.Infof("Importing %s - found %v records", "experiments", count)
		return nil
	})
	if err != nil {
		return eris.Wrap(err, "error copying experiments table")
	}
	return nil
}

// importTablewill copy the contents of one table (model) from sourceDB
// while updating the experiment_id to destDB.
func (s *Importer) importTable(table string) error {
	// Start transaction in the destDB
	err := s.destDB.Transaction(func(destTX *gorm.DB) error {
		// Query data from the source database
		rows, err := s.sourceDB.Table(table).Rows()
		if err != nil {
			return eris.Wrap(err, "error creating Rows instance from source")
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var item map[string]any
			if err := s.sourceDB.ScanRows(rows, &item); err != nil {
				return eris.Wrap(err, "error scanning source row")
			}
			item, err = s.translateFields(item)
			if err != nil {
				return eris.Wrap(err, "error translating fields")
			}
			if err := destTX.
				Table(table).
				Clauses(clause.OnConflict{DoNothing: true}).
				Create(&item).Error; err != nil {
				return eris.Wrap(err, "error creating destination row")
			}
			count++
		}
		log.Infof("Importing %s - found %v records", table, count)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// saveExperimentInfo will relate the source and destination experiment for later id mapping.
func (s *Importer) saveExperimentInfo(source, dest Experiment) {
	s.experimentInfos = append(s.experimentInfos, experimentInfo{
		sourceID: int64(*source.ID),
		destID:   int64(*dest.ID),
	})
}

// translateFields will alter row before creation as needed (especially, replacing old experiment_id with new).
func (s *Importer) translateFields(item map[string]any) (map[string]any, error) {
	// boolean is numeric when coming from sqlite
	if isNaN, ok := item["is_nan"]; ok {
		switch v := isNaN.(type) {
		case bool:
			break
		default:
			item["is_nan"] = (v != 0.0)
		}
	}
	// items with experiment_id fk need to reference the new ID
	if expID, ok := item["experiment_id"]; ok {
		id, ok := expID.(int64)
		if !ok {
			return nil, eris.Errorf("unable to assert experiment_id as int64: %v", expID)
		}
		for _, expInfo := range s.experimentInfos {
			if expInfo.sourceID == id {
				item["experiment_id"] = expInfo.destID
			}
		}
	}
	return item, nil
}
