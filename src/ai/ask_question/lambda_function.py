# A function that receives a topic from the user and returns a question from the topic using openai api model gpt4o-mini
from openai import OpenAI
import json
import os


def ask_question(topic):

    client = OpenAI(api_key=os.environ.get("OPENAI_API_KEY"))

    chat_completion = client.chat.completions.create(
        messages=[
            {
                "role": "user",
                "content": "Give me a question about " + topic + ". No context",
            }
        ],
        model="gpt-4o-mini",
    )

    response = chat_completion.choices[0].message.content

    return response


def lambda_handler(event, context):
    # Extract topic from query parameters
    topic = event['queryStringParameters']['topic']
    question = ask_question(topic)
    return {
        'statusCode': 200,
        'body': json.dumps({'question': question})
    }

# if __name__ == "__main__":
#     topic = input("Enter a topic: ")
#     print(ask_question(topic))
