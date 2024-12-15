package utils

import "github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network"

func GetQueryParameter(page *int, perPage *int) []*network.QueryParameter {
	var queryParams []*network.QueryParameter
	if page != nil {
		queryParams = append(queryParams, &network.QueryParameter{Name: "page", Value: page})
	}
	if perPage != nil {
		queryParams = append(queryParams, &network.QueryParameter{Name: "perPage", Value: perPage})
	}
	return queryParams
}

