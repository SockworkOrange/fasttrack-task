package client

import (
	"bufio"
	"context"
	"fasttrack-task/pkg/gen"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func Run() {
	conn, err := grpc.NewClient(
		"127.0.0.1:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	client := gen.NewQuizServiceClient(conn)
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)

	quiz := getQuiz(client)
	choices := answerQuiz(quiz)
	submitQuiz(client, choices)
}

func getQuiz(client gen.QuizServiceClient) *gen.Quiz {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	quiz, err := client.Get(ctx, &gen.GetRequest{})
	if err != nil {
		log.Fatalf("Failed to get quiz: %v\n", err)
	}
	return quiz
}

func answerQuiz(quiz *gen.Quiz) []*gen.Choice {
	var choices []*gen.Choice
	for questionId, question := range quiz.Questions {
		fmt.Printf("%d - %s\n", questionId+1, question.Text)
		for answerId, answerText := range question.Possibilities {
			fmt.Printf("\t%d - %s\n", answerId+1, answerText)
		}
		choice := &gen.Choice{
			QuestionId: int32(questionId),
			AnswerId:   int32(readChoiceNumber(len(question.Possibilities))),
		}
		choices = append(choices, choice)
		fmt.Println()
	}
	return choices
}

func readChoiceNumber(n int) int {
	var inputNumber int
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Enter your choice (1..%d): ", n)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input, please try again.")
			continue
		}

		if num < 1 || num > n {
			fmt.Printf("Number must be between 1 and %d.\n", n)
			continue
		}
		inputNumber = num
		break
	}

	return inputNumber
}

func submitQuiz(client gen.QuizServiceClient, choices []*gen.Choice) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	submitReq := &gen.SubmitRequest{
		Choices: choices,
	}
	quizResults, err := client.Submit(ctx, submitReq)
	if err != nil {
		log.Fatalf("Failed to submit quiz: %v\n", err)
	}
	prettyPrintQuizResults(quizResults)
}

func prettyPrintQuizResults(result *gen.QuizResults) {
	fmt.Printf("You answered %+v questions correctly!\n", result.CorrectCount)
	fmt.Printf("You were better than %+v%% of all quizzers!\n", result.ScorePercentile)
}
