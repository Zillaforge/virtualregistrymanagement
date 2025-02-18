package tables

import (
	"time"

	"gorm.io/datatypes"
)

// Registry struct merge the all tags information
type Registry struct {
	RepositoryID        string
	TagID               string
	Creator             string
	ProjectID           string
	Namespace           string
	RepositoryName      string
	TagName             string
	Description         string
	OperatingSystem     string
	Type                string
	Size                uint64
	Status              string
	Extra               datatypes.JSON
	ReferenceTarget     string
	MemberAclID         string
	AllowUserID         string
	ProjectAclID        string
	AllowProjectID      *string
	RepositoryCreatedAt time.Time
	RepositoryUpdatedAt time.Time
	TagCreatedAt        time.Time
	TagUpdatedAt        time.Time

	Project    Project    `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Repository Repository `gorm:"foreignKey:RepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tag        Tag        `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MemberAcl  MemberAcl  `gorm:"foreignKey:MemberAclID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectAcl ProjectAcl `gorm:"foreignKey:ProjectAclID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	_          struct{}
}

// RegistryView ...
func RegistryView() string {
	return `CREATE VIEW registry AS 
	SELECT
		repository.id AS repository_id,
		tag.id AS tag_id,
		repository.creator AS creator,
		repository.project_id AS project_id,
		repository.namespace AS namespace,
        
		repository.name AS repository_name,
		tag.name AS tag_name,
		tag.reference_target AS reference_target,
		repository.description AS description,
		repository.operating_system AS operating_system,
		tag.type AS type,
		tag.size AS size,
		tag.status AS status,
		tag.extra AS extra,

		member_acl.id AS member_acl_id,
		member_acl.user_id AS allow_user_id,
		project_acl.id AS project_acl_id,
		project_acl.project_id AS allow_project_id,

		repository.created_at AS repository_created_at,
		repository.updated_at AS repository_updated_at,
		tag.created_at AS tag_created_at,
		tag.updated_at AS tag_updated_at
	FROM
		repository
	LEFT JOIN tag ON
		repository.id = tag.repository_id
	LEFT JOIN member_acl ON
		member_acl.tag_id = tag.id
	LEFT JOIN project_acl ON
		project_acl.tag_id = tag.id`
}
