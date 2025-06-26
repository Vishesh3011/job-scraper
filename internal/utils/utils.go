package utils

import (
	"database/sql"
	"fmt"
	"regexp"
)

func ExtractJobID(urn string) (string, error) {
	re := regexp.MustCompile(`urn:li:fsd_jobPostingCard:\((\d+),JOB_DETAILS\)`)
	matches := re.FindStringSubmatch(urn)
	if len(matches) < 2 {
		return "", fmt.Errorf("ID not found in input: %s", urn)
	}
	return matches[1], nil
}

func ToPtr[T any](value T) *T {
	return &value
}

func NullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
