package common

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"pegasus-cloud.com/aes/toolkits/flatten"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

type snapshotStatus string

const (
	SnapshotStatusCreating      snapshotStatus = "creating"
	SnapshotStatusAvailable     snapshotStatus = "available"
	SnapshotStatusBackingUp     snapshotStatus = "backing-up"
	SnapshotStatusDeleting      snapshotStatus = "deleting"
	SnapshotStatusError         snapshotStatus = "error"
	SnapshotStatusDeleted       snapshotStatus = "deleted"
	SnapshotStatusUnmanaging    snapshotStatus = "unmanaging"
	SnapshotStatusRestoring     snapshotStatus = "restoring"
	SnapshotStatusErrorDeleting snapshotStatus = "error_deleting"
)

// Snapshot Status
//
//	creating: The snapshot is being created.
//	available: The snapshot is ready to use.
//	backing-up: The snapshot is being backed up.
//	deleting: The snapshot is being deleted.
//	error: A snapshot creation error occurred.
//	deleted: The snapshot has been deleted.
//	unmanaging: The snapshot is being unmanaged.
//	restoring: The snapshot is being restored to a volume.
//	error_deleting: A snapshot deletion error occurred.
func SnapshotStatus(status string) snapshotStatus {
	return snapshotStatus(status)
}

func (status snapshotStatus) String() string {
	return string(status)
}

type SnapshotInfo struct {
	ID            string            // Unique identifier.
	VolumeID      string            // ID of the Volume from which this Snapshot was created.
	Name          string            // Display name.
	Description   string            // Display description.
	Status        snapshotStatus    // Current status of the Snapshot.
	SizeBytes     int64             // Size depends on scope setting, in Bytes.
	FullSizeBytes int64             // Size of the Snapshot, in Bytes.
	Metadata      map[string]string // User-defined key-value pairs.
	CreatedAt     time.Time         // Date created.
	UpdatedAt     time.Time         // Date updated.
}

func (s *SnapshotInfo) ExtractSnapshot(snapshot *snapshots.Snapshot) *SnapshotInfo {
	sizeSource := mviper.GetString("VirtualRegistryManagement.scopes.snapshot_size")

	s.ID = snapshot.ID
	s.VolumeID = snapshot.VolumeID
	s.Name = snapshot.Name
	s.Description = snapshot.Description
	s.Status = SnapshotStatus(snapshot.Status)
	s.SizeBytes = GetSnapshotSize(snapshot, sizeSource) * 1024 * 1024 * 1024
	s.FullSizeBytes = int64(snapshot.Size) * 1024 * 1024 * 1024
	s.Metadata = snapshot.Metadata
	s.CreatedAt = snapshot.CreatedAt
	s.UpdatedAt = snapshot.UpdatedAt
	return s
}

func GetSnapshotSize(snapshot *snapshots.Snapshot, sizeSource string) int64 {
	m := map[string]interface{}{}
	b, _ := json.Marshal(snapshot)
	json.Unmarshal(b, &m)

	snapshotInfo, _ := flatten.Flatten(m, "", flatten.DotStyle)
	if val, exist := snapshotInfo[sizeSource]; exist {
		switch iVal := val.(type) {
		case string:
			if v, err := strconv.ParseInt(iVal, 10, 64); err == nil {
				return v
			}
		case int:
			return int64(iVal)
		case int8:
			return int64(iVal)
		case int16:
			return int64(iVal)
		case int32:
			return int64(iVal)
		case int64:
			return int64(iVal)
		case uint:
			return int64(iVal)
		case uint8:
			return int64(iVal)
		case uint16:
			return int64(iVal)
		case uint32:
			return int64(iVal)
		case uint64:
			return int64(iVal)
		case float32:
			return int64(iVal)
		case float64:
			return int64(iVal)
		}
	}
	return 0
}

type (
	CreateSnapshotInput struct {
		VolumeID    string
		Force       bool
		Name        string
		Description string
		Metadata    map[string]string

		Creator      *string
		RepositoryID *string
		TagID        *string
	}

	CreateSnapshotOutput struct {
		Snapshot SnapshotInfo
	}
)

func (input *CreateSnapshotInput) Tag(namespace, projectID string) map[string]string {
	return InsertSystemLabelToMap(input.Metadata, namespace, projectID, input.Creator, input.RepositoryID, input.TagID)
}

type (
	GetSnapshotInput struct {
		ID string

		DeleteSystemTag bool

		Creator      *string
		RepositoryID *string
		TagID        *string
	}

	GetSnapshotOutput struct {
		Snapshot SnapshotInfo
	}
)

