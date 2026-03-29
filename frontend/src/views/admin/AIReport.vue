<template>
  <div class="ai-report-page">
    <Navbar />

    <div class="container custom-container">
      <!-- 顶部返回导航 -->
      <div class="page-nav">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回管理后台</span>
        </div>
      </div>

      

      <!-- 核心数据面板 -->
      <div class="premium-card">
        <div class="panel-header">
          <div class="panel-title">
            <div class="icon-box">
              <el-icon><DataAnalysis /></el-icon>
            </div>
            <span>智能日报汇编</span>
          </div>
          <div class="header-actions">
            <button class="action-btn btn-success" :disabled="generating" @click="handleGenerate(true)">
              <el-icon v-if="generating" class="is-loading"><Loading /></el-icon>
              <el-icon v-else><Plus /></el-icon>
              {{ generating ? 'AI 生成中...' : '手动生成报告' }}
            </button>
            <button class="action-btn btn-primary" :disabled="listLoading" @click="fetchList">
              <el-icon :class="{ 'is-loading': listLoading }"><Refresh /></el-icon>
              刷新列表
            </button>
          </div>
        </div>

        <!-- 深度定制的报表列表 -->
        <el-table :data="list" v-loading="listLoading" class="custom-table" :empty-text="'暂无生成的 AI 报告'">
          <el-table-column label="报告生成时间" width="220">
            <template #default="{ row }">
              <div class="time-cell">
                <el-icon><Clock /></el-icon>
                <span>{{ formatDateTime(row.created_at) }}</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="report_summary" label="核心摘要" min-width="500">
            <template #default="{ row }">
              <span class="summary-text">{{ row.report_summary }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="140" align="center" fixed="right">
            <template #default="{ row }">
              <button class="btn-text" @click="openDetail(row.id)">查看完整研报</button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页器 -->
        <div class="pagination-wrap" v-if="total > 0">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchList"
            @current-change="fetchList"
            class="custom-pagination"
          />
        </div>
      </div>
    </div>

    <!-- 高级定制抽屉：研报详情 -->
    <el-drawer
      v-model="detailVisible"
      size="70%"
      title="AI 深度分析研报"
      destroy-on-close
      class="premium-drawer"
    >
      <div v-loading="detailLoading" class="detail-wrap custom-scrollbar">
        <template v-if="activeReport">
          
          <!-- 研报核心指标概览 -->
          <div class="detail-meta">
            <span class="meta-tag time">
              <el-icon><Clock /></el-icon> {{ formatDateTime(activeReport.created_at) }}
            </span>
            <span class="meta-tag repair-new">
              新增报修：{{ activeReport.repair_new_count }} 单
            </span>
            <span class="meta-tag repair-pending">
              未处理报修：{{ activeReport.repair_pending_count }} 单
            </span>
            <span class="meta-tag visitor">
              新增访客：{{ activeReport.visitor_new_count }} 人次
            </span>
            <span class="meta-tag money">
              物业收缴：¥{{ formatAmount(activeReport.property_paid_amount) }}
            </span>
          </div>

          <!-- Markdown 研报正文 -->
          <article class="markdown-body" v-html="activeReportHtml"></article>
          
        </template>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
// 引入必要的图标，包含 Loading 和 ArrowLeft
import { DataAnalysis, Plus, Refresh, Clock, ArrowLeft, Loading } from '@element-plus/icons-vue'
import Navbar from '@/components/layout/Navbar.vue'
import { generateAIReport, getAIReportDetail, getAIReportList } from '@/api/admin'

const listLoading = ref(false)
const generating = ref(false)
const detailLoading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(10)

const detailVisible = ref(false)
const activeReport = ref(null)
const activeReportHtml = computed(() => markdownToHtml(activeReport.value?.report || ''))

function formatDateTime(value) {
  return dayjs(value).format('YYYY-MM-DD HH:mm:ss')
}

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function isTimeoutError(error) {
  if (!error) return false
  return String(error.code || '') === 'ECONNABORTED' || /timeout/i.test(String(error.message || ''))
}

function isNetworkError(error) {
  if (!error) return false
  if (error.response) return false
  if (isTimeoutError(error)) return false
  return String(error.code || '') === 'ERR_NETWORK' || /Network Error/i.test(String(error.message || ''))
}

async function fetchList() {
  listLoading.value = true
  try {
    const res = await getAIReportList({
      page: page.value,
      size: size.value
    })
    list.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    if (isNetworkError(error)) {
      ElMessage.error('后端服务未连接，请先确认 smartcomunity 已启动')
      return
    }
    ElMessage.error(error.response?.data?.msg || error.message || '获取报表列表失败')
  } finally {
    listLoading.value = false
  }
}

async function handleGenerate(showMessage) {
  if (generating.value) return
  generating.value = true
  try {
    const data = await generateAIReport()
    if (showMessage) {
      ElMessage.success('AI 报告生成成功')
    }
    await fetchList()
    if (showMessage && data?.id) {
      await openDetail(data.id)
    }
  } catch (error) {
    if (isTimeoutError(error)) {
      ElMessage.warning('请求超时，后台可能仍在生成。请稍后点击“刷新列表”查看结果')
      return
    }
    if (isNetworkError(error)) {
      ElMessage.error('后端服务未连接，请先确认 smartcomunity 已启动')
      return
    }
    ElMessage.error(error.response?.data?.msg || error.message || 'AI 报告生成失败')
  } finally {
    generating.value = false
  }
}

async function openDetail(id) {
  detailVisible.value = true
  detailLoading.value = true
  try {
    activeReport.value = await getAIReportDetail(id)
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '获取报告详情失败')
  } finally {
    detailLoading.value = false
  }
}

