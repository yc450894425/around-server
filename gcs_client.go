package main

import (
	"context"
	"fmt"
	"io"
	"cloud.google.com/go/storage"
)

const (
	BUCKET_NAME = "around-zhangyuchao"
)

func saveToGCS(r io.Reader, id string) (string, error) {
	// create a client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	// write
	object := client.Bucket(BUCKET_NAME).Object(id)
	wc := object.NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	// make the file public
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}
	
	// get MediaLink from attrs
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	fmt.Printf("File saved to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil
}

func deleteFromGCS(id string) error {
	// create a client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	// delete
	object := client.Bucket(BUCKET_NAME).Object(id)
	if err := object.Delete(ctx); err != nil {
		return err
	}
	return nil
}