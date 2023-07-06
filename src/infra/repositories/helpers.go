package repositories

import "github.com/castmetal/golang-token-api-authorizer/src/infra/storage/postgres"

func GetQuerieTransaction(qtx *postgres.Queries, pgRepositoryQuery *postgres.Queries) *postgres.Queries {
	if qtx != nil {
		return qtx
	}

	return pgRepositoryQuery
}
