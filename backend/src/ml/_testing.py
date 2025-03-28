import json
from datetime import datetime, timedelta

from ollama import ChatResponse, chat
from pydantic import BaseModel


class NFO(BaseModel):
    season: int


# or access fields directly from the response object
start_time = datetime.now() + timedelta(seconds=300)
i = 0
failed_prompts = 0
while start_time > datetime.now():
    response: ChatResponse = chat(model='qwen2.5:0.5b',
                                  format=NFO.model_json_schema()
                                  , messages=[
            {
                'role': 'USER',
                'content':
                    "which season does a torrent with the following NFO contain? output the season number, which is an integer in json please\n" +
                    "The.Big.Bang.Theory.(2007).Season.9.S09.(1080p.BluRay.x265.HEVC.10bit.AAC.5.1.Vyndros)"
            },
        ])
    i += 1
    print("prompt #", i)
    print("remaining time: ", start_time - datetime.now())
    try:
        json2 = json.loads(response.message.content)
        print(json2)
    except Exception as e:
        print("prompt failed", e)
        print(response.message.content)
        failed_prompts += 1

    if json2['season'] != 9:
        failed_prompts += 1

print("prompts: ", i, " total time: 120s")
print("failed prompts: ", failed_prompts)
print("average time per prompt: ", 300 / i)
print("average time per successful prompt: ", 300 / (i - failed_prompts))
print("ratio successful/failed prompts: ", failed_prompts / (i - failed_prompts))
