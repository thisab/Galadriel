package datastore

import (
	"net/url"
	"time"

	management "github.com/HewlettPackard/galadriel/pkg/server/api/management"
	"github.com/labstack/echo/v4"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/spire/proto/spire/common"
)

// DataStore defines the data storage interface.
type DataStore interface {
	CreateOrganization(echo.Context, *management.Organization) (*management.Organization, error)
	CreateBridge(echo.Context, *management.FederationGroup, uint) (*management.FederationGroup, error)
	CreateMember(echo.Context, *management.SpireServer) (*management.SpireServer, error)
	CreateMembership(ctx echo.Context, membership *management.FederationGroupMembership, memberID uint, bridgeID uint) (*management.FederationGroupMembership, error)
	CreateTrustBundle(ctx echo.Context, trust *common.Bundle, memberID uint) (*common.Bundle, error)
	CreateRelationship(ctx echo.Context, newrelation *FederationRelationship, sourceID uint, targetID uint) error
	RetrieveOrganizationbyID(ctx echo.Context, orgID uint) (*management.Organization, error)
	RetrieveBridgebyID(ctx echo.Context, brID uint) (*management.FederationGroup, error)
	RetrieveAllBridgesbyOrgID(ctx echo.Context, orgID uint) (*[]management.FederationGroup, error)
	RetrieveAllMembershipsbyBridgeID(ctx echo.Context, bridgeID uint) (*[]management.FederationGroupMembership, error)
	RetrieveAllMembersbyBridgeID(ctx echo.Context, bridgeID uint) (*[]management.SpireServer, error)
	RetrieveAllBridgesbyMemberID(ctx echo.Context, memberID uint) (*[]management.FederationGroup, error)
	RetrieveAllMembershipsbyMemberID(ctx echo.Context, memberID uint) (*[]management.FederationGroupMembership, error)
	RetrieveAllRelationshipsbyMemberID(ctx echo.Context, memberID uint) (*[]FederationRelationship, error)
	RetrieveAllTrustBundlesbyMemberID(ctx echo.Context, memberID uint) (*[]common.Bundle, error)
	RetrieveMemberbyID(ctx echo.Context, memberID uint) (*management.SpireServer, error)
	RetrieveMembershipbyCreationDate(ctx echo.Context, date time.Time) (*management.FederationGroupMembership, error)
	RetrieveMembershipbyToken(ctx echo.Context, token string) (*management.FederationGroupMembership, error)
	RetrieveRelationshipbySourceandTargetID(ctx echo.Context, source uint, target uint) (*FederationRelationship, error)
	RetrieveTrustbundlebyMemberID(ctx echo.Context, memberID string) (*common.Bundle, error)
	UpdateBridge(echo.Context, management.FederationGroup) error
	UpdateOrganization(echo.Context, management.Organization) error
	UpdateMember(echo.Context, management.SpireServer) error
	UpdateMembership(echo.Context, management.FederationGroupMembership) error
	UpdateTrust(echo.Context, common.Bundle) error
	DeleteOrganizationbyID(ctx echo.Context, orgID uint) error
	DeleteBridgebyID(ctx echo.Context, bridgeID uint) error
	DeleteMemberbyID(ctx echo.Context, memberID uint) error
	DeleteAllMembershipsbyMemberID(ctx echo.Context, memberid uint) error
	DeleteAllMembershipsbyBridgeID(ctx echo.Context, bridgeid uint) error
	DeleteAllRelationshipsbyMemberID(ctx echo.Context, memberid uint) error
	DeleteAllTrustbundlesbyMemberID(ctx echo.Context, memberid uint) error
	DeleteMembershipbyToken(ctx echo.Context, name string) error
	DeleteRelationshipbySourceTargetID(ctx echo.Context, source uint, target uint) error
	DeleteTrustBundlebyMemberID(ctx echo.Context, memberID string) error
}
type BundleEndpointType string

type FederationRelationship struct {
	TrustDomain           spiffeid.TrustDomain
	BundleEndpointURL     *url.URL
	BundleEndpointProfile BundleEndpointType
	TrustDomainBundle     *common.Bundle

	// Fields only used for 'https_spiffe' bundle endpoint profile
	EndpointSPIFFEID spiffeid.ID
}
