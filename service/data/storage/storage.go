package storage

import (
	"context"

	"gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/field"
	"gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/query"
)

//go:generate mockgen -destination=../mock/storage_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/iu-common-go/service/data/storage Storage

/*
Storage is a generic data persistance interface that aims to wrap interaction with data persistance
platforms so consumers do not need to consider storage platform/method when persisting data.
*/
type Storage interface {
	CreateOne(ctx context.Context, data interface{}, opts ...query.Option) error
	First(ctx context.Context, purpose field.FieldPurpose, uid string, dest interface{}, opts ...query.Option) error
	UpdateFirst(ctx context.Context, purpose field.FieldPurpose, uid string, data interface{}, opts ...query.Option) error
	Get(ctx context.Context, purpose field.FieldPurpose, UIDs []string, dest *[]interface{}, opts ...query.Option) (bool, error)
	DeleteFirst(ctx context.Context, purpose field.FieldPurpose, uid string) error
	Delete(ctx context.Context, purpose field.FieldPurpose, uids []string) error
	setOwnerUIDFieldName(f string)
	setInternalUIDFieldName(f string)
	setExternalUIDFieldName(f string)
}
