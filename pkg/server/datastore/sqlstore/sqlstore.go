package sqlstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	// MySQL database type
	MySQL = "mysql"
	// PostgreSQL database type
	PostgreSQL = "postgres"
	// SQLite database type
	SQLite = "sqlite3"
)

type Plugin struct {
	db *gorm.DB
}

func (ds *Plugin) OpenDB(connectionString, dbtype string) (err error) {
	var dialectvar dialect

	switch dbtype {
	case SQLite:
		dialectvar = sqliteDB{}
	case PostgreSQL:
		dialectvar = postgresDB{}
	case MySQL:
		dialectvar = mysqlDB{}
	default:
		return fmt.Errorf("unsupported database_type: %s" + dbtype)
	}
	ds.db, err = dialectvar.connect(connectionString)
	if err != nil {
		return fmt.Errorf("error connecting to: %s", connectionString)
	}
	return nil
}

func (ds *Plugin) CreateAllTablesinDB() error {
	if err := ds.createOrganizationTableInDB(); err != nil {
		return err
	}
	if err := ds.createBridgeTableInDB(); err != nil {
		return err
	}
	if err := ds.createMemberTableInDB(); err != nil {
		return err
	}
	if err := ds.createMembershipTableInDB(); err != nil {
		return err
	}
	if err := ds.createRelationshipTableInDB(); err != nil {
		return err
	}
	if err := ds.createTrustbundleTableInDB(); err != nil {
		return err
	}
	return nil
}

// Creates the Table for the Organization Model
// Returns Error if AutoMigrate fails
func (ds *Plugin) createOrganizationTableInDB() error {
	err := ds.db.AutoMigrate(&Organization{})
	if err != nil {
		return fmt.Errorf("sqlstorage error: automigrate: %v", err)
	}
	return nil
}

// Creates the Table for the Bridge Model
// Returns Error if AutoMigrate fails
func (ds *Plugin) createBridgeTableInDB() error {
	err := ds.db.AutoMigrate(&Bridge{})
	if err != nil {
		return fmt.Errorf("sqlstorage error: automigrate: %v", err)
	}
	return nil
}

// Creates the Table for the Member Model
// Returns Error if AutoMigrate fails
func (ds *Plugin) createMemberTableInDB() error {
	err := ds.db.AutoMigrate(&Member{})
	if err != nil {
		return fmt.Errorf("sqlstorage error: automigrate: %v", err)
	}
	return nil
}

// Creates the Table for the Membership Model
// Returns Error if AutoMigrate fails
func (ds *Plugin) createMembershipTableInDB() error {
	err := ds.db.AutoMigrate(&Membership{})
	if err != nil {
		return fmt.Errorf("sqlstore error: automigrate: %v", err)
	}
	return nil
}

// Creates the Table for the Relationship Model
// Returns Error if AutoMigrate fails
func (ds *Plugin) createRelationshipTableInDB() error {

	err := ds.db.AutoMigrate(&Relationship{})
	if err != nil {
		return fmt.Errorf("sqlstore error: automigrate: %v", err)
	}
	return nil
}

// Creates the Table for the Trustbundle Model
func (ds *Plugin) createTrustbundleTableInDB() error {

	err := ds.db.AutoMigrate(&TrustBundle{})
	if err != nil {
		return fmt.Errorf("sqlstore error: automigrate: %v", err)
	}
	return nil
}

func (ds *Plugin) CreateOrganization(ctx context.Context, org *Organization) (*Organization, error) {
	return ds.createOrganization(org)
}

// Insert a new Organization into the DB.
// Ignores and returns nil if entry already exists. Returns an error if creation fails
func (ds *Plugin) createOrganization(org *Organization) (*Organization, error) {

	err := ds.db.Where(&org).FirstOrCreate(org).Error
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return org, nil
}

// Creates a new Bridge or ATB  from an Organization Name
// Ignores and returns nil if it already exists. Returns an error if creation fails
func (ds *Plugin) CreateBridge(br *Bridge, orgID uint) (*Bridge, error) {

	org, err := ds.RetrieveOrganizationbyID(orgID)
	if err != nil {
		return nil, err
	}
	br.OrganizationID = org.ID // Fill in the OrgID for the bridge
	err = ds.db.Where(&br).FirstOrCreate(br).Error
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return br, nil
}

