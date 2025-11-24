package usecase

import (
	"context"

	auditlog "github.com/Flood-project/backend-flood/internal/audit_log"
	"github.com/Flood-project/backend-flood/internal/audit_log/repository"
)

type AuditLogUseCase interface {
	Create(ctx context.Context, log *auditlog.AuditLog) error
	Fetch() ([]auditlog.AuditLog, error)
}

type auditLogUseCase struct {
	logRepo repository.AuditLogManagement
}

func NewAuditLogUseCase(br repository.AuditLogManagement) AuditLogUseCase {
	return &auditLogUseCase {
		logRepo: br,
	}
}

func (us *auditLogUseCase) Create(ctx context.Context, log *auditlog.AuditLog) error {
	err := us.logRepo.Create(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

func (us *auditLogUseCase) Fetch() ([]auditlog.AuditLog, error) {
	logs, err := us.logRepo.Fetch()
	if err != nil {
		return nil, err
	}

	return logs, nil
}