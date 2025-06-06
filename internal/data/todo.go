package data

import (
	"context"

	"bubble/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// todoRepo实现了biz层定义的repo接口
type todoRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewTodoRepo(data *Data, logger log.Logger) biz.TodoRepo {
	return &todoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *todoRepo) Save(ctx context.Context, t *biz.Todo) (*biz.Todo, error) {
	// 实现数据库的操作
	// gorm操作
	err := r.data.db.Create(t).Error
	return t, err
}

func (r *todoRepo) Delete(ctx context.Context, id int64) error {
	t := biz.Todo{ID: id}
	err := r.data.db.Delete(&t).Error
	return err
}

func (r *todoRepo) Update(ctx context.Context, t *biz.Todo) error {
	err := r.data.db.Model(t).Update("status", t.Status).Error
	return err
}

func (r *todoRepo) FindByID(ctx context.Context, id int64) (*biz.Todo, error) {
	t := biz.Todo{ID: id}
	err := r.data.db.First(&t).Error
	return &t, err
}

func (r *todoRepo) ListAll(ctx context.Context) ([]*biz.Todo, error) {
	var todoList []*biz.Todo
	err := r.data.db.WithContext(ctx).Find(&todoList).Error
	return todoList, err
}
