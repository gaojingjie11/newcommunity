package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"smartcommunity-microservices/app/agent/rpc/internal/config"
	"smartcommunity-microservices/app/agent/rpc/internal/model"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	arkmodel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultEmbeddingBaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	defaultEmbeddingModel   = "doubao-embedding-text-240715"
	defaultEmbeddingRegion  = "cn-beijing"
	defaultRAGTopK          = 4
	defaultRAGSyncInterval  = 24 * time.Hour
)

type KnowledgeHit struct {
	SourceType string
	SourceID   int64
	Title      string
	Summary    string
	Content    string
	Visibility string
	UpdatedAt  time.Time
	Distance   float64
}

type LatestAIReportResult struct {
	ID        int64
	Title     string
	Summary   string
	Content   string
	UpdatedAt time.Time
}

type arkEmbedder struct {
	client *arkruntime.Client
	model  string
}

func newArkEmbedder(apiKey, baseURL, modelName string) *arkEmbedder {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = defaultEmbeddingBaseURL
	}
	if strings.TrimSpace(modelName) == "" {
		modelName = defaultEmbeddingModel
	}

	client := arkruntime.NewClientWithApiKey(
		apiKey,
		arkruntime.WithBaseUrl(baseURL),
		arkruntime.WithRegion(defaultEmbeddingRegion),
	)
	return &arkEmbedder{
		client: client,
		model:  modelName,
	}
}

func (e *arkEmbedder) EmbedTexts(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	cleaned := make([]string, 0, len(texts))
	for _, text := range texts {
		text = strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
		if text == "" {
			text = "-"
		}
		cleaned = append(cleaned, text)
	}

	if shouldUseMultimodalEmbedding(e.model) {
		return e.embedTextsWithMultimodal(ctx, cleaned)
	}

	resp, err := e.client.CreateEmbeddings(ctx, arkmodel.EmbeddingRequestStrings{
		Input: cleaned,
		Model: e.model,
	})
	if err != nil {
		if shouldFallbackToMultimodal(e.model, err) {
			return e.embedTextsWithMultimodal(ctx, cleaned)
		}
		return nil, err
	}

	vectors := make([][]float32, len(cleaned))
	for _, item := range resp.Data {
		if item.Index < 0 || item.Index >= len(vectors) {
			continue
		}
		vectors[item.Index] = item.Embedding
	}
	for i, vector := range vectors {
		if len(vector) == 0 {
			return nil, fmt.Errorf("embedding response missing vector at index %d", i)
		}
	}
	return vectors, nil
}

func (e *arkEmbedder) embedTextsWithMultimodal(ctx context.Context, texts []string) ([][]float32, error) {
	vectors := make([][]float32, 0, len(texts))
	for _, text := range texts {
		content := text
		resp, err := e.client.CreateMultiModalEmbeddings(ctx, arkmodel.MultiModalEmbeddingRequest{
			Model: e.model,
			Input: []arkmodel.MultimodalEmbeddingInput{{
				Type: arkmodel.MultiModalEmbeddingInputTypeText,
				Text: &content,
			}},
		})
		if err != nil {
			return nil, err
		}
		if len(resp.Data.Embedding) == 0 {
			return nil, fmt.Errorf("multimodal embedding response is empty")
		}
		vectors = append(vectors, resp.Data.Embedding)
	}
	return vectors, nil
}

type KnowledgeService struct {
	db           *gorm.DB
	embedder     *arkEmbedder
	maxResults   int
	syncInterval time.Duration

	mu       sync.Mutex
	lastSync time.Time
}

func NewKnowledgeService(db *gorm.DB, cfg config.AgentConfig) (*KnowledgeService, error) {
	apiKey, baseURL, modelName := cfg.GetModelConfig(cfg.Models.Embedding)
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("embedding API key is not configured")
	}

	maxResults := cfg.RAGMaxResults
	if maxResults <= 0 {
		maxResults = defaultRAGTopK
	}

	syncInterval := time.Duration(cfg.RAGSyncIntervalSeconds) * time.Second
	if syncInterval <= 0 {
		syncInterval = defaultRAGSyncInterval
	}

	return &KnowledgeService{
		db:           db,
		embedder:     newArkEmbedder(apiKey, baseURL, modelName),
		maxResults:   maxResults,
		syncInterval: syncInterval,
	}, nil
}

