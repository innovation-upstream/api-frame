package identity

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	firestoreadapter "github.com/innovation-upstream/casbin-firestore-adapter"
	"unknwon.dev/clog/v2"
)

func GetAccessControl(ctx context.Context, rawACModel string, db *firestore.Client) AccessControl {
	a := firestoreadapter.NewAdapter(db)
	acModel, err := model.NewModelFromString(rawACModel)
	if err != nil {
		clog.Fatal("%+v", err)
	}

	enforcer, err := casbin.NewEnforcer(acModel, a)
	if err != nil {
		clog.Fatal("%+v", err)
	}

	enforcer.EnableAutoSave(true)
	ac := NewAccessControl(enforcer)
	return ac
}
