import requests
from pydantic_settings import BaseSettings

from media_manager.config import AllEncompassingConfig
from media_manager.notification.schemas import MessageNotification
from media_manager.notification.service_providers.abstractNotificationServiceProvider import (
    AbstractNotificationServiceProvider,
)

class GotifyConfig(BaseSettings):
    enabled: bool = False
    api_key: str | None = None
    url: str | None = (
        None  # e.g. https://gotify.example.com (note lack of trailing slash)
    )

class GotifyNotificationServiceProvider(AbstractNotificationServiceProvider):
    """
    Gotify Notification Service Provider
    """

    def __init__(self):
        self.config = AllEncompassingConfig().notifications.gotify

    def send_notification(self, message: MessageNotification) -> bool:
        response = requests.post(
            url=f"{self.config.url}/message?token={self.config.api_key}",
            json={
                "message": message.message,
                "title": message.title,
            },
        )
        if response.status_code not in range(200, 300):
            return False
        return True
