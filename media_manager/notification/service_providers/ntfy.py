import requests

from media_manager.notification.config import NotificationConfig
from media_manager.notification.schemas import MessageNotification
from media_manager.notification.service_providers.abstractNotificationServiceProvider import (
    AbstractNotificationServiceProvider,
)


class NtfyNotificationServiceProvider(AbstractNotificationServiceProvider):
    """
    Ntfy Notification Service Provider
    """

    def __init__(self):
        self.config = NotificationConfig()

    def send_notification(self, message: MessageNotification) -> bool:
        response = requests.post(
            url=self.config.ntfy_url,
            data=message.message.encode(encoding="utf-8"),
            headers={
                "Title": "MediaManager - " + message.title,
            },
        )
        if response.status_code not in range(200, 300):
            return False
        return True
