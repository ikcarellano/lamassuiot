package postgres

import (
	"context"

	"github.com/lamassuiot/lamassuiot/pkg/v3/models"
	"github.com/lamassuiot/lamassuiot/pkg/v3/resources"
	"github.com/lamassuiot/lamassuiot/pkg/v3/storage"
	"gorm.io/gorm"
)

type PostgresEventsStore struct {
	db      *gorm.DB
	querier *postgresDBQuerier[models.AlertLatestEvent]
}

func NewEventsPostgresRepository(db *gorm.DB) (storage.EventRepository, error) {
	querier, err := CheckAndCreateTable(db, "events", "event_type", models.AlertLatestEvent{})
	if err != nil {
		return nil, err
	}

	return &PostgresEventsStore{
		db:      db,
		querier: (*postgresDBQuerier[models.AlertLatestEvent])(querier),
	}, nil
}

func (db *PostgresEventsStore) InsertUpdateEvent(ctx context.Context, ev *models.AlertLatestEvent) (*models.AlertLatestEvent, error) {
	return db.querier.Update(*ev, string(ev.EventType))
}

func (db *PostgresEventsStore) GetLatestEventByEventType(ctx context.Context, eventType models.EventType) (bool, *models.AlertLatestEvent, error) {
	return db.querier.SelectExists(string(eventType), nil)
}

func (db *PostgresEventsStore) GetLatestEvents(ctx context.Context) ([]*models.AlertLatestEvent, error) {
	evs := []*models.AlertLatestEvent{}
	_, err := db.querier.SelectAll(&resources.QueryParameters{}, []gormWhereParams{}, true, func(elem models.AlertLatestEvent) {
		derefElem := elem
		evs = append(evs, &derefElem)
	})

	if err != nil {
		return nil, err
	}

	return evs, nil
}
