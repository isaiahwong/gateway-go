package services

import (
	"fmt"

	accountsV1 "github.com/isaiahwong/gateway-go/protogen/accounts/v1"
)

// NewAccountsClient returns a new MailServiceClient
func NewAccountsClient(opts ...Option) (accountsV1.AccountsServiceClient, error) {
	conn, err := CreateClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("NewAccountsClient: %v", err)
	}
	client := accountsV1.NewAccountsServiceClient(conn)
	return client, nil
}
