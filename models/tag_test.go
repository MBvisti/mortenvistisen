package models

import (
	"strings"
	"testing"

	"mortenvistisen/internal/validation"
)

func TestTagEntityValidate(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		wantCode string
	}{
		{name: "valid", title: "Go"},
		{name: "title is required", title: " ", wantCode: "required"},
		{
			name:     "title is too long",
			title:    strings.Repeat("a", tagTitleMaxLength+1),
			wantCode: "max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := TagEntity{Title: tt.title}
			err := tag.Validate()
			if tt.wantCode == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}

			errors, ok := validation.As(err)
			if !ok {
				t.Fatalf("Validate() error = %v, want validation errors", err)
			}
			for _, validationErr := range errors {
				if validationErr.Field == "title" && validationErr.Code == tt.wantCode {
					return
				}
			}

			t.Fatalf("Validate() errors = %#v, want title code %q", errors, tt.wantCode)
		})
	}
}

func TestUniquePositiveIDs(t *testing.T) {
	got := uniquePositiveIDs([]int32{3, 0, 2, 3, -1, 2, 7})
	want := []int32{3, 2, 7}
	if len(got) != len(want) {
		t.Fatalf("uniquePositiveIDs() = %#v, want %#v", got, want)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("uniquePositiveIDs() = %#v, want %#v", got, want)
		}
	}
}