onMounted(async () => {
  await fetchList()
})

function escapeHtml(text) {
  return String(text || '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function inlineMarkdown(text) {
  return text
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.+?)\*/g, '<em>$1</em>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
}

function normalizeMarkdownInput(markdown) {
  let text = String(markdown || '').replace(/\r\n/g, '\n').trim()
  if (!text) return ''

  if (text.startsWith('```')) {
    const lines = text.split('\n')
    lines.shift()
    while (lines.length && !lines[lines.length - 1].trim()) lines.pop()
    if (lines.length && lines[lines.length - 1].trim() === '```') {
      lines.pop()
    }
    text = lines.join('\n').trim()
  }

  const lines = text.split('\n')
  if (lines.length && /^(markdown|md)$/i.test(lines[0].trim())) {
    lines.shift()
    text = lines.join('\n').trim()
  }

  return text
}

function isTableHeader(line) { return /\|/.test(line) }
function isTableDivider(line) { return /^\|?[\s:-]+\|[\s|:-]*$/.test(line) }
function isTableRow(line) { return /\|/.test(line) }

function splitTableRow(line) {
  let content = line.trim()
  if (content.startsWith('|')) content = content.slice(1)
  if (content.endsWith('|')) content = content.slice(0, -1)
  return content.split('|').map((cell) => cell.trim())
}

function renderTable(headers, rows) {
  const head = `<thead><tr>${headers.map((h) => `<th>${inlineMarkdown(h)}</th>`).join('')}</tr></thead>`
  const body = `<tbody>${rows.map((row) => `<tr>${row.map((cell) => `<td>${inlineMarkdown(cell)}</td>`).join('')}</tr>`).join('')}</tbody>`
  return `<table>${head}${body}</table>`
}

