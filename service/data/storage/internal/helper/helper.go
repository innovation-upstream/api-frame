package helper

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

func ParseIterator(it *firestore.DocumentIterator, limit int, dest *[]interface{}) (bool, error) {
	var hasMoreResults bool
	counter := 0

	for {
		var data interface{}
		doc, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			} else if grpc.Code(err) == codes.NotFound {
				return hasMoreResults, nil
			} else {
				return hasMoreResults, err
			}
		}

		counter++
		if limit > 0 && counter > limit {
			hasMoreResults = true
			break
		}

		err = doc.DataTo(&data)
		if err != nil {
			return hasMoreResults, err
		}

		*dest = append(*dest, data)
	}

	return hasMoreResults, nil
}
