from typing import Any

from fastapi import FastAPI
from pydantic import BaseModel, Field

from .clients import service_clients
from .config import settings


app = FastAPI(title="Smart Community Agent Service", version="0.1.0")


class ChatRequest(BaseModel):
    user_id: int | None = None
    message: str = Field(default="", description="User message")


class RepairClassifyRequest(BaseModel):
    content: str = ""
    images: list[str] = []


class ComplaintRiskRequest(BaseModel):
    content: str = ""
    user_id: int | None = None


class RecommendRequest(BaseModel):
    user_id: int | None = None
    scene: str = "home"
    context: dict[str, Any] = {}


def success(data: dict) -> dict:
    return {"code": 0, "message": "success", "data": data}


@app.get("/health")
def health() -> dict:
    return success({"service": settings.service_name, "status": "ok"})


@app.get("/agent/health")
def agent_health() -> dict:
    return health()


@app.post("/agent/chat")
def chat(req: ChatRequest) -> dict:
    return success(
        {
            "intent": "general_community_consultation",
            "reply": "您好，我是社区客服 Agent 占位服务。当前阶段尚未接入真实大模型和业务服务。",
            "called_services": [],
            "service_clients": service_clients.endpoints(),
            "input_message": req.message,
        }
    )


@app.post("/agent/repair-classify")
def repair_classify(req: RepairClassifyRequest) -> dict:
    return success(
        {
            "category": "unknown",
            "urgency": "normal",
            "suggested_department": "property-maintenance",
            "reason": "占位分类结果；后续将接入 LLM、历史工单和规则库。",
            "image_count": len(req.images),
        }
    )


@app.post("/agent/complaint-risk")
def complaint_risk(req: ComplaintRiskRequest) -> dict:
    return success(
        {
            "category": "general",
            "risk_level": "low",
            "suggested_action": "create_follow_up_ticket",
            "reason": "占位风险评估；后续将结合投诉内容、用户历史和人工规则。",
        }
    )


@app.post("/agent/recommend")
def recommend(req: RecommendRequest) -> dict:
    return success(
        {
            "recommended_items": [
                {"type": "service", "id": "notice", "name": "公告通知"},
                {"type": "service", "id": "repair", "name": "报修服务"},
            ],
            "reason": "占位推荐结果；后续将接入用户画像、商品服务数据和 LLM。",
            "scene": req.scene,
        }
    )
