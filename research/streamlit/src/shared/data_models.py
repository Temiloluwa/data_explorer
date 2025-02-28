from datetime import datetime
from typing import Optional, Dict
from pydantic import BaseModel
import pytz

class Message(BaseModel):
    """A class representing a message in a conversation."""
    role: str = "User"
    content: str = ""
    timestamp: str = None

    def __init__(self, **data):
        """Initialize a Message object."""
        if data.get('timestamp') is None:
            berlin_timezone = pytz.timezone('Europe/Berlin')
            berlin_time = datetime.now(berlin_timezone)
            data['timestamp'] = berlin_time.strftime("%Y-%m-%d %H:%M:%S")
        
        super().__init__(**data)