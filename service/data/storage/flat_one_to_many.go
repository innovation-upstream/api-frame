package storage

import (
	"context"

	"cloud.google.com/go/firestore"
	"gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/field"
	"gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/internal/helper"
	"gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/internal/query"
	externQuery "gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/query"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type flatOneToManyCollectionStorage struct {
	collectionName   string
	db               *firestore.Client
	ownerUIDField    string
	internalUIDField string
	externalUIDField string
}

/*
NewFlatOneToManyCollectionStorage is intended for use with data that needs to be queriable by any
field
Data is stored in a single top level firestore collection.
*/
func NewFlatOneToManyCollectionStorage(db *firestore.Client, collection string, opts ...Option) Storage {
	s := &flatOneToManyCollectionStorage{
		collectionName:   collection,
		db:               db,
		ownerUIDField:    "owner_uid",
		internalUIDField: "uid",
		externalUIDField: "external_uid",
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(s)
		}
	}
	return s
}

func (s *flatOneToManyCollectionStorage) getUIDField(idType field.FieldPurpose) string {
	uidField := s.ownerUIDField
	switch idType {
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

func (s *flatOneToManyCollectionStorage) ref(refUID string, idType field.FieldPurpose) *firestore.Query {
	q := s.db.Collection(s.collectionName).Query
	if idType != field.PurposeReferenceNone {
		uidField := s.getUIDField(idType)
		q = q.Where(uidField, "==", refUID)
	}

	return &q
}

func (s *flatOneToManyCollectionStorage) refMany(idType field.FieldPurpose, uids []string) *firestore.Query {
	var q firestore.Query
	q = s.db.Collection(s.collectionName).Query
	if idType != field.PurposeReferenceNone {
		uidField := s.getUIDField(idType)
		q = q.Where(uidField, "in", uids)
	}

	return &q
}

func (s *flatOneToManyCollectionStorage) firstRef(ctx context.Context, idType field.FieldPurpose, refUID string) (*firestore.DocumentSnapshot, error) {
	docs := s.ref(refUID, idType).Limit(1).Documents(ctx)
	doc, err := docs.Next()
	err = helper.CheckIteratorNextError(err)
	if err != nil {
		return doc, err
	}

	return doc, nil
}

func (s *flatOneToManyCollectionStorage) getCollection() *firestore.CollectionRef {
	return s.db.Collection(s.collectionName)
}

// Get the first document in the collection that has a matching uid
func (s *flatOneToManyCollectionStorage) First(ctx context.Context, idType field.FieldPurpose, uid string, dest interface{}, opts ...externQuery.Option) error {
	q := s.ref(uid, idType)
	q.Limit(1)
	c := query.NewFirestoreQueryCustomize(q)
	c.ApplyOptions(opts...)
	docs := q.Documents(ctx)
	doc, err := docs.Next()
	err = helper.CheckIteratorNextError(err)
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

// Update the first document in the collection that has a matching uid
func (s *flatOneToManyCollectionStorage) UpdateFirst(ctx context.Context, purp field.FieldPurpose, uid string, data interface{}, opts ...externQuery.Option) error {
	q := s.ref(uid, purp)
	q.Limit(1)
	c := query.NewFirestoreQueryCustomize(q)
	c.ApplyOptions(opts...)
	docs := s.ref(uid, purp).Documents(ctx)
	doc, err := docs.Next()
	err = helper.CheckIteratorNextError(err)
	if err != nil {
		return err
	}

	var setOpts []firestore.SetOption
	optFields := c.GetQueryFields()
	if len(optFields) > 0 {
		var fieldPaths []firestore.FieldPath
		for _, field := range optFields {
			fieldPaths = append(fieldPaths, firestore.FieldPath([]string{field}))
		}
		setOpts = append(setOpts, firestore.Merge(fieldPaths...))
	}

	if doc == nil || !doc.Exists() {
		return s.CreateOne(ctx, data)
	} else {
		_, err = doc.Ref.Set(ctx, data, setOpts...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *flatOneToManyCollectionStorage) CreateOne(ctx context.Context, data interface{}) error {
	_, _, err := s.getCollection().Add(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// Get many documents from a collection, returns true if there are more documents matching the query
func (s *flatOneToManyCollectionStorage) Get(ctx context.Context, idType field.FieldPurpose, UIDs []string, dest *[]interface{}, opts ...externQuery.Option) (bool, error) {
	var q *firestore.Query
	q = s.refMany(idType, UIDs)
	c := query.NewFirestoreQueryCustomize(q)
	c.ApplyOptions(opts...)
	docs := q.Documents(ctx)
	if docs == nil {
		return false, nil
	}

	// subtract one from limit so we know if there are more results.
	// (by default, limit + 1 is retrieved)
	// TODO(zach): implement GetAfter so we never pull redundent records
	hasMoreResults, err := helper.ParseIterator(docs, c.GetLimit()-1, dest)
	if err != nil {
		return hasMoreResults, err
	}

	return hasMoreResults, nil
}

// Delete from the collection the first document with a matching uid
func (s *flatOneToManyCollectionStorage) DeleteFirst(ctx context.Context, purp field.FieldPurpose, uid string, opts ...externQuery.Option) error {
	q := s.refMany(purp, []string{uid})
	c := query.NewFirestoreQueryCustomize(q)
	c.ApplyOptions(opts...)
	q.Limit(1)
	docs := q.Documents(ctx)

	for {
		doc, err := docs.Next()
		if err != nil {
			if err == iterator.Done {
				if grpc.Code(err) == codes.NotFound {
					return nil
				}
				break
			} else {
				return err
			}
		}

		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete from the collection the first document with a matching uid
func (s *flatOneToManyCollectionStorage) Delete(ctx context.Context, purp field.FieldPurpose, uids []string, opts ...externQuery.Option) error {
	q := s.refMany(purp, uids)
	c := query.NewFirestoreQueryCustomize(q)
	c.ApplyOptions(opts...)
	docs := q.Documents(ctx)

	for {
		doc, err := docs.Next()
		if err != nil {
			if err == iterator.Done {
				if grpc.Code(err) == codes.NotFound {
					return nil
				}
				break
			} else {
				return err
			}
		}

		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *flatOneToManyCollectionStorage) setOwnerUIDFieldName(f string) {
	s.ownerUIDField = f
}

func (s *flatOneToManyCollectionStorage) setInternalUIDFieldName(f string) {
	s.internalUIDField = f
}

func (s *flatOneToManyCollectionStorage) setExternalUIDFieldName(f string) {
	s.externalUIDField = f
}
