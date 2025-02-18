package common

// Pagination ...
type Pagination struct {
	Limit  int
	Offset int
}

func Paginate(limit, offset int32) *Pagination {
	pagination := &Pagination{
		Limit:  -1,
		Offset: 0,
	}

	if limit > 0 {
		pagination.Limit = int(limit)
	}
	if offset > 0 {
		pagination.Offset = int(offset)

		if pagination.Limit == -1 {
			pagination.Limit = 100
		}
	}
	return pagination
}
