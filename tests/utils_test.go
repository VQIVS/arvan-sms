package tests

import (
	"testing"
	"time"
)

func TestTimeUtilities(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Hour)
	past := now.Add(-1 * time.Hour)

	if !future.After(now) {
		t.Error("Future time should be after now")
	}
	if !past.Before(now) {
		t.Error("Past time should be before now")
	}
}

func TestPhoneNumberValidation(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		isValid     bool
	}{
		{"Valid international format", "+1234567890", true},
		{"Valid international format with country code", "+989123456789", true},
		{"Invalid - no plus sign", "1234567890", false},
		{"Invalid - too short", "+123", false},
		{"Invalid - empty", "", false},
		{"Invalid - contains letters", "+123abc4567", false},
		{"Invalid - contains spaces", "+123 456 7890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := isValidE164PhoneNumber(tt.phoneNumber)
			if isValid != tt.isValid {
				t.Errorf("Expected %s to be valid: %v, got: %v", tt.phoneNumber, tt.isValid, isValid)
			}
		})
	}
}

func isValidE164PhoneNumber(phone string) bool {
	if len(phone) == 0 {
		return false
	}
	if phone[0] != '+' {
		return false
	}
	if len(phone) < 5 || len(phone) > 16 {
		return false
	}
	for i := 1; i < len(phone); i++ {
		if phone[i] < '0' || phone[i] > '9' {
			return false
		}
	}
	return true
}

func TestSMSContentValidation(t *testing.T) {
	tests := []struct {
		name    string
		content string
		isValid bool
	}{
		{"Valid short message", "Hello", true},
		{"Valid max length message", generateString(160), true},
		{"Invalid - too long", generateString(161), false},
		{"Valid - empty allowed in some cases", "", true},
		{"Valid - with emojis", "Hello 😊", true},
		{"Valid - with numbers", "Your code is 123456", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := isValidSMSContent(tt.content)
			if isValid != tt.isValid {
				t.Errorf("Expected content validity to be %v, got %v for content length %d", tt.isValid, isValid, len(tt.content))
			}
		})
	}
}

func isValidSMSContent(content string) bool {
	return len(content) <= 160
}

func generateString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}

func TestStatusTransitions(t *testing.T) {
	tests := []struct {
		name              string
		initialStatus     string
		targetStatus      string
		isValidTransition bool
	}{
		{"Pending to Delivered", "pending", "delivered", true},
		{"Pending to Failed", "pending", "failed", true},
		{"Delivered to Failed", "delivered", "failed", false},
		{"Failed to Delivered", "failed", "delivered", false},
		{"Failed to Pending", "failed", "pending", false},
		{"Delivered to Pending", "delivered", "pending", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := isValidStatusTransition(tt.initialStatus, tt.targetStatus)
			if isValid != tt.isValidTransition {
				t.Errorf("Expected transition from %s to %s to be valid: %v, got: %v",
					tt.initialStatus, tt.targetStatus, tt.isValidTransition, isValid)
			}
		})
	}
}

func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		"pending":   {"delivered", "failed"},
		"delivered": {},
		"failed":    {},
	}

	validTargets, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, validTarget := range validTargets {
		if validTarget == to {
			return true
		}
	}
	return false
}
