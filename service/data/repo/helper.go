package repo

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func CheckIteratorNextError(err error) error {
	if err != nil {
		if err == iterator.Done || grpc.Code(err) == codes.NotFound {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func ParseIterator(it *firestore.DocumentIterator, limit int) (*[]interface{}, bool, error) {
	var hasMoreResults bool
	var list []interface{}
	counter := 0

	for {
		var data interface{}
		doc, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				if grpc.Code(err) == codes.NotFound {
					return &list, hasMoreResults, nil
				}
				break
			} else {
				return &list, hasMoreResults, err
			}
		}

		counter++
		if counter > limit {
			hasMoreResults = true
			break
		}

		err = doc.DataTo(&data)
		if err != nil {
			return &list, hasMoreResults, err
		}

		list = append(list, &data)
	}

	return &list, hasMoreResults, nil
}
