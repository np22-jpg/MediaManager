import requests

from media_manager.notification.config import NotificationConfig
from media_manager.notification.schemas import MessageNotification
from media_manager.notification.service_providers.abstractNotificationServiceProvider import \
    AbstractNotificationServiceProvider


class PushoverNotificationServiceProvider(AbstractNotificationServiceProvider):
    def __init__(self):
        self.config = NotificationConfig()


    def send_notification(self, message: MessageNotification) -> bool:
        response = requests.post(
            url = "https://api.pushover.net/1/messages.json",
            params={
                "token": self.config.pushover_api_key,
                "user": self.config.pushover_user,
                "message": message.message,
                "title": "MediaManager - "+ message.title,
            }
        )
        if response.status_code not in range(200,300):
            return False
        return True
