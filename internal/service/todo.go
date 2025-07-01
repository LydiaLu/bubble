package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	pb "bubble/api/bubble/v1"
	v1 "bubble/api/bubble/v1"
	"bubble/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type TodoService struct {
	pb.UnimplementedTodoServer
	// 每个UseCase代表一个完整的业务功能
	uc  *biz.TodoUsecase
	log *log.Helper
	// 存储评估结果的map
	evaluationResults sync.Map
}

func NewTodoService(uc *biz.TodoUsecase) *TodoService {
	return &TodoService{uc: uc, evaluationResults: sync.Map{}}
}

func (s *TodoService) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoReply, error) {
	// 请求来了
	// 1. 请求参数的校验
	if len(req.GetTitle()) == 0 {
		return &pb.CreateTodoReply{}, errors.New("无效的title")
	}
	// 2. 调用业务逻辑
	data, err := s.uc.Create(ctx, &biz.Todo{Title: req.Title})
	if err != nil {
		return nil, err
	}
	// 3. 返回响应
	return &pb.CreateTodoReply{
		Id:     data.ID,
		Title:  data.Title,
		Status: data.Status,
	}, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoReply, error) {
	// 1. 参数处理
	// 2. 调用biz层业务逻辑
	err := s.uc.Update(ctx, &biz.Todo{
		ID:     req.Id,
		Title:  req.Title,
		Status: req.Status,
	})
	return &pb.UpdateTodoReply{}, err
}

func (s *TodoService) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoReply, error) {
	// 1. 参数处理
	if req.Id <= 0 {
		return &pb.DeleteTodoReply{}, errors.New("无效的id")
	}
	// 2. 调用biz层业务逻辑
	err := s.uc.Delete(ctx, req.Id)
	return &pb.DeleteTodoReply{}, err
}
func (s *TodoService) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoReply, error) {
	// 1. 参数处理
	if req.Id <= 0 {
		return &pb.GetTodoReply{}, errors.New("无效的id")
	}
	// 2. 调用biz层业务逻辑
	ret, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		// return nil, err
		// 返回自定义错误
		return nil, v1.ErrorTodoNotFound("id: %v to do is not found", req.Id)
	}
	// 3. 按格式返回响应
	return &pb.GetTodoReply{Todo: &pb.Todo{
		Id:     ret.ID,
		Title:  ret.Title,
		Status: ret.Status,
	}}, nil
}
func (s *TodoService) ListTodo(ctx context.Context, req *pb.ListTodoRequest) (*pb.ListTodoReply, error) {
	// 1. 参数处理 => 没有参数
	// 2. 调用biz层业务逻辑
	dataList, err := s.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	reply := &pb.ListTodoReply{}
	for _, data := range dataList {
		reply.Data = append(reply.Data, &pb.Todo{
			Id:     data.ID,
			Title:  data.Title,
			Status: data.Status,
		})
	}
	return reply, nil
}

func (s *TodoService) EvaluateTodo(ctx context.Context, req *pb.EvaluateTodoRequest) (*pb.EvaluateTodoReply, error) {
	if req.Id <= 0 {
		return nil, errors.New("无效的待办事项ID")
	}

	todoID := req.Id
	todoIDStr := fmt.Sprintf("%d", todoID)

	// 检查任务状态
	if status, exists := s.evaluationResults.Load(todoIDStr); exists {
		if status == "进行中" {
			return &pb.EvaluateTodoReply{
				Message: "评估任务已在执行中",
			}, nil
		} else if status == "已完成" {
			return &pb.EvaluateTodoReply{
				Message: "评估任务已完成",
			}, nil
		}
	}

	// 更新状态为"进行中"
	s.evaluationResults.Store(todoIDStr, "进行中")

	// 模拟耗时操作（同步执行）
	time.Sleep(20 * time.Second)

	// 更新评估状态为"已完成"
	s.evaluationResults.Store(todoIDStr, "已完成")
	return &pb.EvaluateTodoReply{
		Message: "评估成功完成",
	}, nil
}

func (s *TodoService) GetEvaluationStatus(ctx context.Context, req *pb.GetEvaluationStatusRequest) (*pb.GetEvaluateStatusdoReply, error) {
	todoIDStr := fmt.Sprintf("%d", req.Id)

	// 获取评估状态
	if status, exists := s.evaluationResults.Load(todoIDStr); exists {
		return &pb.GetEvaluateStatusdoReply{
			Status: status.(string),
		}, nil
	}

	return &pb.GetEvaluateStatusdoReply{
		Status: "未开始",
	}, nil
}
