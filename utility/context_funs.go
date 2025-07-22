package utility

import (
	"context"

	"github.com/Zillaforge/toolkits/tracer"
)

// MustGetContextRequestID ...
func MustGetContextRequestID(ctx context.Context) string {
	requestID := ctx.Value(tracer.RequestID)
	if requestID == nil {
		return tracer.EmptyRequestID
	}
	rid, ok := requestID.(string)
	if ok {
		return rid
	} else {
		return tracer.EmptyRequestID
	}
}