func (s *KnowledgeService) Init(ctx context.Context) error {
	if err := s.db.WithContext(ctx).Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		return fmt.Errorf("enable pgvector extension failed: %w", err)
	}
	if err := s.db.WithContext(ctx).AutoMigrate(&model.RAGDocument{}, &model.RAGChunk{}); err != nil {
		return fmt.Errorf("migrate rag tables failed: %w", err)
	}
	if err := s.db.WithContext(ctx).
		Exec("ALTER TABLE rag_chunks ALTER COLUMN embedding TYPE vector USING embedding::vector").Error; err != nil {
		return fmt.Errorf("adjust rag embedding column failed: %w", err)
	}
	return nil
}

func (s *KnowledgeService) StartBackgroundSync() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		if err := s.SyncAll(ctx); err != nil {
			log.Printf("[RAG] initial sync failed: %v", err)
		}
		cancel()

		ticker := time.NewTicker(s.syncInterval)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
			if err := s.SyncAll(ctx); err != nil {
				log.Printf("[RAG] periodic sync failed: %v", err)
			}
			cancel()
		}
	}()
}

func (s *KnowledgeService) EnsureFresh(ctx context.Context) error {
	s.mu.Lock()
	shouldSync := s.lastSync.IsZero() || time.Since(s.lastSync) >= s.syncInterval
	s.mu.Unlock()
	if !shouldSync {
		return nil
	}
	return s.SyncAll(ctx)
}

func (s *KnowledgeService) SyncAll(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var notices []model.NoticeSource
	if err := s.db.WithContext(ctx).
		Where("status = ?", 1).
		Order("updated_at DESC, id DESC").
		Find(&notices).Error; err != nil {
		return fmt.Errorf("load notices failed: %w", err)
	}
	for _, notice := range notices {
		if err := s.upsertDocument(ctx, "notice", notice.ID, notice.Title, notice.Content, notice.Content, "public"); err != nil {
			return err
		}
	}
	if err := s.cleanupStaleDocuments(ctx, "notice", collectSourceIDs(notices)); err != nil {
		return err
	}

	var reports []model.AIReportSource
	if err := s.db.WithContext(ctx).
		Order("updated_at DESC, id DESC").
		Find(&reports).Error; err != nil {
		return fmt.Errorf("load ai reports failed: %w", err)
	}
	for _, report := range reports {
		title := fmt.Sprintf("AI 报告 #%d", report.ID)
		content := strings.TrimSpace(report.ReportSummary + "\n\n" + report.Report)
		if err := s.upsertDocument(ctx, "ai_report", report.ID, title, report.ReportSummary, content, "admin_only"); err != nil {
			return err
		}
	}
	if err := s.cleanupStaleDocuments(ctx, "ai_report", collectSourceIDs(reports)); err != nil {
		return err
	}

	s.lastSync = time.Now()
	return nil
}