// Creates a new Member
// Ignores and returns nil if entry already exists. Returns an error if creation fails
func (ds *Plugin) CreateMember(mem *Member) (*Member, error) {
	err := ds.db.Where(mem).FirstOrCreate(mem).Error
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return mem, nil
}

// Creates a new Membership from a Member
// Ignores and returns nil if entry already exists. Returns an error if creation fails
func (ds *Plugin) CreateMembership(memb *Membership, memberID uint, bridgeID uint) (*Membership, error) {
	// Check if Member exists in DB
	mem, err := ds.RetrieveMemberbyID(memberID)
	if err != nil {
		return nil, err
	}
	var br *Bridge
	// Check if Bridge exists in DB
	br, err = ds.RetrieveBridgebyID(bridgeID)
	if err != nil {
		return nil, err
	}
	memb.MemberID = mem.ID // Fill in the BridgeID for the bridge
	memb.BridgeID = br.ID  // Fill in the BridgeID for the bridge
	err = ds.db.Where(memb).FirstOrCreate(memb).Error
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return memb, nil
}

// Create a new Trustbundle from a Member
// Ignores and returns nil if entry already exists. Returns an error if creation fails
func (ds *Plugin) CreateTrustBundle(trust *TrustBundle, memberID uint) (*TrustBundle, error) {

	mem, err := ds.RetrieveMemberbyID(memberID)
	if err != nil {
		return nil, err
	}
	trust.MemberID = mem.ID // Fill in the BridgeID for the bridge
	err = ds.db.Where(trust).FirstOrCreate(trust).Error
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return trust, nil
}

// Adds a new Relation between Two SPIRE Servers in DB using description as reference for the IDs
// Ignores and returns nil if entry already exists. Returns an error if creation fails
func (ds *Plugin) CreateRelationship(newrelation *Relationship, sourceID uint, targetID uint) error {

	// Check if Member exists in DB
	sourcemember, err := ds.RetrieveMemberbyID(sourceID)
	if err != nil {
		return err
	}
	var targetmember *Member
	// Check if Member exists in DB
	targetmember, err = ds.RetrieveMemberbyID(targetID)
	if err != nil {
		return err
	}
	newrelation.MemberID = sourcemember.ID       // Fill in the Source MemberID (Foreign Key)
	newrelation.TargetMemberID = targetmember.ID // Fill in the Target MemberID
	err = ds.db.Where(&newrelation).FirstOrCreate(&newrelation).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// retrieves an Organization from the Database by Name. returns an error if query fails
func (ds *Plugin) RetrieveOrganizationbyID(orgID uint) (*Organization, error) {
	var org *Organization = &Organization{}
	err := ds.db.Where("id = ?", orgID).First(org).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // If does not, throw an error
		return nil, fmt.Errorf("sqlstore error: organization %v does not exist in db", orgID)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return org, nil
}

// RetrieveBridgebyDescription retrieves a Bridge from the Database by description. returns an error if query fails
func (ds *Plugin) RetrieveBridgebyID(brID uint) (*Bridge, error) {
	var br *Bridge = &Bridge{}
	err := ds.db.Where("id = ?", brID).First(br).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("sqlstore error: bridge ID %v does not exist in db", brID)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return br, nil
}

// Retrieves all Bridges from the Database using Organization ID as reference. returns an error if the query fails
func (ds *Plugin) RetrieveAllBridgesbyOrgID(orgID uint) (*[]Bridge, error) {
	var org *Organization = &Organization{}
	err := ds.db.Preload("Bridges").Where("ID = ?", orgID).Find(org).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("sqlstore error: organization ID %d does not exist in db", orgID)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return &org.Bridges, nil
}

// Retrieves all Members from the Database using bridge ID as reference. returns an error if the query fails
func (ds *Plugin) RetrieveAllMembershipsbyBridgeID(bridgeID uint) (*[]Membership, error) {
	var br *Bridge = &Bridge{}
	err := ds.db.Preload("Memberships").Where("ID = ?", bridgeID).Find(br).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("sqlstore error: bridge %d does not exist in db", bridgeID)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v. bridge query failed with bridge ID %d", err, bridgeID)
	}
	return &br.Memberships, nil
}

