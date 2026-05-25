package domain

import "testing"

func TestDraftStatus_IsValid(t *testing.T) {
	tests := []struct {
		status DraftStatus
		want   bool
	}{
		{DraftStatusDraft, true},
		{DraftStatusPending, true},
		{DraftStatusApproved, true},
		{DraftStatusRejected, true},
		{DraftStatus("invalid"), false},
		{DraftStatus(""), false},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("DraftStatus(%q).IsValid() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestDraftStatus_IsEditable(t *testing.T) {
	tests := []struct {
		status DraftStatus
		want   bool
	}{
		{DraftStatusDraft, true},
		{DraftStatusRejected, true},
		{DraftStatusPending, false},
		{DraftStatusApproved, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsEditable(); got != tt.want {
				t.Errorf("DraftStatus(%q).IsEditable() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestDraftStatus_String(t *testing.T) {
	if got := DraftStatusDraft.String(); got != "draft" {
		t.Errorf("DraftStatusDraft.String() = %q, want %q", got, "draft")
	}
}
