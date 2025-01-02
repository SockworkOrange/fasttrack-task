package cli

import (
	"context"
	"gen"
	"google.golang.org/grpc"
)

// Ensure our quizServiceClient implements the QuizServiceClient interface at compile time.
var _ QuizServiceClient = (*quizServiceClient)(nil)

// NewQuizServiceClient creates a new QuizServiceClient wrapping an existing gRPC connection.
func NewQuizServiceClient(cc grpc.ClientConnInterface) QuizServiceClient {
	return &quizServiceClient{cc}
}

// Get calls the remote Get RPC method.
func (c *quizServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Quiz, error) {
	out := new(Quiz)
	err := c.cc.Invoke(ctx, "/YourProtobufPackage.QuizService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Submit calls the remote Submit RPC method.
func (c *quizServiceClient) Submit(ctx context.Context, in *SubmitRequest, opts ...grpc.CallOption) (*QuizResults, error) {
	out := new(QuizResults)
	err := c.cc.Invoke(ctx, "/YourProtobufPackage.QuizService/Submit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
