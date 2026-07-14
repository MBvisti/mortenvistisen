package admin

import "testing"

func TestArticleStatusFilter(t *testing.T) {
	tests := []struct {
		status string
		want   *bool
	}{
		{status: "published", want: boolPointer(true)},
		{status: "draft", want: boolPointer(false)},
		{status: "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := articleStatusFilter(tt.status).Published
			if got == nil || tt.want == nil {
				if got != nil || tt.want != nil {
					t.Fatalf("Published = %v, want %v", got, tt.want)
				}
				return
			}
			if *got != *tt.want {
				t.Fatalf("Published = %t, want %t", *got, *tt.want)
			}
		})
	}
}

func boolPointer(value bool) *bool {
	return &value
}