function markdownToHtml(markdown) {
  const normalized = normalizeMarkdownInput(markdown)
  if (!normalized) return '<p class="empty-report">暂无内容分析</p>'

  const lines = escapeHtml(normalized).split('\n')
  const html = []
  let inUl = false
  let inOl = false

  const closeLists = () => {
    if (inUl) { html.push('</ul>'); inUl = false; }
    if (inOl) { html.push('</ol>'); inOl = false; }
  }

  for (let i = 0; i < lines.length; i += 1) {
    const line = lines[i].trim()
    if (!line) { closeLists(); continue; }

    if (isTableHeader(line) && i + 1 < lines.length && isTableDivider(lines[i + 1].trim())) {
      closeLists()
      const headers = splitTableRow(line)
      const rows = []
      i += 2
      for (; i < lines.length; i += 1) {
        const rowLine = lines[i].trim()
        if (!rowLine || !isTableRow(rowLine)) { i -= 1; break; }
        if (isTableDivider(rowLine)) continue
        rows.push(splitTableRow(rowLine))
      }
      html.push(renderTable(headers, rows))
      continue
    }

    if (/^###\s+/.test(line)) { closeLists(); html.push(`<h3>${inlineMarkdown(line.replace(/^###\s+/, ''))}</h3>`); continue; }
    if (/^##\s+/.test(line)) { closeLists(); html.push(`<h2>${inlineMarkdown(line.replace(/^##\s+/, ''))}</h2>`); continue; }
    if (/^#\s+/.test(line)) { closeLists(); html.push(`<h1>${inlineMarkdown(line.replace(/^#\s+/, ''))}</h1>`); continue; }
    if (/^>\s+/.test(line)) { closeLists(); html.push(`<blockquote>${inlineMarkdown(line.replace(/^>\s+/, ''))}</blockquote>`); continue; }
    
    if (/^[-*]\s+/.test(line)) {
      if (inOl) { html.push('</ol>'); inOl = false; }
      if (!inUl) { html.push('<ul>'); inUl = true; }
      html.push(`<li>${inlineMarkdown(line.replace(/^[-*]\s+/, ''))}</li>`)
      continue
    }
    
    if (/^\d+\.\s+/.test(line)) {
      if (inUl) { html.push('</ul>'); inUl = false; }
      if (!inOl) { html.push('<ol>'); inOl = true; }
      html.push(`<li>${inlineMarkdown(line.replace(/^\d+\.\s+/, ''))}</li>`)
      continue
    }
    
    if (/^---+$/.test(line)) { closeLists(); html.push('<hr />'); continue; }

    closeLists()
    html.push(`<p>${inlineMarkdown(line)}</p>`)
  }

  closeLists()
  return html.join('')
}
</script>

<style scoped>
/* 全局背景与布局 */
.ai-report-page {
  min-height: 100vh;
  background-color: #f4f7f9;
  padding-bottom: 80px;
}

.custom-container {
  max-width: 1280px;
  margin: 0 auto;
}

/* 顶部返回导航 */
.page-nav {
  padding: 24px 0 16px;
}

.back-btn {
  display: inline-flex;
  align-items: center;
  color: #606266;
  font-size: 15px;
  cursor: pointer;
  transition: color 0.3s;
  padding: 8px 16px 8px 0;
}

.back-btn:hover {
  color: #2d597b;
}

.back-icon {
  margin-right: 6px;
  font-size: 16px;
}

/* 高光标题 */
.page-header {
  padding: 16px 0 32px;
}

.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 32px;
  color: #2c3e50;
  font-weight: 700;
  margin: 0;
  z-index: 1;
}

.highlight-title::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: -5%;
  width: 110%;
  height: 14px;
  background-color: #2d597b; 
  opacity: 0.15;
  border-radius: 6px;
  z-index: -1;
}

/* 核心面板卡片 */
.premium-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 32px 40px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0,0,0,0.02);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px dashed #ebeef5;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 22px;
  font-weight: 700;
  color: #2c3e50;
}

.icon-box {
  width: 48px;
  height: 48px;
  background: #f0f7ff;
  color: #2d597b;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.header-actions {
  display: flex;
  gap: 16px;
}

