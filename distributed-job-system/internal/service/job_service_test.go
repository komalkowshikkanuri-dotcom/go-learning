package service

import "testing"

func TestIsRetryable(t *testing.T) {
	svc := &JobService{}

	tests := []struct {
		currentAttempt  int
		maximumAttempts int
		expected        bool
	}{
		{1, 3, true},
		{2, 3, true},
		{3, 3, false},
		{1, 1, false},
		{0, 3, false},
	}

	for _, tc := range tests {
		got := svc.CanRetry(tc.currentAttempt, tc.maximumAttempts)

		if got != tc.expected {
			t.Errorf("CanRetry(%d, %d) = %v; want %v", tc.currentAttempt, tc.maximumAttempts, got, tc.expected)
		}
	}
}
