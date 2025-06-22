class MediaAlreadyExists(ValueError):
    """Raised when a show already exists"""

    pass


class NotFoundError(Exception):
    """Custom exception for when an entity is not found."""

    pass


class InvalidConfigError(Exception):
    """Custom exception for when an entity is not found."""

    pass
