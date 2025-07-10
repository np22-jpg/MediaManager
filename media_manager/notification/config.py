from pydantic_settings import BaseSettings

from media_manager.notification.service_providers.email import EmailNotificationsConfig
from media_manager.notification.service_providers.gotify import GotifyConfig
from media_manager.notification.service_providers.ntfy import NtfyConfig
from media_manager.notification.service_providers.pushover import PushoverConfig


class EmailConfig(BaseSettings):
    smtp_host: str = ""
    smtp_port: int = 587
    smtp_user: str = ""
    smtp_password: str = ""
    from_email: str = ""
    use_tls: bool = False


class NotificationConfig(BaseSettings):
    smtp_config: EmailConfig = EmailConfig()
    email_notifications: EmailNotificationsConfig = EmailNotificationsConfig()
    gotify: GotifyConfig = GotifyConfig()
    ntfy: NtfyConfig = NtfyConfig()
    pushover: PushoverConfig = PushoverConfig()
