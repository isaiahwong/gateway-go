package services

import (
	"fmt"

	"gitlab.com/eco_system/gateway/api/gen/go/accounts/v1"
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
