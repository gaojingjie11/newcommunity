import os
from dataclasses import dataclass


@dataclass(frozen=True)
class Settings:
    service_name: str = os.getenv("SERVICE_NAME", "agent-service")
    llm_api_key: str = os.getenv("LLM_API_KEY", "")
    llm_base_url: str = os.getenv("LLM_BASE_URL", "")
    llm_model: str = os.getenv("LLM_MODEL", "")
    user_service_url: str = os.getenv("USER_SERVICE_URL", "http://user-service:8001")
    mall_service_url: str = os.getenv("MALL_SERVICE_URL", "http://mall-service:8002")
    community_service_url: str = os.getenv("COMMUNITY_SERVICE_URL", "http://community-service:8003")
    workorder_service_url: str = os.getenv("WORKORDER_SERVICE_URL", "http://community-service:8003")
    statistics_service_url: str = os.getenv("STATISTICS_SERVICE_URL", "http://community-service:8003")


settings = Settings()
