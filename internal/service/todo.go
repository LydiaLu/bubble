package service

import (
	"context"
	"errors"

	pb "bubble/api/bubble/v1"
	"bubble/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type TodoService struct {
	pb.UnimplementedTodoServer
	// 每个UseCase代表一个完整的业务功能
	uc  *biz.TodoUsecase
	log *log.Helper
}

func NewTodoService(uc *biz.TodoUsecase) *TodoService {
	return &TodoService{uc: uc}
}

func (s *TodoService) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoReply, error) {
	// 请求来了
	if len(req.GetTitle()) == 0 {
		return &pb.CreateTodoReply{}, errors.New("无效的title")
	}
	// 调用业务逻辑
	data, err := s.uc.Create(ctx, &biz.Todo{Title: req.Title})
	if err != nil {
		return nil, err
	}
	// 返回响应
	return &pb.CreateTodoReply{
		Id:     data.ID,
		Title:  data.Title,
		Status: data.Status,
	}, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoReply, error) {
	return &pb.UpdateTodoReply{}, nil
}
func (s *TodoService) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoReply, error) {
	return &pb.DeleteTodoReply{}, nil
}
func (s *TodoService) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoReply, error) {
	return &pb.GetTodoReply{}, nil
}
func (s *TodoService) ListTodo(ctx context.Context, req *pb.ListTodoRequest) (*pb.ListTodoReply, error) {
	return &pb.ListTodoReply{}, nil
}
