package utils

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

func ReadMultilineInput(prompt string) (string, error) {
	fmt.Println(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}

func WaitForTermination(cancel context.CancelFunc, logger *slog.Logger) <-chan struct{} {
	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{})
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
		logger.Info(fmt.Sprintf("Received termination signal, shutting down..."))
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

func ToSQLNullStr(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func CheckIfResponseIsJSON(body []byte) error {
	bodyStr := strings.TrimSpace(string(body))
	if !strings.HasPrefix(bodyStr, "{") && !strings.HasPrefix(bodyStr, "[") {
		return fmt.Errorf("API returned error: %s", bodyStr)
	}
	return nil
}
