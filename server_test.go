package nexmo

import (
	"testing"
	"time"
)

func TestParseSCTS(t *testing.T) {
	if result, err := parseSCTS(""); err != nil {
		t.Errorf("Failed to parse empty string: %v", err)
	} else if !result.IsZero() {
		t.Errorf("Did not get a zero time from an empty string: %v", result)
	}
}

func TestParseMessageTimestamp(t *testing.T) {
	if result, err := parseMessageTimestamp(""); err != nil {
		t.Errorf("Failed to parse empty string: %v", err)
	} else if !result.IsZero() {
		t.Errorf("Did not get a zero time from an empty string: %v", result)
	}

	if result, err := parseMessageTimestamp("2022-05-05 14:35:48 0000"); err != nil {
		t.Errorf("Failed to parse empty string: %v", err)
	} else if !result.Equal(time.Date(2022, 5, 5, 14, 35, 48, 0, time.UTC)) {
		t.Errorf("Wrong time: %v", result)
	}

	if result, err := parseMessageTimestamp("2022-05-05 14:35:48  0000"); err != nil {
		t.Errorf("Failed to parse empty string: %v", err)
	} else if !result.Equal(time.Date(2022, 5, 5, 14, 35, 48, 0, time.UTC)) {
		t.Errorf("Wrong time: %v", result)
	}

	if result, err := parseMessageTimestamp("2022-05-05 16:05:13 +0000"); err != nil {
		t.Errorf("Failed to parse empty string: %v", err)
	} else if !result.Equal(time.Date(2022, 5, 5, 16, 5, 13, 0, time.UTC)) {
		t.Errorf("Wrong time: %v", result)
	}
}
