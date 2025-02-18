package api

const (
	uuidRegexpString = `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
)

type ResourceIDInput struct {
	ID string
}
