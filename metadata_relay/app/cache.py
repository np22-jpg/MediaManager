import json
import hashlib
import logging
from typing import Any, Optional
import valkey.asyncio as valkey
import os
from functools import wraps

log = logging.getLogger(__name__)

redis_client = valkey.Redis(
    host=os.getenv("VALKEY_HOST", "localhost"),
    port=int(os.getenv("VALKEY_PORT", 6379)),
    db=int(os.getenv("VALKEY_DB", 0)),
    decode_responses=True,
)


def generate_cache_key(prefix: str, *args, **kwargs) -> str:
    key_data = f"{prefix}:{args}:{sorted(kwargs.items())}"
    return hashlib.md5(key_data.encode()).hexdigest()


async def get_cached_response(cache_key: str) -> Optional[Any]:
    try:
        cached_data = await redis_client.get(cache_key)
        if cached_data:
            return json.loads(cached_data)
        return None
    except Exception as e:
        log.error(f"Error getting cached response: {e}")
        return None


async def set_cached_response(cache_key: str, data: Any, ttl: int = 3600) -> bool:
    try:
        await redis_client.setex(cache_key, ttl, json.dumps(data, default=str))
        return True
    except Exception as e:
        log.error(f"Error setting cached response: {e}")
        return False


def cache_response(prefix: str, ttl: int = 3600):
    def decorator(func):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            cache_key = generate_cache_key(prefix, *args, **kwargs)

            cached_response = await get_cached_response(cache_key)
            if cached_response is not None:
                log.info(f"Cache hit for key: {cache_key}")
                return cached_response

            log.info(f"Cache miss for key: {cache_key}")

            response = await func(*args, **kwargs)

            await set_cached_response(cache_key, response, ttl)

            return response

        return wrapper

    return decorator
