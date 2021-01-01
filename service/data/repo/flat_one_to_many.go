package repo

import (
	"context"

	"cloud.google.com/go/firestore"
)

type flatOneToManyRepo struct {
	containerName string
	db            *firestore.Client
	ownerUIDField string
}

/**
For small/less frequently accessed data.
Stores all data in a single top level collection and queries using `Where(ownerUIDField, "==")`
*/
func NewFlatOneToManyRepo(db *firestore.Client, collection string, opts ...RepoOption) Repo {
	s := &flatOneToManyRepo{
		containerName: collection,
		db:            db,
		ownerUIDField: "owner_uid",
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(s)
		}
	}
	return s
}

func (s *flatOneToManyRepo) firstRef(ctx context.Context, ownerUID string) (*firestore.DocumentSnapshot, error) {
	docs := s.db.Collection(s.containerName).Where(s.ownerUIDField, "==", ownerUID).Documents(ctx)
	doc, err := docs.Next()
	err = CheckIteratorNextError(err)
	if err != nil {
		return doc, err
	}

	return doc, nil
}

func (s *flatOneToManyRepo) getContainer() *firestore.CollectionRef {
	return s.db.Collection(s.containerName)
}

// Get the first document in the collection with an ownerUID
func (s *flatOneToManyRepo) First(ctx context.Context, ownerUID string, dest interface{}) error {
	doc, err := s.firstRef(ctx, ownerUID)
	if err != nil {
		return err
	}
	if doc == nil || doc.Exists() == false {
		return nil
	}

	err = doc.DataTo(dest)
	if err != nil {
		return err
	}

	return nil
}

func (s *flatOneToManyRepo) SetOwnerUIDField(fieldName string) {
	s.ownerUIDField = fieldName
}

func (s *flatOneToManyRepo) UpdateFirst(ctx context.Context, ownerUID string, data interface{}, opts ...firestore.SetOption) error {
	doc, err := s.firstRef(ctx, ownerUID)
	if err != nil {
		return err
	}

	_, err = doc.Ref.Set(ctx, data, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (s *flatOneToManyRepo) CreateOne(ctx context.Context, data interface{}) error {
	_, _, err := s.getContainer().Add(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
