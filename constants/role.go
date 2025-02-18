package constants

type TenantRole string

const (
	TenantAdmin  TenantRole = "TENANT_ADMIN"
	TenantMember TenantRole = "TENANT_MEMBER"
	TenantOwner  TenantRole = "TENANT_OWNER"
)

func (t TenantRole) String() string {
	return string(t)
}
