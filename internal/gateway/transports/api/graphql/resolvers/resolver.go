package resolvers

import (
	dataloaders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/graphql/dataloaders"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

// This file is GraphQL transport dependency injection.

type Resolver struct {
	dataloader dataloaders.Dataloaders
	coreClient *coreadapters.CoreClient
}

func NewResolver(coreClient *coreadapters.CoreClient) *Resolver {
	return &Resolver{
		dataloader: dataloaders.NewDataloaders(coreClient),
		coreClient: coreClient,
	}
}
