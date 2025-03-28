import json
import logging
from collections import Counter
from typing import List

from ollama import ChatResponse, chat
from pydantic import BaseModel

from ml.config import MachineLearningConfig


class NFO(BaseModel):
    season: int


class Contains(BaseModel):
    contains: bool

def get_season(nfo: str) -> int | None:
    responses: List[ChatResponse] = []
    parsed_responses: List[int] = []

    for i in range(0, 5):
        responses.append(chat(
            model=config.ollama_model_name,
            format=NFO.model_json_schema(),
            messages=[
                {
                    'role': 'USER',
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


def contains_season(season_number: int, string_to_analyze: str) -> bool:
    responses: List[ChatResponse] = []
    parsed_responses: List[bool] = []

    for i in range(0, 3):
        responses.append(chat(
            model=config.ollama_model_name,
            format=Contains.model_json_schema(),
            messages=[
                {
                    'role': 'USER',
                    'content':
                        "Does this torrent contain the season " + season_number.__str__() + " ?" +
                        " output a boolean json format" +
                        string_to_analyze
                },
            ]))

    for response in responses:
        try:
            answer: bool = json.loads(response.message.content)['contains']
            log.debug(f"extracted contains: {answer}")
        except Exception as e:
            log.warning(f"failed to parse season number: {e}")
            break
        parsed_responses.append(answer)

    most_common = Counter(parsed_responses).most_common(1)
    log.debug(f"according to AI {string_to_analyze} contains season {season_number} {most_common[0][0]}")
    return most_common[0][0]


config = MachineLearningConfig
log = logging.getLogger(__name__)
