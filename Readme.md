# QRiosity

QRiosity is an AI-powered learning platform designed to deliver personalized learning experiences, daily challenges, and curated roadmaps for learners. It leverages advanced AI models (OpenAI, LLMs) to generate questions, assess answers, and recommend resources tailored to user-selected topics. 

## Features

- **Personalized Learning Roadmaps**: Generate detailed learning paths for any topic, offering curated courses, resources, and skill progression.
- **Daily Challenges**: Receive AI-generated questions tailored to your interests and topics, encouraging daily practice and engagement.
- **AI-Powered Feedback**: Answers are rated by an AI system, providing instant feedback, insights, and motivation to help you improve.
- **Course Creation**: Users can create and add their own courses, contributing to a growing library of resources.
- **Modern Web Interface**: Built with TypeScript and React for a seamless user experience.

## Technologies Used

- **Frontend**: React, TypeScript
- **Backend**: Go (Golang)
- **AI Integration**: Python Lambda functions using OpenAI GPT-4o-mini
- **APIs & Microservices**: RESTful communication between backend and AI services
- **Cloud Functions**: AWS Lambda for AI-powered question and answer evaluation

## Project Structure

```
src/
  ai/              # Python Lambda functions for AI capabilities
    ask_question/
    rate_question/
    get-roadmap/
  backend/         # Go backend for API and business logic
    cmd/
    internal/
  frontend/        # React + TypeScript frontend application
```

## Getting Started

1. **Clone the repository**  
   `git clone https://github.com/hmrguez/QRiosity.git`

2. **Install dependencies**  
   - For the frontend, run `npm install` in `src/frontend`.
   - For the backend, ensure Go is installed and set up required env variables.
   - For AI services, set up Python environment and configure OpenAI keys as environment variables.

3. **Run the Application**  
   - Start backend and frontend servers.
   - Deploy Lambda functions or run locally for development.

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

This project is currently not licensed. See repository info for updates.

---
**Repository:** https://github.com/hmrguez/QRiosity
