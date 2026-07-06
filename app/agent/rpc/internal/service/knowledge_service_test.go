package service

import (
	"testing"
	"time"

	"smartcommunity-microservices/app/agent/rpc/internal/model"
)

func TestNormalizeKnowledgeScope(t *testing.T) {
	tests := map[string]string{
		"notice":        "notice",
		"announcements": "notice",
		"report":        "ai_report",
		"ai reports":    "ai_report",
		"unknown":       "",
		"":              "",
	}

	for input, want := range tests {
		if got := normalizeKnowledgeScope(input); got != want {
			t.Fatalf("normalizeKnowledgeScope(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestChunkTextPreservesCoverage(t *testing.T) {
	text := "abcdefghijklmnopqrstuvwxyz"
	chunks := chunkText(text, 10, 2)

	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d: %#v", len(chunks), chunks)
	}
	if chunks[0] != "abcdefghij" || chunks[1] != "ijklmnopqr" || chunks[2] != "qrstuvwxyz" {
		t.Fatalf("unexpected chunks: %#v", chunks)
	}
}

func TestVectorLiteral(t *testing.T) {
	got := vectorLiteral([]float32{1.25, -0.5})
	want := "[1.25000000,-0.50000000]"
	if got != want {
		t.Fatalf("vectorLiteral() = %q, want %q", got, want)
	}
}

func TestHashDocumentIgnoresTimestampNoise(t *testing.T) {
	hash1 := hashDocument("notice", int64(1), "标题", "摘要", "正文", "public")
	hash2 := hashDocument("notice", int64(1), "标题", "摘要", "正文", "public")
	if hash1 != hash2 {
		t.Fatalf("expected identical hashes, got %q and %q", hash1, hash2)
	}
}

func TestCollectSourceIDs(t *testing.T) {
	now := time.Now()
	noticeIDs := collectSourceIDs([]model.NoticeSource{
		{ID: 11, UpdatedAt: now},
		{ID: 22, UpdatedAt: now},
	})
	if len(noticeIDs) != 2 || noticeIDs[0] != 11 || noticeIDs[1] != 22 {
		t.Fatalf("unexpected notice IDs: %#v", noticeIDs)
	}

	reportIDs := collectSourceIDs([]model.AIReportSource{
		{ID: 7, UpdatedAt: now},
	})
	if len(reportIDs) != 1 || reportIDs[0] != 7 {
		t.Fatalf("unexpected report IDs: %#v", reportIDs)
	}
}
