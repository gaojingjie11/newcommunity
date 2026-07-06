package config

import (
	"os"
	"strings"

	"smartcommunity-microservices/common/db"

	"github.com/zeromicro/go-zero/zrpc"
)

type LLMModelDetail struct {
	ApiKey  string `json:"ApiKey,optional" yaml:"ApiKey"`
	BaseUrl string `json:"BaseUrl,optional" yaml:"BaseUrl"`
	Model   string `json:"Model,optional" yaml:"Model"`
}

type ModelsConfig struct {
	ChatDefault    LLMModelDetail `json:",optional"`
	AgentReasoning LLMModelDetail `json:",optional"`
	Embedding      LLMModelDetail `json:",optional"`
	ImageVision    LLMModelDetail `json:",optional"`
	ASR            LLMModelDetail `json:",optional"`
	TTS            LLMModelDetail `json:",optional"`
}

type AgentConfig struct {
	LlmApiKey              string       `json:",optional"`
	LlmBaseUrl             string       `json:",optional"`
	LlmModel               string       `json:",optional"`
	RAGMaxResults          int          `json:",optional"`
	RAGSyncIntervalSeconds int          `json:",optional"`
	Models                 ModelsConfig `json:",optional"`
}

func isPlaceholder(val string) bool {
	val = strings.TrimSpace(val)
	return val == "" || strings.Contains(val, "${") || (strings.HasPrefix(val, "$") && len(val) > 1)
}

// GetModelConfig resolves model credentials, falling back to global settings if empty or containing unreplaced placeholders
func (c AgentConfig) GetModelConfig(detail LLMModelDetail) (apiKey, baseUrl, modelName string) {
	apiKey = os.ExpandEnv(detail.ApiKey)
	if isPlaceholder(apiKey) {
		apiKey = os.ExpandEnv(c.LlmApiKey)
	}
	if isPlaceholder(apiKey) {
		apiKey = ""
	}

	baseUrl = os.ExpandEnv(detail.BaseUrl)
	if isPlaceholder(baseUrl) {
		baseUrl = os.ExpandEnv(c.LlmBaseUrl)
	}
	if isPlaceholder(baseUrl) {
		baseUrl = ""
	}

	modelName = os.ExpandEnv(detail.Model)
	if isPlaceholder(modelName) {
		modelName = os.ExpandEnv(c.LlmModel)
	}
	if isPlaceholder(modelName) {
		modelName = ""
	}

	return apiKey, baseUrl, modelName
}

type Config struct {
	zrpc.RpcServerConf
	Postgres     db.PostgresConfig
	UserRpc      zrpc.RpcClientConf
	MallRpc      zrpc.RpcClientConf
	CommunityRpc zrpc.RpcClientConf
	WorkorderRpc zrpc.RpcClientConf
	Agent        AgentConfig `json:",optional"`
}
