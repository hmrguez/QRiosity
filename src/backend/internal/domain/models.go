package domain

type Content interface {
	IsContent()
}

type AuthPayload struct {
	Token string `json:"token"`
}

type ChallengeResponse struct {
	UserID   string `json:"userId"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Rating   int    `json:"rating"`
	Insight  string `json:"insight"`
}

type ContentInput struct {
	Quiz   *QuizInput   `json:"quiz,omitempty"`
	Lesson *LessonInput `json:"lesson,omitempty"`
}

type Course struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Author  string    `json:"author"`
	Content []Content `json:"content"`
}

type CourseInput struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Author  string          `json:"author"`
	Content []*ContentInput `json:"content"`
}

type Lesson struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Author string `json:"author"`
}

func (Lesson) IsContent() {}

type LessonInput struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Author string `json:"author"`
}

type Mutation struct {
}

type Problem struct {
	Question   string   `json:"question"`
	Categories []string `json:"categories"`
	Type       string   `json:"type"`
}

type ProblemInput struct {
	Question   string   `json:"question"`
	Categories []string `json:"categories"`
	Type       string   `json:"type"`
}

type Query struct {
}

type Quiz struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (Quiz) IsContent() {}

type QuizInput struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Roadmap struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Author  string    `json:"author"`
	Courses []*Course `json:"courses"`
}

type RoadmapInput struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Author  string         `json:"author"`
	Courses []*CourseInput `json:"courses"`
}

type Topic struct {
	Name string `json:"name"`
}

type User struct {
	Name                    string   `json:"name"`
	Email                   string   `json:"email"`
	Topics                  []string `json:"topics,omitempty"`
	DailyChallengeAvailable bool     `json:"dailyChallengeAvailable"`
}
