package resolvers

import (
	dataloaders "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/graphql/dataloaders"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

// This file is GraphQL transport dependency injection.

type Resolver struct {
	dataloader dataloaders.Dataloaders
	coreClient *coreadapters.CoreAdapter
}

func NewResolver(coreClient *coreadapters.CoreAdapter) *Resolver {
	return &Resolver{
		dataloader: dataloaders.NewDataloaders(coreClient),
		coreClient: coreClient,
	}
}
