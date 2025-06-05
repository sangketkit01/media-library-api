package util

import "fmt"

const (
	NotNullViolationErrCode    = "23502"
	ForeignKeyViolationErrCode = "23503"
	UniqueViolationErrCode     = "23505"
)

func RouteCustomError(err error, path string) string {
	return fmt.Sprintf("[%s] Error: %v", path, err)
}