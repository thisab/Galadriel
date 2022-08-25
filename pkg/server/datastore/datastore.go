package datastore

import (
	"context"
	"net/url"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/spire/proto/spire/common"
)

// DataStore defines the data storage interface.
type DataStore interface {
	CreateOrganization(ctx context.Context, org *Organization) (*Organization, error)
	CreateBridge(ctx context.Context, br *Bridge, orgID uint) (*Bridge, error)
	CreateMember(ctx context.Context, mem *Member) (*Member, error)
	CreateMembership(ctx context.Context, memb *Membership, memberID uint, bridgeID uint) (*Membership, error)
	CreateTrustBundle(ctx context.Context, trust *TrustBundle, memberID uint) (*TrustBundle, error)
	CreateRelationship(ctx context.Context, newrelation *FederationRelationship, sourceID uint, targetID uint) error
	RetrieveOrganizationbyID(ctx context.Context, orgID uint) (*Organization, error)
	RetrieveBridgebyID(ctx context.Context, brID uint) (*Bridge, error)
	RetrieveAllBridgesbyOrgID(ctx context.Context, orgID uint) (*[]Bridge, error)
	RetrieveAllMembershipsbyBridgeID(ctx context.Context, bridgeID uint) (*[]Membership, error)
	RetrieveAllMembersbyBridgeID(ctx context.Context, bridgeID uint) (mem *[]Member, err error)
	RetrieveAllBridgesbyMemberID(ctx context.Context, memberID uint) (mem *[]Bridge, err error)
	RetrieveAllMembershipsbyMemberID(ctx context.Context, memberID uint) (*[]Membership, error)
	RetrieveAllRelationshipsbyMemberID(ctx context.Context, memberID uint) (*[]FederationRelationship, error)
	RetrieveAllTrustBundlesbyMemberID(ctx context.Context, memberID uint) (*[]TrustBundle, error)
	RetrieveMemberbyID(ctx context.Context, memberID uint) (*Member, error)
	RetrieveMembershipbyCreationDate(ctx context.Context, date time.Time) (*Membership, error)
	RetrieveMembershipbyToken(ctx context.Context, token string) (*Membership, error)
	RetrieveRelationshipbySourceandTargetID(ctx context.Context, source uint, target uint) (*FederationRelationship, error)
	RetrieveTrustbundlebyMemberID(ctx context.Context, memberID string) (*TrustBundle, error)
	UpdateBridge(ctx context.Context, br Bridge) error
	UpdateOrganization(ctx context.Context, org Organization) error
	UpdateMember(ctx context.Context, member Member) error
	UpdateMembership(ctx context.Context, membership Membership) error
	UpdateTrust(ctx context.Context, trust TrustBundle) error
	DeleteOrganizationbyID(ctx context.Context, orgID uint) error
	DeleteBridgebyID(ctx context.Context, bridgeID uint) error
	DeleteMemberbyID(ctx context.Context, memberID uint) error
	DeleteAllMembershipsbyMemberID(ctx context.Context, memberid uint) error
	DeleteAllMembershipsbyBridgeID(ctx context.Context, bridgeid uint) error
	DeleteAllRelationshipsbyMemberID(ctx context.Context, memberid uint) error
	DeleteAllTrustbundlesbyMemberID(ctx context.Context, memberid uint) error
	DeleteMembershipbyToken(ctx context.Context, name string) error
	DeleteRelationshipbySourceTargetID(ctx context.Context, source uint, target uint) error
	DeleteTrustBundlebyMemberID(ctx context.Context, memberID string) error
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
