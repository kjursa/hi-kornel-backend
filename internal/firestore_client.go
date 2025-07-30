package internal

import (
	"cloud.google.com/go/firestore"

	"context"

	"google.golang.org/api/option"
)

type FirestoreClient struct {
	firestoreClient *firestore.Client
	config          FirestoreConfig
}

func NewFirestoreClient(config FirestoreConfig) *FirestoreClient {
	return &FirestoreClient{
		config: config,
	}
}

func (f *FirestoreClient) Init() error {
	ctx := context.Background()
	op := option.WithCredentialsFile(f.config.Path)
	client, err := firestore.NewClientWithDatabase(ctx, f.config.ProjectId, f.config.DatabaseId, op)
	f.firestoreClient = client

	return err
}

func (f *FirestoreClient) Close() {
	if f.firestoreClient != nil {
		f.firestoreClient.Close()
	}
}

func (f *FirestoreClient) Get() *firestore.Client {
	return f.firestoreClient
}
