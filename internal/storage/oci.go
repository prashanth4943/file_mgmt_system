package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

type OCIStorage struct {
	Client     objectstorage.ObjectStorageClient
	Namespace  string
	BucketName string
	Region     string
}

// NewOCIStorage creates a new OCIStorage instance
func NewOCIStorage(provider common.ConfigurationProvider, bucketName string) (*OCIStorage, error) {
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create OCI client: %w", err)
	}

	namespaceResp, err := client.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}
	region, err := provider.Region()
	if err != nil {
		return nil, fmt.Errorf("failed to get region from provider: %w", err)
	}
	return &OCIStorage{
		Client:     client,
		Namespace:  *namespaceResp.Value,
		BucketName: bucketName,
		Region:     region,
	}, nil
}

func (o *OCIStorage) UploadFile(objectName string, content io.Reader, contentLength int64) (string, error) {
	_, err := o.Client.PutObject(context.Background(), objectstorage.PutObjectRequest{
		NamespaceName: &o.Namespace,
		BucketName:    &o.BucketName,
		ObjectName:    &objectName,
		ContentLength: &contentLength,
		PutObjectBody: io.NopCloser(content),
	})
	if err != nil {
		return "", err
	}
	ociReference := fmt.Sprintf("https://objectstorage.%s.oraclecloud.com/n/%s/b/%s/o/%s",
		o.Region, // Replace with the region name (e.g., "us-ashburn-1")
		o.Namespace,
		o.BucketName,
		objectName,
	)

	return ociReference, nil
}

func (o *OCIStorage) DownloadFile(objectName string) (io.Reader, error) {
	resp, err := o.Client.GetObject(context.Background(), objectstorage.GetObjectRequest{
		NamespaceName: &o.Namespace,
		BucketName:    &o.BucketName,
		ObjectName:    &objectName,
	})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

func (o *OCIStorage) DeleteFile(objectName string) error {
	_, err := o.Client.DeleteObject(context.Background(), objectstorage.DeleteObjectRequest{
		NamespaceName: &o.Namespace,
		BucketName:    &o.BucketName,
		ObjectName:    &objectName,
	})
	return err
}
