import requests
from pydantic import HttpUrl

from media_manager.notification.config import NotificationConfig
from media_manager.notification.schemas import MessageNotification
from media_manager.notification.service_providers.abstractNotificationServiceProvider import \
    AbstractNotificationServiceProvider


class GotifyNotificationServiceProvider(AbstractNotificationServiceProvider):
    """
    Gotify Notification Service Provider
    """

    def __init__(self):
        self.config = NotificationConfig()

    def send_notification(self, message: MessageNotification) -> bool:
        response = requests.post(
            url=f"{self.config.gotify_url}/message?token={self.config.gotify_api_key}",
            json={
                "message": message.message,
                "title": message.title,
            },
        )
        if response.status_code not in range(200,300):
            return False
        return True