// Retrieves all Members from the Database using bridge ID as reference. returns an error if the query fails
func (ds *Plugin) RetrieveAllMembersbyBridgeID(bridgeID uint) (mem *[]Member, err error) {
	var br *Bridge = &Bridge{}
	err = ds.db.Preload("Memberships.member").Where("ID = ?", bridgeID).Find(br).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("sqlstore error: bridge %d does not exist in db", bridgeID)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v. bridge query failed with bridge ID %d", err, bridgeID)
	}
	for _, membership := range br.Memberships {
		*mem = append(*mem, membership.member)
	}
	return mem, nil
}

// Retrieves all Members from the Database using bridge ID as reference. returns an error if the query fails
func (ds *Plugin) RetrieveAllBridgesbyMemberID(memberID uint) (mem *[]Bridge, err error) {
	var member *Member = &Member{}
	err = ds.db.Preload("Memberships.bridge").Where("ID = ?", memberID).Find(member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("sqlstore error: bridge %d does not exist in db", memberID)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v. bridge query failed with bridge ID %d", err, memberID)
	}
	for _, membership := range member.Memberships {
		*mem = append(*mem, membership.bridge)
	}
	return mem, nil
}

// Retrieves all Memberships from the Database using memberID as reference. returns an error if the query fails
func (ds *Plugin) RetrieveAllMembershipsbyMemberID(memberID uint) (*[]Membership, error) {
	var member *Member = &Member{}
	err := ds.db.Preload("Memberships").Where("ID = ?", memberID).Find(member).Error
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v. Membership query failed with member ID %d", err, memberID)
	}
	return &member.Memberships, nil
}

/// Retrieves all Relationships from the Database using memberID as reference. returns an error if the query fails
func (ds *Plugin) RetrieveAllRelationshipsbyMemberID(memberID uint) (*[]Relationship, error) {
	var member *Member = &Member{}
	err := ds.db.Preload("Relationships").Where("ID = ?", memberID).Find(member).Error
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v. relationship query failed with member ID %d", err, memberID)
	}
	return &member.Relationships, nil
}

// Retrieves all Trusts from the Database using memberID as reference. returns an error if the query fails
func (ds *Plugin) RetrieveAllTrustBundlesbyMemberID(memberID uint) (*[]TrustBundle, error) {
	var member *Member = &Member{}
	err := ds.db.Preload("TrustBundles").Where("ID = ?", memberID).Find(member).Error
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v. trust query failed with member id %d", err, memberID)
	}
	return &member.TrustBundles, nil
}

// Retrieves a Member from the Database by description. returns an error if the query fails
func (ds *Plugin) RetrieveMemberbyID(memberID uint) (*Member, error) {
	var member *Member = &Member{}
	err := ds.db.Where("id = ?", memberID).First(member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // If does not, throw an error
		return nil, fmt.Errorf("sqlstore error: Member with ID=%v does not exist in DB", memberID)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return member, nil
}

// RetrieveMembershipbyToken retrieves a Membership from the Database bigger than an specific date. returns an error if something goes wrong.
func (ds *Plugin) RetrieveMembershipbyCreationDate(date time.Time) (*Membership, error) {
	var membership *Membership = &Membership{}
	err := ds.db.Where("created_at >= ?", date).First(membership).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // If does not, throw an error
		return nil, fmt.Errorf("sqlstore error: member created_at=%v does not exist in DB", date)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return membership, nil
}

// RetrieveMembershipbyToken retrieves a Membership from the Database by Token. returns an error if something goes wrong.
func (ds *Plugin) RetrieveMembershipbyToken(token string) (*Membership, error) {
	var membership *Membership = &Membership{}
	err := ds.db.Where("join_token = ?", token).First(membership).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // If does not, throw an error
		return nil, fmt.Errorf("sqlstore error: Member with Token=%v does not exist in DB", token)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return membership, nil
}

// retrieves a Relationship from the Database by Source and Target IDs. returns an error if something goes wrong.
func (ds *Plugin) RetrieveRelationshipbySourceandTargetID(source uint, target uint) (*Relationship, error) {
	var relationship *Relationship = &Relationship{}
	err := ds.db.Where("MemberID = ? AND TargetMemberID = ?", source, target).First(relationship).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // If does not, throw an error
		return nil, fmt.Errorf("sqlstore error: Member with SourceMemberID=%v and/or TargetMemberID=%v does not exist in DB", source, target)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return relationship, nil
}

// retrieves a TrustBundle from the Database by Token. returns an error if something goes wrong.
func (ds *Plugin) RetrieveTrustbundlebyMemberID(memberID string) (*TrustBundle, error) {
	var trustbundle *TrustBundle = &TrustBundle{}
	err := ds.db.Where("MemberID = ?", memberID).First(trustbundle).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // If does not, throw an error
		return nil, fmt.Errorf("sqlstore error: Member with Token=%v does not exist in DB", memberID)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore error: %v", err)
	}
	return trustbundle, nil
}

