import json
from openai import OpenAI
from pydantic import BaseModel
from typing import List

client = OpenAI()


class Course(BaseModel):
    title: str
    url: str
    description: str
    source: str
    difficulty: str
    topics: List[str]
    is_free: bool
    duration: int
    language: str


class Roadmap(BaseModel):
    title: str
    description: str
    courses: List[Course]
    topics: List[str]
    difficulty: str


response_format = {
    "type": "json_schema",
    "json_schema": {
        "name": "learning_roadmap",
        "strict": True,
        "schema": {
            "type": "object",
            "additionalProperties": False,
            "properties": {
                "title": {"type": "string"},
                "description": {"type": "string"},
                "courses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/course"
                    }
                },
                "topics": {
                    "type": "array",
                    "items": {"type": "string"}
                },
                "difficulty": {"type": "string"}
            },
            "required": ["title", "description", "courses", "topics", "difficulty"],
            "definitions": {
                "course": {
                    "type": "object",
                    "additionalProperties": False,
                    "properties": {
                        "title": {"type": "string"},
                        "url": {"type": "string"},
                        "description": {"type": "string"},
                        "source": {"type": "string"},
                        "difficulty": {"type": "string"},
                        "topics": {
                            "type": "array",
                            "items": {"type": "string"}
                        },
                        "is_free": {"type": "boolean"},
                        "duration": {"type": "integer"},
                        "language": {"type": "string"}
                    },
                    "required": [
                        "title",
                        "url",
                        "description",
                        "source",
                        "difficulty",
                        "topics",
                        "is_free",
                        "duration",
                        "language"
                    ]
                }
            }
        }
    }
}


def get_roadmap(topic: str) -> str:
    chat_completion = client.chat.completions.create(
        messages=[
            {
                "role": "user",
                "content": f"Create a roadmap for {topic}",
            }
        ],
        model="gpt-4o-mini",
        response_format=response_format,
    )

    response = chat_completion.choices[0].message.content

    print(response)

    roadmap_data = json.loads(response)
    roadmap = Roadmap(**roadmap_data)
    return roadmap.json()


def lambda_handler(event, context):
    try:
        topic = event['queryStringParameters'].get('topic')
        if not topic:
            return {
                "statusCode": 400,
                "body": json.dumps({"error": "Topic is required"})
            }

        roadmap_json = get_roadmap(topic)

        return {
            "statusCode": 200,
            "body": roadmap_json
        }
    except Exception as e:
        return {
            "statusCode": 500,
            "body": json.dumps({"error": str(e)})
        }
