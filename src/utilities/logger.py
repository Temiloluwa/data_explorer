import logging
from logging.config import dictConfig

from config import settings

logging_config = {
    "version": 1,
    "formatters": {
        "default": {"format": "%(asctime)s - %(name)s - %(levelname)s - %(message)s"}
    },
    "handlers": {
        "console": {
            "class": "logging.StreamHandler",
            "formatter": "default",
        }
    },
    "root": {"level": "INFO", "handlers": ["console"]},
}

dictConfig(logging_config)
logger = logging.getLogger(__name__)
