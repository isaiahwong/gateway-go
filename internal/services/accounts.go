package services

import (
	"fmt"
	"github.com/isaiahwong/gateway-go/api/gen/accounts/v1"
)

// NewAccountsClient returns a new MailServiceClient
func NewAccountsClient(opts ...Option) (accounts.AccountsServiceClient, error) {
	conn, err := CreateClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("NewAccountsClient: %v", err)
	}
	client := accounts.NewAccountsServiceClient(conn)
	return client, nil
}
