package repo

import (
	"context"

	"cloud.google.com/go/firestore"
)

//go:generate mockgen -destination=../mock/repo_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/iu-common-go/service/data/repo Repo

// Repo is a wrapper around `firestore.Client` that aims to simplify interaction with firestore.
type Repo interface {
	CreateOne(ctx context.Context, data interface{}) error
	First(ctx context.Context, ownerUID string, dest interface{}) error
	FirstData(ctx context.Context, dataUID string, dest interface{}) error
	UpdateFirst(ctx context.Context, ownerUID string, data interface{}, opts ...firestore.SetOption) error
	UpdateFirstData(ctx context.Context, dataUID string, data interface{}, opts ...firestore.SetOption) error
}
