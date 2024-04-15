package credentialHandling

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type ServiceAccountSecretResponse struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

var (
	projectId             string = "methodologysandbox"
	serviceAccountKeyName string = "PAP_PULL_SA_KEY"
	secretNameString      string = fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectId, serviceAccountKeyName)
)

func AccessSecretVersion(name string) (ServiceAccountSecretResponse, error) {
	// create nil response
	var secretObj ServiceAccountSecretResponse

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return secretObj, err
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return secretObj, err
	}

	// marshal string secret into object
	secretString := string(result.Payload.Data)
	if err = json.Unmarshal([]byte(secretString), &secretObj); err != nil {
		return secretObj, err
	}

	return secretObj, nil
}

func GetServiceAccountSecret() (ServiceAccountSecretResponse, error) {
	sa_secret, err := AccessSecretVersion(secretNameString)
	if err != nil {
		log.Printf("failed to access secret version: %v", err)
		return sa_secret, err
	}

	log.Println("successfully accessed service account secret")
	return sa_secret, nil
}

func credsMain() error { // expect to return nil for testing
	_, err := AccessSecretVersion(secretNameString)
	if err != nil {
		log.Printf("failed to access secret version: %v", err)
		return err
	}
	return nil
}
