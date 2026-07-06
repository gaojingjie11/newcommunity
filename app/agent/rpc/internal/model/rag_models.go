package model

import "time"

type RAGDocument struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	SourceType  string    `gorm:"type:varchar(32);not null;uniqueIndex:uk_rag_source,priority:1;index" json:"source_type"`
	SourceID    int64     `gorm:"not null;uniqueIndex:uk_rag_source,priority:2;index" json:"source_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Summary     string    `gorm:"type:text" json:"summary"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Visibility  string    `gorm:"type:varchar(16);not null;index" json:"visibility"`
	ContentHash string    `gorm:"type:char(64);not null;index" json:"content_hash"`
	SyncedAt    time.Time `json:"synced_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (RAGDocument) TableName() string { return "rag_documents" }

type RAGChunk struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	DocumentID int64     `gorm:"not null;index;uniqueIndex:uk_rag_chunk,priority:1" json:"document_id"`
	ChunkIndex int       `gorm:"not null;uniqueIndex:uk_rag_chunk,priority:2" json:"chunk_index"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Embedding  string    `gorm:"type:vector;not null" json:"embedding"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (RAGChunk) TableName() string { return "rag_chunks" }

// NoticeSource mirrors the community notice table for read-only indexing.
type NoticeSource struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Publisher string    `json:"publisher"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NoticeSource) TableName() string { return "notices" }

func (n NoticeSource) GetID() int64 { return n.ID }

// AIReportSource mirrors the statistics report table for read-only indexing.
type AIReportSource struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	ReportSummary string    `gorm:"column:report_summary" json:"report_summary"`
	Report        string    `gorm:"column:report_markdown" json:"report"`
	GeneratedBy   int64     `gorm:"column:generated_by" json:"generated_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (AIReportSource) TableName() string { return "cms_ai_report" }

func (r AIReportSource) GetID() int64 { return r.ID }

// SysUser is a lean read-only reference for shared RBAC checks.
type SysUser struct {
	ID   int64  `gorm:"primaryKey" json:"id"`
	Role string `gorm:"type:varchar(32)" json:"role"`
}

func (SysUser) TableName() string { return "sys_user" }

type SysUserRole struct {
	UserID int64 `gorm:"primaryKey" json:"user_id"`
	RoleID int64 `gorm:"primaryKey" json:"role_id"`
}

func (SysUserRole) TableName() string { return "sys_user_role" }

type SysRole struct {
	ID   int64  `gorm:"primaryKey" json:"id"`
	Code string `gorm:"type:varchar(64)" json:"code"`
}

func (SysRole) TableName() string { return "sys_role" }
