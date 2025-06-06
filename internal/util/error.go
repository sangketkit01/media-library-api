package util

import (
	"log"
)

const (
	NotNullViolationErrCode    = "23502"
	ForeignKeyViolationErrCode = "23503"
	UniqueViolationErrCode     = "23505"
)

func RouteCustomError(err error, path string) {
	log.Printf("[%s] Error: %v\n", path, err)
}
