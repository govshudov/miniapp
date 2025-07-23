package postgresql

import (
	"database/sql"
	"fmt"
	"miniapp/models"
)

type Repository struct {
	client *sql.DB
}

func NewPostgreSQLRepository(client *sql.DB) *Repository {
	return &Repository{
		client: client,
	}
}
func (repository *Repository) UpsertClient(client models.Client) error {
	var usd, eur float64
	q := `select usd,eur from clients where passport = $1 limit 1`
	err := repository.client.QueryRow(q, client.Passport).Scan(&usd, &eur)
	if err != nil {
		fmt.Println("select clienr error: ", err)
	}
	tx, _ := repository.client.Begin()
	defer tx.Rollback()
	if usd == 0 && eur == 0 {
		q = `insert into clients(user_id,name,passport,usd,eur,currency) values($1,$2,$3,$4,$5,$6)`
		_, err = tx.Exec(q, client.UserId, client.Name, client.Passport, client.USD, client.EUR, client.Currency)
		if err != nil {
			return err
		}
		q = `insert into client_histories(user_id,passport,usd,eur,currency,is_plus) values($1,$2,$3,$4,$5,$6)`
		_, err = tx.Exec(q, client.UserId, client.Passport, client.USD, client.EUR, client.Currency, true)
		if err != nil {
			return err
		}

	} else {
		var isPlus bool
		if (usd > 0 && client.USD < usd) || (eur > 0 && client.EUR < eur) {
			isPlus = false
		} else {
			isPlus = true
		}
		q = `update clients set(name,usd,eur,currency)=($1,$2,$3,$4) where passport = $5`
		_, err = tx.Exec(q, client.Name, client.USD, client.EUR, client.Currency, client.Passport)
		if err != nil {
			return err
		}
		q = `insert into client_histories(user_id,passport,usd,eur,currency,is_plus) values($1,$2,$3,$4,$5,$6)`
		_, err = tx.Exec(q, client.UserId, client.Passport, client.USD, client.EUR, client.Currency, isPlus)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("commit", err)
	}
	return nil
}
func (repository *Repository) GetOwnClients(id int) ([]models.Client, error) {
	var clients []models.Client
	q := `select user_id,name,passport,usd,eur,currency,users.username,clients.created_at from clients 
			left join users on users.id = clients.user_id
			where user_id = $1 order by created_at desc`
	rows, err := repository.client.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var client models.Client
		err = rows.Scan(&client.UserId, &client.Name, &client.Passport, &client.USD, &client.EUR, &client.Currency, &client.ReceivedName, &client.CreatedDate)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}
func (repository *Repository) SearchClient(passport string) (models.Client, error) {
	var client models.Client
	q := `select user_id,name,passport,usd,eur,currency,users.username,clients.created_at from clients
			left join users on users.id = clients.user_id
			where passport = $1 limit 1`
	err := repository.client.QueryRow(q, passport).Scan(&client.UserId, &client.Name, &client.Passport, &client.USD, &client.EUR, &client.Currency, &client.ReceivedName, &client.CreatedDate))
	if err != nil {
		return client, err
	}

	return client, nil
}
func (repository *Repository) GetUserData(username string) (int, string, error) {
	var (
		userId   int
		password string
	)

	q := `select id,password from users where username=$1 limit 1`
	err := repository.client.QueryRow(q, username).Scan(&userId, &password)
	if err != nil {
		return 0, "", err
	}

	return userId, password, nil
}
