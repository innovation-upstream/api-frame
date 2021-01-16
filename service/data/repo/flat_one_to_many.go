package repo

import (
	"context"

	"cloud.google.com/go/firestore"
)

type flatOneToManyRepo struct {
	containerName string
	db            *firestore.Client
	ownerUIDField string
	dataUIDField  string
}

// TODO(@zach): Find a better home for these
type FlatOneToManyRepoOption func(*flatOneToManyRepo)

// WithOwnerUIDField sets the `ownerUIDField` used to lookup data by owner uid
func WithOwnerUIDField(f string) FlatOneToManyRepoOption {
	return func(r *flatOneToManyRepo) {
		r.ownerUIDField = f
	}
}

// WithDataUIDField sets the `dataUIDField` used to lookup data by data uid
func WithDataUIDField(f string) FlatOneToManyRepoOption {
	return func(r *flatOneToManyRepo) {
		r.dataUIDField = f
	}
}

/*
NewFlatOneToManyRepo is intended for use with data that may be accessed without knowledge of the
`ownerUIDField` value.
Data is stored in a single top level collection.
*/
func NewFlatOneToManyRepo(db *firestore.Client, collection string, opts ...FlatOneToManyRepoOption) Repo {
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

func (s *flatOneToManyRepo) ref(ownerUID string) *firestore.Query {
	q := s.db.Collection(s.containerName).Where(s.ownerUIDField, "==", ownerUID)

	return &q
}

func (s *flatOneToManyRepo) firstRef(ctx context.Context, ownerUID string) (*firestore.DocumentSnapshot, error) {
	docs := s.ref(ownerUID).Documents(ctx)
	doc, err := docs.Next()
	err = CheckIteratorNextError(err)
	if err != nil {
		return doc, err
	}

	return doc, nil
}

func (s *flatOneToManyRepo) dataRef(dataUID string) *firestore.Query {
	q := s.db.Collection(s.containerName).Where(s.dataUIDField, "==", dataUID)

	return &q
}

func (s *flatOneToManyRepo) firstDataRef(ctx context.Context, dataUID string) (*firestore.DocumentSnapshot, error) {
	docs := s.dataRef(dataUID).Documents(ctx)
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

// Get the first document in the collection where `ownerUID == ownerUIDField`
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

// Update the first document in the collection where `ownerUID == ownerUIDField`
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

// Update the first document in the collection where `dataUID == dataUIDField`
func (s *flatOneToManyRepo) UpdateFirstData(ctx context.Context, ownerUID string, data interface{}, opts ...firestore.SetOption) error {
	doc, err := s.firstDataRef(ctx, ownerUID)
	if err != nil {
		return err
	}

	_, err = doc.Ref.Set(ctx, data, opts...)
	if err != nil {
		return err
	}

	return nil
}

// Get the first document in the collection where `dataUID == dataUIDField`
func (s *flatOneToManyRepo) FirstData(ctx context.Context, dataUID string, dest interface{}) error {
	doc, err := s.firstDataRef(ctx, dataUID)
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

// Get many documents from a collection, returns true if there are more documents matching the query
func (s *flatOneToManyRepo) Get(ctx context.Context, ownerUID string, dest *[]interface{}, opts ...QueryOption) (bool, error) {
	q := s.ref(ownerUID)
	qm := &firestoreQueryMutate{
		query: q,
	}
	qm.ApplyMutations(opts...)
	docs := q.Documents(ctx)
	if docs == nil {
		return false, nil
	}

	list, hasMoreResults, err := ParseIterator(docs, qm.GetLimit())
	if err != nil {
		return hasMoreResults, err
	}

	*dest = *list

	return hasMoreResults, nil
}

// Get many documents from a collection, returns true if there are more documents matching the query
func (s *flatOneToManyRepo) GetData(ctx context.Context, dataUID string, dest *[]interface{}, opts ...QueryOption) (bool, error) {
	q := s.ref(dataUID)
	qm := &firestoreQueryMutate{
		query: q,
	}
	qm.ApplyMutations(opts...)
	docs := q.Documents(ctx)
	if docs == nil {
		return false, nil
	}

	list, hasMoreResults, err := ParseIterator(docs, qm.GetLimit())
	if err != nil {
		return hasMoreResults, err
	}

	*dest = *list

	return hasMoreResults, nil
}
