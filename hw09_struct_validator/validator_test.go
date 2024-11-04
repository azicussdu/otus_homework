package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid User",
			in: User{
				ID:     "12345678-1234-5678-1234-567812345678",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid User - wrong ID length",
			in: User{
				ID:     "short-id",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("length must be 36")},
			},
		},
		{
			name: "invalid User - age below min",
			in: User{
				ID:     "12345678-123324-5678-1234-567812345678",
				Age:    17,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("length must be 36")},
				{Field: "Age", Err: fmt.Errorf("must be >= 18")},
			},
		},
		{
			name: "invalid User - email format",
			in: User{
				ID:     "12345678-1234-5678-1234-567812345678",
				Age:    25,
				Email:  "invalid-email",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Email", Err: fmt.Errorf("must match regexp ^\\w+@\\w+\\.\\w+$")},
			},
		},
		{
			name: "valid App version",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "invalid App version length",
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("length must be 5")},
			},
		},
		{
			name: "valid Token struct - no validation tags",
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil, // no errors expected
		},
		{
			name: "valid Response code",
			in: Response{
				Code: 200,
			},
			expectedErr: nil, // no errors expected
		},
		{
			name: "invalid Response code",
			in: Response{
				Code: 403,
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("must be one of [200, 404, 500]")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)

			// Comparing expected and actual errors
			if tt.expectedErr == nil && err != nil {
				t.Errorf("expected no error, but got: %v", err)
				return
			}

			var actErrors, expErrors ValidationErrors
			if tt.expectedErr != nil && errors.As(err, &actErrors) {
				for j, veErr := range actErrors {
					if !errors.As(tt.expectedErr, &expErrors) && (veErr.Field != expErrors[j].Field ||
						veErr.Err.Error() != expErrors.Error()) {
						t.Errorf("expected error for field %s to be %v, but got %v",
							veErr.Field, expErrors[j].Err, veErr.Err)
					}
				}
			}
		})
	}
}
