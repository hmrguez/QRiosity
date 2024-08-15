import json
import os

from pydantic import BaseModel
from openai import OpenAI


class RateQuestionRequest(BaseModel):
    question: str
    answer: str


class RateQuestionResponse(BaseModel):
    rating: int
    insight: str

client = OpenAI(api_key=os.environ.get("OPENAI_API_KEY"))


def rate_question(request: RateQuestionRequest) -> str:
    completion = client.chat.completions.create(
        model="gpt-4o-mini",
        messages=[
            {"role": "system",
             "content": "Return a json object with rating and insight of the following question and answer"},
            {"role": "user", "content": "Question: " + request.question + "\nAnswer: " + request.answer},
        ],
    )

    return completion.choices[0].message.content


def lambda_handler(event, context):
    try:
        body = json.loads(event['body'])
        request = RateQuestionRequest(**body)
        response = rate_question(request)
        return {
            'statusCode': 200,
            'body': response
        }
    except Exception as e:
        return {
            'statusCode': 400,
            'body': json.dumps({'error': str(e)})
        }