func (s *KnowledgeService) Search(ctx context.Context, userID int64, query string, scope string, limit int) ([]KnowledgeHit, error) {
	if err := s.EnsureFresh(ctx); err != nil {
		return nil, err
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	if limit <= 0 {
		limit = s.maxResults
	}

	vectors, err := s.embedder.EmbedTexts(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("embed query failed: %w", err)
	}
	if len(vectors) == 0 {
		return nil, nil
	}

	visibility := []string{"public"}
	isAdmin, err := s.isAdminUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		visibility = append(visibility, "admin_only")
	}

	vectorLiteral := vectorLiteral(vectors[0])
	scopeFilter := normalizeKnowledgeScope(scope)

	var rows []struct {
		SourceType string
		SourceID   int64
		Title      string
		Summary    string
		Content    string
		Visibility string
		UpdatedAt  time.Time
		Distance   float64
	}

	sql := `
SELECT *
FROM (
	SELECT DISTINCT ON (d.source_type, d.source_id)
		d.source_type, d.source_id, d.title, d.summary, c.content, d.visibility, d.updated_at,
		(CAST(c.embedding AS vector) <=> CAST(? AS vector)) AS distance
	FROM rag_chunks c
	JOIN rag_documents d ON d.id = c.document_id
	WHERE d.visibility IN ?
`
	args := []interface{}{vectorLiteral, visibility}

	if scopeFilter != "" {
		sql += " AND d.source_type = ?"
		args = append(args, scopeFilter)
	}

	sql += `
	ORDER BY d.source_type, d.source_id, distance ASC
) ranked
ORDER BY distance ASC
LIMIT ?`
	args = append(args, limit)

	if err := s.db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("search rag chunks failed: %w", err)
	}

	hits := make([]KnowledgeHit, 0, len(rows))
	seen := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		key := fmt.Sprintf("%s:%d", row.SourceType, row.SourceID)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		hits = append(hits, KnowledgeHit{
			SourceType: row.SourceType,
			SourceID:   row.SourceID,
			Title:      row.Title,
			Summary:    row.Summary,
			Content:    row.Content,
			Visibility: row.Visibility,
			UpdatedAt:  row.UpdatedAt,
			Distance:   row.Distance,
		})
	}
	return hits, nil
}

func (s *KnowledgeService) CanAccessAdminKnowledge(ctx context.Context, userID int64) (bool, error) {
	return s.isAdminUser(ctx, userID)
}

func (s *KnowledgeService) GetLatestAIReport(ctx context.Context, userID int64) (*LatestAIReportResult, error) {
	allowed, err := s.isAdminUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, nil
	}

	var report model.AIReportSource
	if err := s.db.WithContext(ctx).
		Order("updated_at DESC, id DESC").
		First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	content := strings.TrimSpace(report.Report)
	summary := strings.TrimSpace(report.ReportSummary)
	if summary == "" {
		summary = truncateSummary(content, 220)
	}

	return &LatestAIReportResult{
		ID:        report.ID,
		Title:     fmt.Sprintf("AI 报告 #%d", report.ID),
		Summary:   summary,
		Content:   content,
		UpdatedAt: report.UpdatedAt,
	}, nil
}

func (s *KnowledgeService) upsertDocument(
	ctx context.Context,
	sourceType string,
	sourceID int64,
	title string,
	summary string,
	content string,
	visibility string,
) error {
	title = strings.TrimSpace(title)
	summary = strings.TrimSpace(summary)
	content = strings.TrimSpace(content)
	if title == "" || content == "" {
		return nil
	}

	hash := hashDocument(sourceType, sourceID, title, summary, content, visibility)

	var existing model.RAGDocument
	err := s.db.WithContext(ctx).
		Where("source_type = ? AND source_id = ?", sourceType, sourceID).
		First(&existing).Error
	if err == nil && existing.ContentHash == hash {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("load existing rag document failed: %w", err)
	}

	chunks := chunkText(content, 500, 80)
	embeddings, err := s.embedder.EmbedTexts(ctx, chunks)
	if err != nil {
		return fmt.Errorf("embed source %s:%d failed: %w", sourceType, sourceID, err)
	}
	if len(chunks) != len(embeddings) {
		return fmt.Errorf("embedding count mismatch for %s:%d", sourceType, sourceID)
	}

	now := time.Now()
	doc := model.RAGDocument{
		ID:          existing.ID,
		SourceType:  sourceType,
		SourceID:    sourceID,
		Title:       title,
		Summary:     summary,
		Content:     content,
		Visibility:  visibility,
		ContentHash: hash,
		SyncedAt:    now,
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if existing.ID == 0 {
			if err := tx.Create(&doc).Error; err != nil {
				return err
			}
		} else {
			doc.ID = existing.ID
			if err := tx.Model(&model.RAGDocument{}).
				Where("id = ?", existing.ID).
				Updates(map[string]interface{}{
					"title":        doc.Title,
					"summary":      doc.Summary,
					"content":      doc.Content,
					"visibility":   doc.Visibility,
					"content_hash": doc.ContentHash,
					"synced_at":    doc.SyncedAt,
				}).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("document_id = ?", doc.ID).Delete(&model.RAGChunk{}).Error; err != nil {
			return err
		}

		ragChunks := make([]model.RAGChunk, 0, len(chunks))
		for i, chunk := range chunks {
			ragChunks = append(ragChunks, model.RAGChunk{
				DocumentID: doc.ID,
				ChunkIndex: i,
				Content:    chunk,
				Embedding:  vectorLiteral(embeddings[i]),
			})
		}
		if len(ragChunks) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "document_id"}, {Name: "chunk_index"}},
				UpdateAll: true,
			}).Create(&ragChunks).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *KnowledgeService) isAdminUser(ctx context.Context, userID int64) (bool, error) {
	if userID <= 0 {
		return false, nil
	}
	var user model.SysUser
	if err := s.db.WithContext(ctx).Select("id", "role").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if strings.EqualFold(strings.TrimSpace(user.Role), "admin") {
		return true, nil
	}

	var count int64
	if err := s.db.WithContext(ctx).
		Table("sys_user_role ur").
		Joins("JOIN sys_role r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND LOWER(r.code) = ?", userID, "admin").
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func normalizeKnowledgeScope(scope string) string {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "notice", "notices", "announcement", "announcements":
		return "notice"
	case "report", "reports", "ai_report", "ai reports":
		return "ai_report"
	default:
		return ""
	}
}