/* 高级按钮定制 */
.action-btn {
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  border: none;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn-success { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); }
.btn-success:hover:not(:disabled) { background: #1f435d; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3); }

.btn-primary { background: #f0f7ff; color: #2d597b; border: 1px solid #cce3f6; }
.btn-primary:hover:not(:disabled) { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); transform: translateY(-2px); }

.action-btn:disabled { opacity: 0.7; cursor: not-allowed; }
.is-loading { animation: rotating 2s linear infinite; }

/* 深度定制表格 */
:deep(.custom-table) {
  --el-table-border-color: transparent;
  --el-table-header-bg-color: #fbfcfd;
  --el-table-header-text-color: #606266;
  border-radius: 12px;
}

:deep(.custom-table th.el-table__cell) { padding: 18px 0; font-weight: 600; border-bottom: 1px solid #ebeef5; }
:deep(.custom-table td.el-table__cell) { padding: 20px 0; border-bottom: 1px dashed #f0f2f5; }
:deep(.custom-table::before) { display: none; } /* 去掉底部实线 */

.time-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
  font-size: 14px;
}

.summary-text {
  color: #303133;
  font-size: 14px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.btn-text {
  background: transparent;
  border: none;
  color: #2d597b;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: color 0.2s;
}
.btn-text:hover { color: #1f435d; text-decoration: underline; }

/* 分页器 */
.pagination-wrap { display: flex; justify-content: flex-end; margin-top: 32px; padding-top: 16px; border-top: 1px solid #f0f2f5; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b; color: #fff; border-radius: 4px; }
:deep(.custom-pagination .el-pager li:hover) { color: #2d597b; }


/* ================= 抽屉与研报详情 (Markdown 高级排版) ================= */

/* 抽屉整体风格 */
:deep(.premium-drawer .el-drawer) {
  border-radius: 24px 0 0 24px;
  box-shadow: -10px 0 40px rgba(0,0,0,0.08);
}
:deep(.premium-drawer .el-drawer__header) {
  margin-bottom: 0; padding: 32px 40px; border-bottom: 1px solid #f0f2f5;
  font-size: 22px; font-weight: 700; color: #2c3e50;
}
:deep(.premium-drawer .el-drawer__body) { padding: 0; background: #ffffff; }

.detail-wrap {
  padding: 32px 40px 60px;
  height: 100%;
  overflow-y: auto;
}

/* 内部自定义滚动条 */
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #dcdfe6; border-radius: 3px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #c0c4cc; }

/* 指标概览区块 */
.detail-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  padding: 24px;
  background: #fbfcfd;
  border-radius: 12px;
  border: 1px dashed #dcdfe6;
  margin-bottom: 40px;
}

.meta-tag {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.meta-tag.time { background: #ffffff; color: #606266; border: 1px solid #ebeef5; }
.meta-tag.repair-new { background: #fdf6f6; color: #e4393c; }
.meta-tag.repair-pending { background: #fff7ed; color: #d97706; }
.meta-tag.visitor { background: #f4f4f5; color: #606266; }
.meta-tag.money { background: #f0fdf4; color: #166534; }

/* ★ 商业研报级 Markdown 渲染 ★ */
:deep(.markdown-body) {
  color: #303133;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
}

:deep(.markdown-body h1),
:deep(.markdown-body h2) {
  margin: 32px 0 16px;
  color: #2c3e50;
  font-weight: 700;
  border-bottom: 2px solid #f0f2f5;
  padding-bottom: 12px;
}

:deep(.markdown-body h3) {
  margin: 24px 0 12px;
  color: #2d597b; /* 强调蓝 */
  font-weight: 700;
  font-size: 18px;
}

:deep(.markdown-body p),
:deep(.markdown-body li) {
  line-height: 1.8;
  font-size: 15px;
  color: #475569;
  margin-bottom: 12px;
}

:deep(.markdown-body blockquote) {
  margin: 20px 0;
  padding: 16px 24px;
  border-left: 4px solid #2d597b;
  background: #f0f7ff;
  color: #606266;
  border-radius: 0 8px 8px 0;
  font-style: italic;
}

:deep(.markdown-body code) {
  background: #f4f4f5;
  color: #e4393c;
  padding: 2px 8px;
  border-radius: 4px;
  font-family: SFMono-Regular, Consolas, "Liberation Mono", Menlo, Courier, monospace;
}

/* 金融级报表表格样式 */
:deep(.markdown-body table) {
  width: 100%;
  border-collapse: collapse;
  margin: 24px 0;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 0 0 1px #ebeef5; /* 外边框 */
}

:deep(.markdown-body th) {
  background: #fbfcfd;
  font-weight: 600;
  color: #303133;
  text-align: left;
}

:deep(.markdown-body td), 
:deep(.markdown-body th) {
  border: none;
  border-bottom: 1px solid #ebeef5; /* 内部横线 */
  padding: 16px 20px;
}

:deep(.markdown-body tr:last-child td) {
  border-bottom: none;
}

.empty-report {
  text-align: center;
  color: #a4b0be;
  padding: 60px 0;
  font-size: 16px;
}

/* 响应式 */
@media (max-width: 960px) {
  .panel-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
  }
  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>