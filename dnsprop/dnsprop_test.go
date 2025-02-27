package dnsprop

import (
	"strings"
	"sync"
	"testing"
)

func TestLookupDNS_ARecord(t *testing.T) {
	var wg sync.WaitGroup
	results := make(chan []string, 1)

	wg.Add(1)
	go lookupDNS("8.8.8.8", "Google Public DNS", "example.com", "A", &wg, results)

	wg.Wait()
	close(results)

	for result := range results {
		if len(result) < 4 {
			t.Errorf("Expected at least 4 elements in result, got %v", result)
		}
		if result[2] != "A" {
			t.Errorf("Expected record type A, got %s", result[2])
		}
	}
}

func TestLookupDNS_InvalidDomain(t *testing.T) {
	var wg sync.WaitGroup
	results := make(chan []string, 1)

	wg.Add(1)
	go lookupDNS("8.8.8.8", "Google Public DNS", "invalid...invalid", "A", &wg, results)

	wg.Wait()
	close(results)

	for result := range results {
		if !containsError(result[3]) {
			t.Errorf("Expected error result, got %v", result)
		}
	}
}

func TestLookupDNS_UnsupportedRecordType(t *testing.T) {
	var wg sync.WaitGroup
	results := make(chan []string, 1)

	wg.Add(1)
	go lookupDNS("8.8.8.8", "Google Public DNS", "example.com", "XYZ", &wg, results)

	wg.Wait()
	close(results)

	for result := range results {
		if result[3] != "❌ Unsupported record type" {
			t.Errorf("Expected unsupported record type error, got %v", result)
		}
	}
}

func containsError(msg string) bool {
	return strings.HasPrefix(msg, "❌")
}
