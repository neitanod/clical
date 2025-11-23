package calendar

import (
	"testing"
	"time"
)

func TestFilterMatches(t *testing.T) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	entry := &Entry{
		DateTime: now,
		Title:    "Test Meeting",
		Duration: 60,
		Location: "Office",
		Notes:    "Important meeting with client",
		Tags:     []string{"trabajo", "cliente"},
	}

	tests := []struct {
		name    string
		filter  *Filter
		matches bool
	}{
		{
			name:    "empty filter matches all",
			filter:  NewFilter(),
			matches: true,
		},
		{
			name: "date range matches",
			filter: &Filter{
				From: &now,
				To:   &tomorrow,
			},
			matches: true,
		},
		{
			name: "date range doesn't match (before)",
			filter: &Filter{
				From: &tomorrow,
			},
			matches: false,
		},
		{
			name: "query matches title",
			filter: &Filter{
				Query: "meeting",
			},
			matches: true,
		},
		{
			name: "query matches notes",
			filter: &Filter{
				Query: "client",
			},
			matches: true,
		},
		{
			name: "query doesn't match",
			filter: &Filter{
				Query: "nonexistent",
			},
			matches: false,
		},
		{
			name: "title filter matches",
			filter: &Filter{
				Title: "meeting",
			},
			matches: true,
		},
		{
			name: "location filter matches",
			filter: &Filter{
				Location: "office",
			},
			matches: true,
		},
		{
			name: "tags match (all required)",
			filter: &Filter{
				Tags: []string{"trabajo", "cliente"},
			},
			matches: true,
		},
		{
			name: "tags don't match (missing one)",
			filter: &Filter{
				Tags: []string{"trabajo", "missing"},
			},
			matches: false,
		},
		{
			name: "min duration matches",
			filter: &Filter{
				MinDuration: 30,
			},
			matches: true,
		},
		{
			name: "min duration doesn't match",
			filter: &Filter{
				MinDuration: 90,
			},
			matches: false,
		},
		{
			name: "max duration matches",
			filter: &Filter{
				MaxDuration: 90,
			},
			matches: true,
		},
		{
			name: "max duration doesn't match",
			filter: &Filter{
				MaxDuration: 30,
			},
			matches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(entry)
			if got != tt.matches {
				t.Errorf("Matches() = %v, want %v", got, tt.matches)
			}
		})
	}
}

func TestFilterWithDateRange(t *testing.T) {
	from := time.Now()
	to := from.Add(7 * 24 * time.Hour)

	filter := NewFilter().WithDateRange(from, to)

	if filter.From == nil || !filter.From.Equal(from) {
		t.Errorf("Expected From to be %v", from)
	}

	if filter.To == nil || !filter.To.Equal(to) {
		t.Errorf("Expected To to be %v", to)
	}
}

func TestFilterWithQuery(t *testing.T) {
	query := "test query"
	filter := NewFilter().WithQuery(query)

	if filter.Query != query {
		t.Errorf("Expected Query to be %s, got %s", query, filter.Query)
	}
}

func TestFilterWithTags(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	filter := NewFilter().WithTags(tags...)

	if len(filter.Tags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(filter.Tags))
	}

	for i, tag := range tags {
		if filter.Tags[i] != tag {
			t.Errorf("Expected tag %s at position %d, got %s", tag, i, filter.Tags[i])
		}
	}
}
