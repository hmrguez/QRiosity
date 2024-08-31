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

        #         roadmap_json = get_roadmap(topic)
        roadmap = {
            "title": "Web Development Learning Roadmap",
            "description": "This roadmap outlines the essential skills, topics, and courses needed to become a proficient web developer, covering front-end and back-end technologies.",
            "courses": [
                {
                    "title": "The Complete Web Developer Bootcamp",
                    "url": "https://www.udemy.com/course/the-complete-web-developer-bootcamp/",
                    "description": "A comprehensive course covering HTML, CSS, JavaScript, Node.js, and more to build web applications.",
                    "source": "Udemy",
                    "difficulty": "Beginner",
                    "topics": [
                        "HTML",
                        "CSS",
                        "JavaScript",
                        "Node.js",
                        "Bootstrap"
                    ],
                    "is_free": False,
                    "duration": 60,
                    "language": "English"
                },
                {
                    "title": "JavaScript: The Good Parts",
                    "url": "https://www.oreilly.com/library/view/javascript-the-good/9780596805524/",
                    "description": "This book offers insights into the best practices and features of JavaScript.",
                    "source": "O'Reilly",
                    "difficulty": "Intermediate",
                    "topics": [
                        "JavaScript"
                    ],
                    "is_free": False,
                    "duration": 10,
                    "language": "English"
                },
                {
                    "title": "Frontend Development with React",
                    "url": "https://www.codecademy.com/learn/react-101",
                    "description": "An interactive course to learn React.js, a popular JavaScript library for building user interfaces.",
                    "source": "Codecademy",
                    "difficulty": "Intermediate",
                    "topics": [
                        "React.js",
                        "JavaScript",
                        "Frontend"
                    ],
                    "is_free": True,
                    "duration": 30,
                    "language": "English"
                },
                {
                    "title": "Node.js, Express, MongoDB & More: The Complete Bootcamp 2023",
                    "url": "https://www.udemy.com/course/nodejs-express-mongodb-bootcamp/",
                    "description": "Learn back-end development using Node, Express, and MongoDB to create robust server-side applications.",
                    "source": "Udemy",
                    "difficulty": "Intermediate",
                    "topics": [
                        "Node.js",
                        "Express.js",
                        "MongoDB"
                    ],
                    "is_free": False,
                    "duration": 60,
                    "language": "English"
                },
                {
                    "title": "Full-Stack Web Development with React Specialization",
                    "url": "https://www.coursera.org/specializations/full-stack-react",
                    "description": "A series of courses that provide an overview of full-stack development with a focus on React, Node.js, and database management.",
                    "source": "Coursera",
                    "difficulty": "Advanced",
                    "topics": [
                        "React.js",
                        "Node.js",
                        "Express.js",
                        "MongoDB"
                    ],
                    "is_free": True,
                    "duration": 90,
                    "language": "English"
                }
            ],
            "topics": [
                "HTML",
                "CSS",
                "JavaScript",
                "Responsive Design",
                "Version Control (Git)",
                "Web Accessibility",
                "Frontend Frameworks (e.g., React, Angular, Vue)",
                "Backend Development (Node.js, Express)",
                "Database Management (SQL, MongoDB)",
                "Deployment & Hosting (Heroku, Netlify)"
            ],
            "difficulty": "Beginner to Advanced"
        }

        return {
            "statusCode": 200,
            "body": json.dumps(roadmap)
        }
    except Exception as e:
        return {
            "statusCode": 500,
            "body": json.dumps({"error": str(e)})
        }
