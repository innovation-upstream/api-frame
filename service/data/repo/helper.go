package repo

import (
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
