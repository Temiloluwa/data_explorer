import time
from functools import wraps
from typing import Any, Callable

from utilities.logger import logger


def timed_async(func: Callable) -> Callable:
    @wraps(func)
    async def wrapper(*args, **kwargs) -> Any:
        start_time = time.monotonic()
        result = await func(*args, **kwargs)
        duration = time.monotonic() - start_time
        logger.info(f"{func.__name__} executed in {duration:.2f}s")
        return result

    return wrapper
