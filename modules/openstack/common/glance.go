package common

import (
	"time"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

type imageStatus images.ImageStatus

const (
	ImageStatusQueued        imageStatus = "queued"
	ImageStatusSaving        imageStatus = "saving"
	ImageStatusActive        imageStatus = "active"
	ImageStatusKilled        imageStatus = "killed"
	ImageStatusDeleted       imageStatus = "deleted"
	ImageStatusPendingDelete imageStatus = "pending_delete"
	ImageStatusDeactivated   imageStatus = "deactivated"
	ImageStatusUploading     imageStatus = "uploading"
	ImageStatusImporting     imageStatus = "importing"
)

// ImageStatus
//
//	queued: The Image service reserved an image ID for the image in the catalog but did not yet upload any image data.
//	saving: The Image service is in the process of saving the raw data for the image into the backing store.
//	active: The image is active and ready for consumption in the Image service.
//	killed: An image data upload error occurred.
//	deleted: The Image service retains information about the image but the image is no longer available for use.
//	pending_delete: Similar to the deleted status. An image in this state is not recoverable.
//	deactivated: The image data is not available for use.
//	uploading: Data has been staged as part of the interoperable image import process. It is not yet available for use. (Since Image API 2.6)
//	importing: The image data is being processed as part of the interoperable image import process, but is not yet available for use. (Since Image API 2.6)
func ImageStatus(status images.ImageStatus) imageStatus {
	return imageStatus(status)
}

func (status imageStatus) String() string {
	return string(status)
}

type imageVisibility images.ImageVisibility

const (
	ImageVisibilityPublic  imageVisibility = "public"
	ImageVisibilityPrivate imageVisibility = "private"

	// TODO: 待確認使用情境，暫不開放
	// ImageVisibilityShared    ImageVisibility = "shared"
	// ImageVisibilityCommunity ImageVisibility = "community"
)

func ImageVisibility(visibility images.ImageVisibility) imageVisibility {
	return imageVisibility(visibility)
}

func (visibility imageVisibility) Convert() *images.ImageVisibility {
	if visibility == "" {
		visibility = ImageVisibilityPrivate
	}
	return (*images.ImageVisibility)(&visibility)
}

func (visibility imageVisibility) String() string {
	return string(visibility)
}

type ImageInfo struct {
	ID              string                 // ID is the image UUID.
	Name            string                 // Name is the human-readable display name for the image.
	Status          imageStatus            // Status is the image status. It can be "queued" or "active"
	Tags            []string               // Tags is a list of image tags. Tags are arbitrarily defined strings
	ContainerFormat string                 // ContainerFormat is the format of the container.
	DiskFormat      string                 // DiskFormat is the format of the disk.
	Owner           string                 // Owner is the tenant ID the image belongs to.
	Visibility      imageVisibility        // Visibility defines who can see/use the image.
	SizeBytes       int64                  // SizeBytes is the size of the data that's associated with the image.
	Properties      map[string]interface{} // Properties is a set of key-value pairs, if any, that are associated with the image.
	CreatedAt       time.Time              // Date created.
	UpdatedAt       time.Time              // Date updated.
}

func (i *ImageInfo) ExtractImage(img *images.Image) *ImageInfo {
	i.ID = img.ID
	i.Name = img.Name
	i.Status = ImageStatus(img.Status)
	i.Tags = img.Tags
	i.ContainerFormat = img.ContainerFormat
	i.DiskFormat = img.DiskFormat
	i.Owner = img.Owner
	i.Visibility = ImageVisibility(img.Visibility)
	i.SizeBytes = img.SizeBytes
	i.Properties = img.Properties
	i.CreatedAt = img.CreatedAt
	i.UpdatedAt = img.UpdatedAt
	return i
}

type (
	CreateImageInput struct {
		Name            string
		DiskFormat      string
		ContainerFormat string
		Visibility      imageVisibility
		Tags            []string

		Creator      *string
		RepositoryID *string
		TagID        *string
	}

	CreateImageOutput struct {
		Image ImageInfo
	}
)

func (input *CreateImageInput) Tag(namespace, projectID string) []string {
	return InsertSystemLabelToSlice(input.Tags, namespace, projectID, input.Creator, input.RepositoryID, input.TagID)
}

type (
	UploadImageDataInput struct {
		ID       string
		Filepath string
	}

	UploadImageDataOutput struct{}
)

type (
	DownloadImageDataInput struct {
		ID       string
		Filepath string
	}

	DownloadImageDataOutput struct{}
)

type (
	GetImageInput struct {
		ID string
	}

	GetImageOutput struct {
		Image ImageInfo
	}
)

type (
	UpdateImageInput struct {
		ID         string
		Visibility *imageVisibility
		Tags       *[]string

		DeleteSystemLabel bool

		Creator      *string // update to image tags information
		RepositoryID *string // update to image tags information
		TagID        *string // update to image tags information
	}

	UpdateImageOutput struct {
		Image ImageInfo
	}
)

func (input *UpdateImageInput) UpdateOpts(namespace, projectID string) images.UpdateOpts {
	updateOpts := images.UpdateOpts{}
	if input.Visibility != nil {
		v := input.Visibility.Convert()
		updateOpts = append(updateOpts, images.UpdateVisibility{
			Visibility: *v,
		})
	}
	if input.Tags != nil {
		var t []string
		if input.DeleteSystemLabel {
			t = DeleteSystemLabelFromSlice(*input.Tags, input.RepositoryID, input.TagID)
		} else {
			t = InsertSystemLabelToSlice(*input.Tags, namespace, projectID, input.Creator, input.RepositoryID, input.TagID)
		}

		if len(t) != 0 {
			updateOpts = append(updateOpts, images.ReplaceImageTags{
				NewTags: t,
			})
		}
	}
	return updateOpts
}

type (
	DeleteImageInput struct {
		ID string
	}

	DeleteImageOutput struct{}
)

type (
	ListImagesInput struct {
		Creator   *string
		ProjectID *string
		Status    *imageStatus
		Tags      *[]string
	}
	ListImagesOutput struct {
		Images []*ImageInfo
	}
)

func (input *ListImagesInput) ListOpts(namespace string) images.ListOpts {
	listOpts := images.ListOpts{}

	listOpts.Tags = append(listOpts.Tags, namespaceLabel(namespace))
	if input.Creator != nil {
		listOpts.Tags = append(listOpts.Tags, creatorLabel(*input.Creator))
	}
	if input.ProjectID != nil {
		listOpts.Owner = *input.ProjectID
	}
	if input.Status != nil {
		listOpts.Status = images.ImageStatus(*input.Status)
	}
	if input.Tags != nil {
		listOpts.Tags = *input.Tags
	}
	return listOpts
}