// UpdateBridge Updates an existing Bridge with the new Bridge as argument. The ID will be used as reference.
func (ds *Plugin) UpdateBridge(br Bridge) error {
	if br.ID == 0 {
		return fmt.Errorf("sqlstore error: Bridge ID is invalid")
	}
	err := ds.db.Model(&br).Updates(&br).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// Updates an existing Organization with the new Organization as argument. The ID will be used as reference.
func (ds *Plugin) UpdateOrganization(org Organization) error {
	if org.ID == 0 {
		return fmt.Errorf("sqlstore error: organization ID is invalid")
	}
	err := ds.db.Model(&org).Updates(&org).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // If does not, throw an error
		return fmt.Errorf("sqlstore error: organization with ID %d does not exist", org.ID)
	}
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// Updates an existing Member with the new Member as argument. The ID will be used as reference.
func (ds *Plugin) UpdateMember(member Member) error {
	if member.ID == 0 {
		return fmt.Errorf("sqlstore error: member id is invalid")
	}
	err := ds.db.Model(&member).Updates(&member).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// Updates an existing Member with the new Member as argument. The ID will be used as reference.
func (ds *Plugin) UpdateMembership(membership Membership) error {
	if membership.ID == 0 {
		return fmt.Errorf("sqlstore error: membership id is invalid")
	}
	err := ds.db.Model(&membership).Updates(&membership).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// Updates an existing Member with the new Member as argument. The ID will be used as reference.
func (ds *Plugin) UpdateTrust(trust TrustBundle) error {
	if trust.ID == 0 {
		return fmt.Errorf("sqlstore error: membership id is invalid")
	}
	err := ds.db.Model(&trust).Updates(&trust).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// Delete Organization by name from the DB with cascading, returning error if something happens
func (ds *Plugin) DeleteOrganizationbyID(orgID uint) error {
	org, err := ds.RetrieveOrganizationbyID(orgID)
	if err != nil {
		return err
	}

	if ds.db.Name() == "sqlite" {
		// Workaround for https://github.com/mattn/go-sqlite3/pull/802 that
		// might prevent DELETE CASCADE on go-sqlite3 driver from working
		brs, err := ds.RetrieveAllBridgesbyOrgID(orgID)
		if err != nil {
			return err
		}
		for _, br := range *brs {
			err = ds.DeleteBridgebyID(br.ID)
			if err != nil {
				return err
			}
		}
	}
	err = ds.db.Model(&org).Delete(&org).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// Delete Organization by Description from the DB with cascading
func (ds *Plugin) DeleteBridgebyID(bridgeID uint) error {
	br, err := ds.RetrieveBridgebyID(bridgeID)
	if err != nil {
		return err
	}
	if ds.db.Name() == "sqlite" {
		// Workaround for https://github.com/mattn/go-sqlite3/pull/802 that
		// might prevent DELETE CASCADE on go-sqlite3 driver from working
		memberships, err := ds.RetrieveAllMembershipsbyBridgeID(bridgeID)
		if err != nil {
			return err
		}
		for _, membership := range *memberships {
			err = ds.DeleteMembershipbyToken(membership.JoinToken)
			if err != nil {
				return err
			}
		}
	}
	// Deletes the Bridge. If its MySQL or Postgres it will cascade automatically by DB model constraint
	err = ds.db.Model(&br).Delete(&br).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// Delete Organization by name from the DB with cascading
func (ds *Plugin) DeleteMemberbyID(memberID uint) error {
	member, err := ds.RetrieveMemberbyID(memberID)
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	if ds.db.Name() == "sqlite" {
		// Workaround for https://github.com/mattn/go-sqlite3/pull/802 that
		// might prevent DELETE CASCADE on go-sqlite3 driver from working
		err = ds.DeleteAllMembershipsbyMemberID(memberID)
		if err != nil {
			return err
		}
		err = ds.DeleteAllRelationshipsbyMemberID(memberID)
		if err != nil {
			return err
		}
		err = ds.DeleteAllTrustbundlesbyMemberID(memberID)
		if err != nil {
			return err
		}

	}
	// Deletes the Member. If its MySQL or Postgres it will cascade automatically by DB model constraint
	err = ds.db.Model(&member).Delete(&member).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// Deletes all Memberships using memberid as FK
func (ds *Plugin) DeleteAllMembershipsbyMemberID(memberid uint) error {
	memberships, err := ds.RetrieveAllMembershipsbyMemberID(memberid)
	if err != nil {
		return err
	}
	for _, membership := range *memberships {
		err = ds.db.Model(&membership).Delete(&membership).Error
		if err != nil {
			return fmt.Errorf("sqlstore error: %v. Error deleting Relationships from member with id %d", err, memberid)
		}
	}
	return nil
}

// Deletes all Memberships using memberid as FK
func (ds *Plugin) DeleteAllMembershipsbyBridgeID(bridgeid uint) error {
	memberships, err := ds.RetrieveAllMembershipsbyBridgeID(bridgeid)
	if err != nil {
		return err
	}
	for _, membership := range *memberships {
		err = ds.db.Model(&membership).Delete(&membership).Error
		if err != nil {
			return fmt.Errorf("sqlstore error: %v. Error deleting Relationships from bridge with id %d", err, bridgeid)
		}
	}
	return nil
}

// Deletes all Relationships using memberid as FK
func (ds *Plugin) DeleteAllRelationshipsbyMemberID(memberid uint) error {
	relations, err := ds.RetrieveAllRelationshipsbyMemberID(memberid)
	if err != nil {
		return err
	}
	for _, relation := range *relations {
		err = ds.db.Model(&relation).Delete(&relation).Error
		if err != nil {
			return fmt.Errorf("sqlstore error: %v. Error deleting Relationships from member with id %d", err, memberid)
		}
	}
	return nil
}

// Deletes all Trusts using memberid as FK
func (ds *Plugin) DeleteAllTrustbundlesbyMemberID(memberid uint) error {
	trusts, err := ds.RetrieveAllTrustBundlesbyMemberID(memberid)
	if err != nil {
		return err
	}
	for _, trust := range *trusts {
		err = ds.db.Model(&trust).Delete(&trust).Error
		if err != nil {
			return fmt.Errorf("sqlstore error: %v. Not able to fully delete trustbundle %s", err, trust.TrustBundle)
		}
	}
	return nil
}

// Delete membership by Token from the DB
func (ds *Plugin) DeleteMembershipbyToken(name string) error {
	membership, err := ds.RetrieveMembershipbyToken(name)
	if err != nil {
		return err
	}
	err = ds.db.Model(&membership).Delete(&membership).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v. Not able to fully delete membership  with token %s", err, membership.JoinToken)
	}
	return nil
}

// Delete Relationship by Source and Target IDs from the DB
func (ds *Plugin) DeleteRelationshipbySourceTargetID(db *gorm.DB, source uint, target uint) error {
	relationship, err := ds.RetrieveRelationshipbySourceandTargetID(source, target)
	if err != nil {
		return err
	}
	err = db.Model(&relationship).Delete(&relationship).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}

// Delete Trusts by MemberID from the DB
func (ds *Plugin) DeleteTrustBundlebyMemberID(db *gorm.DB, memberID string) error {
	trust, err := ds.RetrieveTrustbundlebyMemberID(memberID)
	if err != nil {
		return err
	}
	err = db.Model(&trust).Delete(&trust).Error
	if err != nil {
		return fmt.Errorf("sqlstore error: %v", err)
	}
	return nil
}
