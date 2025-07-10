import requests
from pydantic_settings import BaseSettings

from media_manager.config import AllEncompassingConfig
from media_manager.notification.schemas import MessageNotification
from media_manager.notification.service_providers.abstractNotificationServiceProvider import (
    AbstractNotificationServiceProvider,
)


class NtfyConfig(BaseSettings):
    enabled: bool = False
    url: str | None = (
        None  # e.g. https://ntfy.sh/your-topic (note lack of trailing slash)
    )


class NtfyNotificationServiceProvider(AbstractNotificationServiceProvider):
    """
    Ntfy Notification Service Provider
    """

    def __init__(self):
        self.config = AllEncompassingConfig().notifications.ntfy

    def send_notification(self, message: MessageNotification) -> bool:
        response = requests.post(
            url=self.config.url,
            data=message.message.encode(encoding="utf-8"),
            headers={
                "Title": "MediaManager - " + message.title,
            },
        )
        if response.status_code not in range(200, 300):
            return False
        return True
