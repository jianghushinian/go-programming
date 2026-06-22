package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"cronx/internal/pkg/meta"
	"cronx/internal/pkg/model"
	"cronx/pkg/logx"
)

type JobStore interface {
	Create(ctx context.Context, job *model.Job) error
	Get(ctx context.Context, jobID int64) (*model.Job, error)
	List(ctx context.Context, opts ...meta.ListOption) (int64, []*model.Job, error)
	Update(ctx context.Context, job *model.Job) error
	UpdateStatus(ctx context.Context, jobID int64, newStatus, oldStatus model.JobStatus) error
	Delete(ctx context.Context, jobID int64) error
}

type jobStore struct {
	ds *datastore
}

var _ JobStore = (*jobStore)(nil)

func newJobStore(ds *datastore) *jobStore {
	return &jobStore{ds}
}

func (s *jobStore) db(ctx context.Context) *gorm.DB {
	return s.ds.Core(ctx)
}

func (s *jobStore) Create(ctx context.Context, job *model.Job) error {
	return s.db(ctx).Create(&job).Error
}

func (s *jobStore) Get(ctx context.Context, jobID int64) (*model.Job, error) {
	job := &model.Job{}
	if err := s.db(ctx).Where("id = ?", jobID).First(&job).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (s *jobStore) List(ctx context.Context, opts ...meta.ListOption) (count int64, ret []*model.Job, err error) {
	o := meta.NewListOptions(opts...)

	ans := s.db(ctx).
		Where(o.Filters).
		Not(o.Not).
		Offset(o.Offset).
		Limit(defaultLimit(o.Limit)).
		Order(defaultOrder(o.Order)).
		Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count)

	return count, ret, ans.Error
}

func (s *jobStore) Update(ctx context.Context, job *model.Job) error {
	return s.db(ctx).Save(job).Error
}

// UpdateStatus performs a CAS-style (Compare-And-Swap) status update on a job.
//
// It attempts to update the job's status from `oldStatus` to `newStatus`.
// The update only succeeds if the current status in the database matches `oldStatus`.
// This prevents race conditions when multiple workers try to update the same job.
//
// If the update fails due to a database error, the error is returned.
// If no rows are affected (i.e., the status was already changed by someone else),
// it is not considered an error, but a debug log is recorded.
func (s *jobStore) UpdateStatus(ctx context.Context, jobID int64, newStatus, oldStatus model.JobStatus) error {
	ret := s.db(ctx).
		Model(&model.Job{}).
		Where("id = ? AND status = ?", jobID, oldStatus).
		Update("status", newStatus)

	if ret.Error != nil {
		return ret.Error
	}

	if ret.RowsAffected == 0 {
		// Status was updated by another worker; not an error, but log for debugging.
		logx.Debug(ctx, "CAS update skipped: status mismatch",
			"jobID", jobID, "expected", oldStatus, "actual changed to", newStatus)
	}

	return nil
}

func (s *jobStore) Delete(ctx context.Context, jobID int64) error {
	err := s.db(ctx).Where("id = ?", jobID).Delete(&model.Job{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
