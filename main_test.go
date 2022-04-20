package main

import (
	"fmt"
	"testing"
)

func TestAuthUser(t *testing.T) {
	testCases := []struct {
		U, P     string
		Expected bool
	}{
		{"", "", false},
		{"notadmin", "", false},
		{"", "pass", false},
		{"admin", "", false},
		{"user", "pass", false},
		{AdminUser, AdminPass, true},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			res := authUser(tc.U, tc.P)
			if res != tc.Expected {
				t.Errorf("expected %t, got %t", tc.Expected, res)
			}
		})
	}
}