type (
	UpdateSnapshotMetadataInput struct {
		ID       string
		Metadata *map[string]string

		DeleteSystemLabel bool

		Creator      *string // update to image tags information
		RepositoryID *string // update to image tags information
		TagID        *string // update to image tags information
	}

	UpdateSnapshotMetadataOutput struct {
		Snapshot SnapshotInfo
	}
)

func (input *UpdateSnapshotMetadataInput) UpdateOpts(namespace, projectID string) snapshots.UpdateMetadataOpts {
	updateOpts := snapshots.UpdateMetadataOpts{}
	if input.Metadata != nil {
		var m map[string]interface{}
		if input.DeleteSystemLabel {
			m = input.stringMap2interfaceMap(DeleteSystemLabelFromMap(*input.Metadata, input.RepositoryID, input.TagID))
		} else {
			m = input.stringMap2interfaceMap(InsertSystemLabelToMap(*input.Metadata, namespace, projectID, input.Creator, input.RepositoryID, input.TagID))
		}

		if len(m) != 0 {
			return snapshots.UpdateMetadataOpts{
				Metadata: m,
			}
		}
	}
	return updateOpts
}

func (UpdateSnapshotMetadataInput) stringMap2interfaceMap(input map[string]string) (output map[string]interface{}) {
	output = map[string]interface{}{}
	for k, v := range input {
		output[k] = v
	}
	return
}

type (
	DeleteSnapshotInput struct {
		ID string
	}

	DeleteSnapshotOutput struct{}
)

type (
	ListSnapshotsInput struct {
		ProjectID *string
		VolumeID  *string
		Metadata  *map[string]string
	}
	ListSnapshotsOutput struct {
		Snapshots []*SnapshotInfo
	}
)

// SnapshotListOpts ...
type SnapshotListOpts struct {
	snapshots.ListOptsBuilder
	Metadata map[string]string `q:"metadata"`
}

// ToSnapshotListQuery ...
func (opts SnapshotListOpts) ToSnapshotListQuery() (out string, err error) {
	var parentQ *url.URL
	var myQ *url.URL
	if parentQ, err = gophercloud.BuildQueryString(opts.ListOptsBuilder); err != nil {
		return
	}
	if myQ, err = gophercloud.BuildQueryString(opts); err != nil {
		return
	}
	params := parentQ.Query()
	if myQ.Query().Has("metadata") {
		params.Add("metadata", myQ.Query().Get("metadata"))
	}
	q := &url.URL{RawQuery: params.Encode()}
	out = q.String()
	return
}

func (input *ListSnapshotsInput) ListOpts(namespace string) snapshots.ListOptsBuilder {
	opts := snapshots.ListOpts{
		AllTenants: true,
	}

	if input.ProjectID != nil {
		opts.TenantID = *input.ProjectID
	}

	if input.VolumeID != nil {
		opts.VolumeID = *input.VolumeID
	}

	listOpts := SnapshotListOpts{
		ListOptsBuilder: opts,
	}

	if input.Metadata != nil {
		listOpts.Metadata = *input.Metadata
	}

	return listOpts
}

type volumeStatus string

const (
	VolumeStatusCreating         volumeStatus = "creating"
	VolumeStatusAvailable        volumeStatus = "available"
	VolumeStatusReserved         volumeStatus = "reserved"
	VolumeStatusAttaching        volumeStatus = "attaching"
	VolumeStatusDetaching        volumeStatus = "detaching"
	VolumeStatusInUse            volumeStatus = "in-use"
	VolumeStatusMaintenance      volumeStatus = "maintenance"
	VolumeStatusDeleting         volumeStatus = "deleting"
	VolumeStatusAwaitingTransfer volumeStatus = "awaiting-transfer"
	VolumeStatusError            volumeStatus = "error"
	VolumeStatusErrorDeleting    volumeStatus = "error_deleting"
	VolumeStatusBackingUp        volumeStatus = "backing-up"
	VolumeStatusRestoringBackup  volumeStatus = "restoring-backup"
	VolumeStatusErrorBackingUp   volumeStatus = "error_backing-up"
	VolumeStatusErrorRestoring   volumeStatus = "error_restoring"
	VolumeStatusErrorExtending   volumeStatus = "error_extending"
	VolumeStatusDownloading      volumeStatus = "downloading"
	VolumeStatusUploading        volumeStatus = "uploading"
	VolumeStatusRetyping         volumeStatus = "retyping"
	VolumeStatusExtending        volumeStatus = "extending"
)

func VolumeStatus(status string) volumeStatus {
	return volumeStatus(status)
}

func (status volumeStatus) String() string {
	return string(status)
}

