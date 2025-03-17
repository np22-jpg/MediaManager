import json
import logging
from collections import Counter
from typing import List

from ollama import ChatResponse
from ollama import chat
from pydantic import BaseModel

log = logging.getLogger(__name__)


class NFO(BaseModel):
    season: int



def get_season(nfo: str) -> int | None:
    responses: List[ChatResponse] = []
    parsed_responses: List[int] = []

    for i in range(0, 4):
        responses.append(chat(
            model='qwen2.5:0.5b',
            format=NFO.model_json_schema(),
            messages=[
                {
                    'role': 'user',
                    'content':
                        "which season does a torrent with the following NFO contain, output a season number, the season number is an integer? output in json please" +
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
