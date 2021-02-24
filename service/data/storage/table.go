package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
	"gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/field"
	"gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/internal/query"
	externQuery "gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/query"
)

type tableStorage struct {
	db               *sql.DB
	ownerUIDField    string
	internalUIDField string
	externalUIDField string
	table            string
}

/*
NewTableStorage is intended for use with data that needs to be stored in a sql table
*/
func NewTableStorage(db *sql.DB, table string, opts ...Option) Storage {
	s := &tableStorage{
		db:               db,
		ownerUIDField:    "owner_uid",
		internalUIDField: "uid",
		externalUIDField: "external_uid",
		table:            table,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(s)
		}
	}
	return s
}

// TODO: move this to a helper of some kind
func (s *tableStorage) getUIDField(purp field.FieldPurpose) string {
	uidField := s.ownerUIDField
	switch purp {
	case field.PurposeReferenceOwner:
		uidField = s.ownerUIDField
		break
	case field.PurposeReferenceInternal:
		uidField = s.internalUIDField
		break
	case field.PurposeReferenceExternal:
		uidField = s.externalUIDField
		break
	}

	return uidField
}

// Get the first document in the collection that has a matching uid
func (s *tableStorage) First(ctx context.Context, purp field.FieldPurpose, uid string, dest interface{}, opts ...externQuery.Option) error {
	uidField := s.getUIDField(purp)
	q := sq.
		Select("*").
		From(s.table).
		Where(fmt.Sprintf("%s = ?", uidField), uid).
		Limit(1)
	c := query.NewMysqlQueryCustomize(&q, nil, nil)
	c.ApplyOptions(opts...)
	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	stmnt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}

	rows, err := stmnt.QueryContext(ctx, args...)
	if err != nil {
		return err
	}

	err = sqlscan.ScanOne(dest, rows)
	if err != nil {
		// If rows is empty
		if sqlscan.NotFound(err) {
			return nil
		}
		return err
	}

	return nil
}

// Update the first document in the collection that has a matching uid
func (s *tableStorage) UpdateFirst(ctx context.Context, purp field.FieldPurpose, uid string, data interface{}, opts ...externQuery.Option) error {
	uidField := s.getUIDField(purp)
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var mapData map[string]interface{}
	err = json.Unmarshal(rawData, &mapData)
	if err != nil {
		return err
	}

	q := sq.
		Update(s.table).
		Where(fmt.Sprintf("%s = ?", uidField), uid).
		SetMap(mapData).
		Limit(1)
	c := query.NewMysqlQueryCustomize(nil, &q, nil)
	c.ApplyOptions(opts...)
	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	stmnt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}

	_, err = stmnt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *tableStorage) CreateOne(ctx context.Context, data interface{}, opts ...externQuery.Option) error {
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var mapData map[string]interface{}
	err = json.Unmarshal(rawData, &mapData)
	if err != nil {
		return err
	}

	q := sq.Insert(s.table).SetMap(mapData)
	c := query.NewMysqlQueryCustomize(nil, nil, &q)
	c.ApplyOptions(opts...)
	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	stmnt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}

	_, err = stmnt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}

// Get many documents from a collection, returns true if there are more documents matching the query
func (s *tableStorage) Get(ctx context.Context, purp field.FieldPurpose, UIDs []string, dest *[]interface{}, opts ...externQuery.Option) (bool, error) {
	var hasMoreResults bool
	uidField := s.getUIDField(purp)
	q := sq.
		Select("*").
		From(s.table).
		Where(fmt.Sprintf("%s in ?", uidField), UIDs)
	c := query.NewMysqlQueryCustomize(&q, nil, nil)
	c.ApplyOptions(opts...)
	sql, args, err := q.ToSql()
	if err != nil {
		return hasMoreResults, err
	}

	stmnt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return hasMoreResults, err
	}

	rows, err := stmnt.QueryContext(ctx, args...)
	if err != nil {
		return hasMoreResults, err
	}

	err = sqlscan.ScanAll(&dest, rows)
	if err != nil {
		// If rows is empty
		if sqlscan.NotFound(err) {
			return hasMoreResults, nil
		}
		return hasMoreResults, err
	}

	return hasMoreResults, nil
}

// Delete from the collection the first document with a matching uid
func (s *tableStorage) DeleteFirst(ctx context.Context, purp field.FieldPurpose, uid string) error {
	uidField := s.getUIDField(purp)
	q := sq.
		Delete(s.table).
		Where(fmt.Sprintf("%s = ?", uidField), uid).
		Limit(1)
	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	stmnt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}

	_, err = stmnt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}

// Delete from the collection the first document with a matching uid
func (s *tableStorage) Delete(ctx context.Context, purp field.FieldPurpose, uids []string) error {
	uidField := s.getUIDField(purp)
	q := sq.
		Delete(s.table).
		Where(fmt.Sprintf("%s = ?", uidField), uids).
		Limit(1)
	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	stmnt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}

	_, err = stmnt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *tableStorage) setOwnerUIDFieldName(f string) {
	s.ownerUIDField = f
}

func (s *tableStorage) setInternalUIDFieldName(f string) {
	s.internalUIDField = f
}

func (s *tableStorage) setExternalUIDFieldName(f string) {
	s.externalUIDField = f
}
