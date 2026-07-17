package models

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"mortenvistisen/internal/validation"
)

func validProjectEntity() ProjectEntity {
	return ProjectEntity{
		Title: "Valid project",
		Slug:  "valid-project",
	}
}

func TestProjectEntityValidate(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*ProjectEntity)
		wantField string
		wantCode  string
	}{
		{name: "valid draft"},
		{
			name: "draft allows empty optional URL",
			mutate: func(project *ProjectEntity) {
				project.ProjectUrl = normalizeOptionalString(sql.NullString{Valid: true})
			},
		},
		{
			name: "title is required",
			mutate: func(project *ProjectEntity) {
				project.Title = " "
			},
			wantField: "title",
			wantCode:  "required",
		},
		{
			name: "title is too long",
			mutate: func(project *ProjectEntity) {
				project.Title = strings.Repeat("a", projectTitleMaxLength+1)
			},
			wantField: "title",
			wantCode:  "max",
		},
		{
			name: "slug is required",
			mutate: func(project *ProjectEntity) {
				project.Slug = ""
			},
			wantField: "slug",
			wantCode:  "required",
		},
		{
			name: "status is too long",
			mutate: func(project *ProjectEntity) {
				project.Status = strings.Repeat("a", projectStatusMaxLength+1)
			},
			wantField: "status",
			wantCode:  "max",
		},
		{
			name: "meta title is too long",
			mutate: func(project *ProjectEntity) {
				project.MetaTitle = strings.Repeat("a", projectMetaTitleMaxLength+1)
			},
			wantField: "metaTitle",
			wantCode:  "max",
		},
		{
			name: "meta description is too long",
			mutate: func(project *ProjectEntity) {
				project.MetaDescription = strings.Repeat(
					"a",
					projectMetaDescriptionMaxLength+1,
				)
			},
			wantField: "metaDescription",
			wantCode:  "max",
		},
		{
			name: "social image URL is invalid",
			mutate: func(project *ProjectEntity) {
				project.MetaImageLink = sql.NullString{String: "not-a-url", Valid: true}
			},
			wantField: "metaImageLink",
			wantCode:  "url",
		},
		{
			name: "project URL is invalid",
			mutate: func(project *ProjectEntity) {
				project.ProjectUrl = sql.NullString{String: "not-a-url", Valid: true}
			},
			wantField: "projectUrl",
			wantCode:  "url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := validProjectEntity()
			if tt.mutate != nil {
				tt.mutate(&project)
			}

			err := project.Validate()
			if tt.wantField == "" {
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
				if validationErr.Field == tt.wantField && validationErr.Code == tt.wantCode {
					return
				}
			}

			t.Fatalf(
				"Validate() errors = %#v, want field %q with code %q",
				errors,
				tt.wantField,
				tt.wantCode,
			)
		})
	}
}

func TestPublishedProjectRequiresCompleteContent(t *testing.T) {
	validPublishedProject := func() ProjectEntity {
		return ProjectEntity{
			Published:       true,
			Title:           "Complete project",
			Slug:            "complete-project",
			StartedAt:       sql.NullTime{Time: time.Now(), Valid: true},
			Status:          "active",
			Description:     "A useful project summary",
			MetaTitle:       "Search title",
			MetaDescription: "Search description",
			ImageLink: sql.NullString{
				String: "https://media.example.com/project.png",
				Valid:  true,
			},
			Content: "# Complete project",
		}
	}

	tests := []struct {
		name      string
		mutate    func(*ProjectEntity)
		wantField string
	}{
		{name: "complete project"},
		{
			name: "start date is required",
			mutate: func(project *ProjectEntity) {
				project.StartedAt = sql.NullTime{}
			},
			wantField: "startedAt",
		},
		{
			name: "status is required",
			mutate: func(project *ProjectEntity) {
				project.Status = ""
			},
			wantField: "status",
		},
		{
			name: "description is required",
			mutate: func(project *ProjectEntity) {
				project.Description = ""
			},
			wantField: "description",
		},
		{
			name: "meta title is required",
			mutate: func(project *ProjectEntity) {
				project.MetaTitle = ""
			},
			wantField: "metaTitle",
		},
		{
			name: "meta description is required",
			mutate: func(project *ProjectEntity) {
				project.MetaDescription = ""
			},
			wantField: "metaDescription",
		},
		{
			name: "image is required",
			mutate: func(project *ProjectEntity) {
				project.ImageLink = sql.NullString{}
			},
			wantField: "imageLink",
		},
		{
			name: "content is required",
			mutate: func(project *ProjectEntity) {
				project.Content = ""
			},
			wantField: "content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := validPublishedProject()
			if tt.mutate != nil {
				tt.mutate(&project)
			}

			err := project.Validate()
			if tt.wantField == "" {
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
				if validationErr.Field == tt.wantField && validationErr.Code == "required" {
					return
				}
			}

			t.Fatalf("Validate() errors = %#v, want required error for %q", errors, tt.wantField)
		})
	}
}
