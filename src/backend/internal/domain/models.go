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

type Topic struct {
	Name string `json:"name"`
}

type User struct {
	Name                    string   `json:"name"`
	Role                    int      `json:"role"`
	Email                   string   `json:"email"`
	Topics                  []string `json:"topics,omitempty"`
	DailyChallengeAvailable bool     `json:"dailyChallengeAvailable"`
	Username                string   `json:"username"`
	Roadmaps                []string `json:"roadmapsIds"`
}

type Course struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Source      string   `json:"source"`
	Difficulty  string   `json:"difficulty"`
	Topics      []string `json:"topics"`
	IsFree      bool     `json:"isFree"`
	Author      string   `json:"author"`
	Duration    int      `json:"duration"`
	Language    string   `json:"language"`
}

type Roadmap struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Courses    []Course `json:"courses"`
	CourseIDs  []string `json:"courseIDs"`
	Topics     []string `json:"topics"`
	IsCustom   bool     `json:"isCustom"`
	CreatedBy  string   `json:"createdBy"`
	Likes      int      `json:"likes"`
	Difficulty string   `json:"difficulty"`
	Liked      bool     `json:"liked,omitempty"`
}
