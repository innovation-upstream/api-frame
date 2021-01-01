package repo

import (
	"context"

	"cloud.google.com/go/firestore"
)

//go:generate mockgen -destination=../mock/repo_mock.go -package=mock gitlab.com/mobl-dev/common/repo Repo
type Repo interface {
	CreateOne(ctx context.Context, data interface{}) error
	First(ctx context.Context, ownerUID string, dest interface{}) error
	UpdateFirst(ctx context.Context, ownerUID string, data interface{}, opts ...firestore.SetOption) error
	SetOwnerUIDField(fieldName string)
}

type RepoOption func(Repo)
