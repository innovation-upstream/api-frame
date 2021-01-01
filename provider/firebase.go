package provider

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type FirebaseClient interface {
	Auth() *auth.Client
}

type firebaseClient struct {
	AuthClient *auth.Client
}

func NewFirebaseClient(ctx context.Context) FirebaseClient {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		panic(err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	return &firebaseClient{
		AuthClient: client,
	}
}

func (e *firebaseClient) Auth() *auth.Client {
	return e.Auth()
}
