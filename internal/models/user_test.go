package models

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name string
		dob  string
		now  string
		age  int
	}{
		{
			name: "birthday has passed this year",
			dob:  "1990-01-10",
			now:  "2025-12-31",
			age:  35,
		},
		{
			name: "birthday is today",
			dob:  "1990-05-10",
			now:  "2025-05-10",
			age:  35,
		},
		{
			name: "birthday has not yet happened this year",
			dob:  "1990-12-31",
			now:  "2025-01-01",
			age:  34,
		},
		{
			name: "future dob",
			dob:  "2100-01-01",
			now:  "2025-01-01",
			age:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dob, _ := time.Parse("2006-01-02", tt.dob)
			now, _ := time.Parse("2006-01-02", tt.now)
			got := CalculateAge(dob, now)
			if got != tt.age {
				t.Fatalf("expected %d, got %d", tt.age, got)
			}
		})
	}
}


