# Task
Fast Track Code Test Quiz - Instructions

Preferred Components:
- REST API or gRPC
- CLI that talks with the API, preferably using https://github.com/spf13/cobra ;( as CLI framework )

User stories/Use cases:
- User should be able to get questions with a number of answers
- User should be able to select just one answer per question.
- User should be able to answer all the questions and then post his/hers answers and get back how many correct 
  answers they had, displayed to the user.
- User should see how well they compared to others who have taken the quiz, eg. "You were better than 60% of all 
  quizzers"

How it should be delivered?
After you complete it, please, share with us the GitHub link to the public repository so they can access it.

# Solution
Running the API (in 1st session):
```bash 
go run ./cmd/api/main.go
```

Running the client (in 2nd session - will work only if server is running):
```bash
go run ./cmd/client/main.go
```