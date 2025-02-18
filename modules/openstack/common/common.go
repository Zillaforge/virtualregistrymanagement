package common

import "fmt"

const (
	_resourcePrefixName = "vrm-"

	CreatedTag     = _resourcePrefixName + "created"
	_namespaceTag  = _resourcePrefixName + "namespace"
	_creatorTag    = _resourcePrefixName + "creator"
	_projectTag    = _resourcePrefixName + "project-id"
	_repositoryTag = _resourcePrefixName + "repository-id"
	_tagTag        = _resourcePrefixName + "tag-id"
)

func AddPrefixName(name string) string {
	return _resourcePrefixName + name
}

func InsertSystemLabelToSlice(label []string, namespace, projectID string, creator, repositoryID, tagID *string) []string {
	label = append(
		label,
		createdLabel(),
		namespaceLabel(namespace),
		projectLabel(projectID),
	)

	if creator != nil {
		label = append(label, creatorLabel(*creator))
	}
	if repositoryID != nil {
		label = append(label, repositoryLabel(*repositoryID))
	}
	if tagID != nil {
		label = append(label, tagLabel(*tagID))
	}

	return label
}

func DeleteSystemLabelFromSlice(label []string, repositoryID, tagID *string) []string {
	s := []string{}

	deletedSlice := []string{}

	if repositoryID != nil {
		deletedSlice = append(deletedSlice, repositoryLabel(*repositoryID))
	}
	if tagID != nil {
		deletedSlice = append(deletedSlice, tagLabel(*tagID))
	}

	for _, v := range label {
		exist := false
		for _, dV := range deletedSlice {
			if v == dV {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		s = append(s, v)
	}

	return s
}

func InsertSystemLabelToMap(label map[string]string, namespace, projectID string, creator, repositoryID, tagID *string) map[string]string {
	if label == nil {
		label = map[string]string{}
	}

	label[CreatedTag] = ""
	label[_namespaceTag] = namespace
	label[_projectTag] = projectID

	if creator != nil {
		label[_creatorTag] = *creator
	}
	if repositoryID != nil {
		label[_repositoryTag] = *repositoryID
	}
	if tagID != nil {
		label[_tagTag] = *tagID
	}

	return label
}

func DeleteSystemLabelFromMap(label map[string]string, repositoryID, tagID *string) map[string]string {
	m := map[string]string{}

	deletedMap := map[string]string{}

	if repositoryID != nil {
		deletedMap[_repositoryTag] = *repositoryID
	}
	if tagID != nil {
		deletedMap[_tagTag] = *tagID
	}

	for k, v := range label {
		if value, ok := deletedMap[k]; ok && v == value {
			continue
		}
		m[k] = v
	}

	return m
}

func createdLabel() string {
	return CreatedTag
}

func namespaceLabel(namespace string) string {
	return fmt.Sprintf("%s=%s", _namespaceTag, namespace)
}

func creatorLabel(creator string) string {
	return fmt.Sprintf("%s=%s", _creatorTag, creator)
}

func projectLabel(projectID string) string {
	return fmt.Sprintf("%s=%s", _projectTag, projectID)
}

func repositoryLabel(repositoryID string) string {
	return fmt.Sprintf("%s=%s", _repositoryTag, repositoryID)
}

func tagLabel(tagID string) string {
	return fmt.Sprintf("%s=%s", _tagTag, tagID)
}
