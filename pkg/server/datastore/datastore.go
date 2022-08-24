package datastore

import "github.com/labstack/echo/v4"

type Datastore struct {
}

// (GET /spireServers)
func (d *Datastore) GetSpireServers(ctx echo.Context) error {
	return nil
}

// (POST /spireServers)
func (d *Datastore) CreateSpireServer(ctx echo.Context) error {
	return nil
}

// (DELETE /spireServers/{spireServerId})
func (d *Datastore) DeleteSpireServer(ctx echo.Context, spireServerId int64) error {
	return nil
}

// (PUT /spireServers/{spireServerId})
func (d *Datastore) UpdateSpireServer(ctx echo.Context, spireServerId int64) error {
	return nil
}

// (PUT /trustBundles/{trustBundleId})
func (d *Datastore) UpdateTrustBundle(ctx echo.Context, trustBundleId int64) error {
	return nil
}

// (GET /federationRelationships)
func (d *Datastore) GetFederationRelationships(ctx echo.Context) error {
	return nil
}

// (POST /federationRelationships)
func (d *Datastore) CreateFederationRelationship(ctx echo.Context) error {
	return nil
}

// (GET /federationRelationships/{relationshipID})
func (d *Datastore) GetFederationRelationshipbyID(ctx echo.Context, relationshipID int64) error {
	return nil
}

// (PUT /federationRelationships/{relationshipID})
func (d *Datastore) UpdateFederationRelationship(ctx echo.Context, relationshipID int64) error {
	return nil
}
