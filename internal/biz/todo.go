package biz

import (
	"context"

	v1 "bubble/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Todo struct {
	ID     int64
	Title  string
	Status bool
}

// TodoRepo is a todo repo.
// biz层对数据操作层提出了以下要求，不管实际的存储是MySQL还是Redis。。。
type TodoRepo interface {
	Save(context.Context, *Todo) (*Todo, error)
	Update(context.Context, *Todo) error
	Delete(context.Context, int64) error
	FindByID(context.Context, int64) (*Todo, error)
	ListAll(context.Context) ([]*Todo, error)
}

// TodoUsecase is a todo usecase.
// 这些是UseCase工作时需要的"工具"或"资源"。
type TodoUsecase struct {
	repo TodoRepo
	log  *log.Helper
}

// NewTodoUsecase new a Todo usecase.
// 组装车间或者新员工入职流程。
func NewTodoUsecase(repo TodoRepo, logger log.Logger) *TodoUsecase {
	return &TodoUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
// 具体的业务操作
// 对外提供的业务函数，实现复杂的业务逻辑
func (uc *TodoUsecase) Create(ctx context.Context, t *Todo) (*Todo, error) {
	uc.log.WithContext(ctx).Infof("Create: %#v", t)
	return uc.repo.Save(ctx, t)
}

func (uc *TodoUsecase) Get(ctx context.Context, id int64) (*Todo, error) {
	uc.log.WithContext(ctx).Infof("Create: %#v", id)
	return uc.repo.FindByID(ctx, id)

}
