import json
import logging
from collections import Counter
from typing import List

from ollama import ChatResponse, chat
from pydantic import BaseModel

import config


class NFO(BaseModel):
    season: int


def get_season(nfo: str) -> int | None:
    responses: List[ChatResponse] = []
    parsed_responses: List[int] = []

    for i in range(0, 5):
        responses.append(chat(
            model=config.model_name,
            format=NFO.model_json_schema(),
            messages=[
                {
                    'role': 'user',
                    'content':
                        "Tell me which season the torrent with this description contains?" +
                        " output a season number in json format, the season number is an integer" +
                        nfo
                },
            ]))

    for response in responses:
        season_number: int
        try:
            season_number: int = json.loads(response.message.content)['season']
        except Exception as e:
            log.warning(f"failed to parse season number: {e}")
            break
        parsed_responses.append(season_number)

    most_common = Counter(parsed_responses).most_common(1)
    log.debug(f"extracted season number: {most_common} from nfo: {nfo}")
    return most_common[0][0]


config = config.MachineLearningConfig()
log = logging.getLogger(__name__)
