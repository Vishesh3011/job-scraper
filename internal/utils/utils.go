package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

func WaitForTermination(cancel context.CancelFunc) <-chan struct{} {
	sig := make(chan os.Signal, 1)
	done := make(chan struct{})

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("Shutting down...at ", time.Now())
		cancel()
		close(done)
	}()

	return done
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
