from dataclasses import dataclass

from .config import settings


@dataclass(frozen=True)
class ServiceClients:
    user_service: str = settings.user_service_url
    mall_service: str = settings.mall_service_url
    community_service: str = settings.community_service_url
    workorder_service: str = settings.workorder_service_url
    statistics_service: str = settings.statistics_service_url

    def endpoints(self) -> dict:
        return {
            "user_service": self.user_service,
            "mall_service": self.mall_service,
            "community_service": self.community_service,
            "workorder_service": self.workorder_service,
            "statistics_service": self.statistics_service,
        }


service_clients = ServiceClients()