func hashDocument(parts ...interface{}) string {
	var sb strings.Builder
	for _, part := range parts {
		sb.WriteString(fmt.Sprint(part))
		sb.WriteString("\n")
	}
	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

func (s *KnowledgeService) cleanupStaleDocuments(ctx context.Context, sourceType string, activeIDs []int64) error {
	query := s.db.WithContext(ctx).Where("source_type = ?", sourceType)
	if len(activeIDs) > 0 {
		query = query.Where("source_id NOT IN ?", activeIDs)
	}

	var staleDocs []model.RAGDocument
	if err := query.Find(&staleDocs).Error; err != nil {
		return fmt.Errorf("load stale rag documents failed: %w", err)
	}
	if len(staleDocs) == 0 {
		return nil
	}

	docIDs := make([]int64, 0, len(staleDocs))
	for _, doc := range staleDocs {
		docIDs = append(docIDs, doc.ID)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("document_id IN ?", docIDs).Delete(&model.RAGChunk{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id IN ?", docIDs).Delete(&model.RAGDocument{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func collectSourceIDs[T interface{ GetID() int64 }](items []T) []int64 {
	ids := make([]int64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.GetID())
	}
	return ids
}

func chunkText(text string, maxRunes int, overlap int) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	if maxRunes <= 0 {
		maxRunes = 500
	}
	if overlap < 0 {
		overlap = 0
	}

	runes := []rune(text)
	if len(runes) <= maxRunes {
		return []string{text}
	}

	chunks := make([]string, 0, (len(runes)/maxRunes)+1)
	step := maxRunes - overlap
	if step <= 0 {
		step = maxRunes
	}

	for start := 0; start < len(runes); start += step {
		end := start + maxRunes
		if end > len(runes) {
			end = len(runes)
		}
		chunk := strings.TrimSpace(string(runes[start:end]))
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
		if end == len(runes) {
			break
		}
	}
	return chunks
}

func vectorLiteral(vector []float32) string {
	if len(vector) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, value := range vector {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%.8f", value))
	}
	sb.WriteByte(']')
	return sb.String()
}

func truncateSummary(text string, maxRunes int) string {
	text = strings.TrimSpace(text)
	if text == "" || maxRunes <= 0 {
		return text
	}
	runes := []rune(text)
	if len(runes) <= maxRunes {
		return text
	}
	return string(runes[:maxRunes]) + "..."
}

func shouldUseMultimodalEmbedding(modelName string) bool {
	name := strings.ToLower(strings.TrimSpace(modelName))
	return strings.Contains(name, "embedding-vision") || strings.Contains(name, "multimodal")
}

func shouldFallbackToMultimodal(modelName string, err error) bool {
	if !shouldUseMultimodalEmbedding(modelName) {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "invalidendpointormodel.notfound") ||
		strings.Contains(message, "not exist") ||
		strings.Contains(message, "do not have access")
}
