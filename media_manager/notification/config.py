from pydantic_settings import BaseSettings, SettingsConfigDict


class EmailConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="EMAIL_")
    smtp_host: str
    smtp_port: int
    smtp_user: str
    smtp_password: str
    from_email: str
    use_tls: bool = False

class NotificationConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="NOTIFICATION_")

    email: str|None = None # the email address to send notifications to

    ntfy_url: str|None = None # e.g. https://ntfy.sh/your-topic (note lack of trailing slash)

    pushover_api_key : str|None = None
    pushover_user: str|None = None

    gotify_api_key: str|None = None
    gotify_url: str|None = None # e.g. https://gotify.example.com (note lack of trailing slash)