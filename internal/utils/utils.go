package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

func GetCorrelationID(ctx context.Context) string {
	val := ctx.Value("cid")
	if cid, ok := val.(string); ok {
		return cid
	}
	return ""
}

func ErrToString(err error) string {
	return fmt.Sprint(err)
}

func WaitForTermination(cancel context.CancelFunc) <-chan struct{} {
	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{})
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
		close(doneChan)
	}()

	return doneChan
}

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
