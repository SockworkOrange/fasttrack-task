package api

import (
	"context"
	"encoding/json"
	"fasttrack-task/pkg/gen"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sort"
	"sync"
)

type Quiz struct {
	Questions []Question `json:"questions"`
}

var submissions []int

type Question struct {
	Text          string   `json:"text"`
	CorrectIdx    int      `json:"correct_idx"`
	Possibilities []string `json:"possibilities"`
}

var quiz *Quiz

func init() {
	loadQuiz()
}

func loadQuiz() {
	file, err := os.Open("default.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatalf("Failed to close file: %v", err)
		}
	}(file)

	if err = json.NewDecoder(file).Decode(&quiz); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	validateQuiz()
}

func validateQuiz() {
	for _, question := range quiz.Questions {
		if question.CorrectIdx >= len(question.Possibilities) {
			log.Fatalf(fmt.Sprintf("Correct answer id %d missing from possibilities", question.CorrectIdx))
		}
	}
}

func Run() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	gen.RegisterQuizServiceServer(grpcServer, newServer())
	log.Printf("Serving requests on %s...\n", lis.Addr().String())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type quizServiceServer struct {
	gen.UnimplementedQuizServiceServer
	mu sync.Mutex
}

func newServer() *quizServiceServer {
	s := &quizServiceServer{}
	return s
}

func (s *quizServiceServer) Get(_ context.Context, _ *gen.GetRequest) (*gen.Quiz, error) {
	questions := make([]*gen.Question, 0)
	for questionId, question := range quiz.Questions {
		questions = append(questions, &gen.Question{
			QuestionId:    int32(questionId),
			Text:          question.Text,
			Possibilities: question.Possibilities,
		})
	}
	response := &gen.Quiz{Questions: questions}
	return response, nil
}

func (s *quizServiceServer) Submit(_ context.Context, request *gen.SubmitRequest) (*gen.QuizResults, error) {
	correctCount := 0
	for _, choice := range request.Choices {
		if CheckQuestionChoice(choice) {
			correctCount++
		}
	}
	submissions = append(submissions, correctCount)
	response := &gen.QuizResults{
		CorrectCount:    int32(correctCount),
		ScorePercentile: calcScorePercentile(correctCount),
	}
	return response, nil
}

func CheckQuestionChoice(c *gen.Choice) bool {
	if int(c.QuestionId) >= len(quiz.Questions) {
		log.Fatalf("Invalid question id %d\n", c.QuestionId)
	}
	question := quiz.Questions[c.QuestionId]
	return int(c.AnswerId) == question.CorrectIdx+1
}

func calcScorePercentile(correctCount int) float32 {
	sort.Ints(submissions)
	idx := sort.Search(len(submissions), func(i int) bool {
		return submissions[i] > correctCount
	})
	percentile := float32(idx) / float32(len(submissions)) * 100
	return percentile
}