type VolumeInfo struct {
	ID                  string            // Unique identifier for the volume.
	Status              volumeStatus      // Current status of the volume.
	Size                int               // Size of the volume in GB.
	CreatedAt           time.Time         // The date when this volume was created.
	UpdatedAt           time.Time         // The date when this volume was last updated
	Name                string            // Human-readable display name for the volume.
	Description         string            // Human-readable description for the volume
	VolumeType          string            // The type of volume to create, either SATA or SSD.
	SnapshotID          string            // The ID of the snapshot from which the volume was created
	SourceVolID         string            // The ID of another block storage volume from which the current volume was created
	BackupID            *string           // The backup ID, from which the volume was restored// This field is supported since 3.47 microversion
	Metadata            map[string]string // Arbitrary key-value pairs defined by the user.
	UserID              string            // UserID is the id of the user who created the volume.
	Bootable            string            // Indicates whether this is a bootable volume.
	Encrypted           bool              // Encrypted denotes if the volume is encrypted.
	ReplicationStatus   string            // ReplicationStatus is the status of replication.
	ConsistencyGroupID  string            // ConsistencyGroupID is the consistency group ID.
	Multiattach         bool              // Multiattach denotes if the volume is multi-attach capable.
	VolumeImageMetadata map[string]string // Image metadata entries, only included for volumes that were created from an image, or from a snapshot of a volume originally created from an image.
}

func (s *VolumeInfo) ExtractVolume(volume *volumes.Volume) *VolumeInfo {
	s.ID = volume.ID
	s.Status = VolumeStatus(volume.Status)
	s.Size = volume.Size
	s.CreatedAt = volume.CreatedAt
	s.UpdatedAt = volume.UpdatedAt
	s.Name = volume.Name
	s.Description = volume.Description
	s.VolumeType = volume.VolumeType
	s.SnapshotID = volume.SnapshotID
	s.SourceVolID = volume.SourceVolID
	s.BackupID = volume.BackupID
	s.Metadata = volume.Metadata
	s.UserID = volume.UserID
	s.Bootable = volume.Bootable
	s.Encrypted = volume.Encrypted
	s.ReplicationStatus = volume.ReplicationStatus
	s.ConsistencyGroupID = volume.ConsistencyGroupID
	s.Multiattach = volume.Multiattach
	s.VolumeImageMetadata = volume.VolumeImageMetadata
	return s
}

type (
	CreateVolumeInput struct {
		SnapshotID  string            // the ID of the existing volume snapshot
		Size        int               // The size of the volume, in GB
		Name        string            // The volume name
		Description string            // The volume description
		Metadata    map[string]string // One or more metadata key and value pairs to associate with the volume
	}

	CreateVolumeOutput struct {
		Volume VolumeInfo
	}
)

type (
	GetVolumeInput struct {
		ID string
	}

	GetVolumeOutput struct {
		Volume VolumeInfo
	}
)

type VolumeImageInfo struct {
	VolumeID        string      // The ID of a volume an image is created from.
	ContainerFormat string      // Container format, may be bare, ofv, ova, etc.
	DiskFormat      string      // Disk format, may be raw, qcow2, vhd, vdi, vmdk, etc.
	Description     string      // Human-readable description for the volume.
	ImageID         string      // The ID of the created image.
	ImageName       string      // Human-readable display name for the image.
	SizeBytes       int         // Size of the volume in Bytes.
	Status          imageStatus // Current status of the volume.
	Visibility      string      // Visibility defines who can see/use the image.	// supported since 3.1 microversion
	Protected       bool        // whether the image is not deletable.	// supported since 3.1 microversion
	UpdatedAt       time.Time   // The date when this volume was last updated.
}

func (s *VolumeImageInfo) ExtractImage(image volumeactions.VolumeImage) *VolumeImageInfo {
	s.VolumeID = image.VolumeID
	s.ContainerFormat = image.ContainerFormat
	s.DiskFormat = image.DiskFormat
	s.Description = image.Description
	s.ImageID = image.ImageID
	s.ImageName = image.ImageName
	s.SizeBytes = image.Size * 1024 * 1024 * 1024
	s.Status = imageStatus(image.Status)
	s.Visibility = image.Visibility
	s.Protected = image.Protected
	s.UpdatedAt = image.UpdatedAt
	return s
}

type (
	UploadImageFromVolumeInput struct {
		ID   string
		Name string
	}

	UploadImageFromVolumeOutput struct {
		Image VolumeImageInfo
	}
)

type (
	DeleteVolumeInput struct {
		ID string
	}

	DeleteVolumeOutput struct{}
)
