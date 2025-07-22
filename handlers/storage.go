package handlers

import "miniapp/models"

type Repository interface {
	UpsertClient(client models.Client) error
	GetOwnClients(id int) ([]models.Client, error)
	SearchClient(passport string) (models.Client, error)
	GetUserData(username string) (string, error)
}
