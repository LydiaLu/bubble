package service

import (
	"context"
	"errors"
	"time"

	pb "bubble/api/bubble/v1"
	v1 "bubble/api/bubble/v1"
	"bubble/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type TodoService struct {
	pb.UnimplementedTodoServer
	// 每个UseCase代表一个完整的业务功能
	uc  *biz.TodoUsecase
	log *log.Helper
	// 存储评估结果的map
	evaluationResults map[uuid.UUID]bool
}

func NewTodoService(uc *biz.TodoUsecase) *TodoService {
	return &TodoService{uc: uc, evaluationResults: make(map[uuid.UUID]bool)}
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

// EvaluateTodo 实现评估待办事项的接口
func (s *TodoService) EvaluateTodo(ctx context.Context, req *pb.EvaluateTodoRequest) (*pb.EvaluateTodoReply, error) {
	evaluation_id := uuid.New()
	go s.performEvaluation(evaluation_id)

	return &pb.EvaluateTodoReply{
		EvaluationId: evaluation_id.String(),
	}, nil
}

// 执行耗时的评估操作
func (s *TodoService) performEvaluation(id uuid.UUID) {
	// 模拟耗时操作
	time.Sleep(20 * time.Second)

	s.evaluationResults[id] = true
}

func (s *TodoService) GetEvaluationStatus(ctx context.Context, req *pb.GetEvaluationStatusRequest) (*pb.GetEvaluateStatusdoReply, error) {

	evaluation_id, _ := uuid.Parse(req.EvaluationId)
	completed, exists := s.evaluationResults[evaluation_id]

	if !exists {
		return &pb.GetEvaluateStatusdoReply{
			Message:   "未找到评估任务",
			Completed: false,
		}, nil
	}

	if completed {
		return &pb.GetEvaluateStatusdoReply{
			Message:   "评估完成",
			Completed: true,
		}, nil
	}

	return &pb.GetEvaluateStatusdoReply{
		Message:   "正在评估中",
		Completed: false,
	}, nil
}
