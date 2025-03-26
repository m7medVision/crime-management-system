package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
	"github.com/m7medVision/crime-management-system/internal/util"
	"github.com/minio/minio-go/v7"
)

type EvidenceService struct {
	evidenceRepo *repository.EvidenceRepository
	caseRepo     *repository.CaseRepository
	userRepo     *repository.UserRepository
	minioClient  *minio.Client
	minioBucket  string
}

func NewEvidenceService(
	evidenceRepo *repository.EvidenceRepository,
	caseRepo *repository.CaseRepository,
	userRepo *repository.UserRepository,
	minioBucket string,
) *EvidenceService {
	return &EvidenceService{
		evidenceRepo: evidenceRepo,
		caseRepo:     caseRepo,
		userRepo:     userRepo,
		minioClient:  util.GetMinioClient(),
		minioBucket:  minioBucket,
	}
}

func (s *EvidenceService) CreateTextEvidence(caseID, userID uint, content, remarks string) (*model.Evidence, error) {
	// Check if case exists
	_, err := s.caseRepo.GetByID(caseID)
	if err != nil {
		return nil, errors.New("case not found")
	}

	// Check if user exists
	_, err = s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	evidence := &model.Evidence{
		CaseID:    caseID,
		Type:      model.EvidenceTypeText,
		Content:   content,
		Remarks:   remarks,
		AddedByID: userID,
		IsDeleted: false,
	}

	if err := s.evidenceRepo.Create(evidence); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := &model.AuditLog{
		UserID:     userID,
		Action:     model.ActionCreate,
		EntityType: "evidence",
		EntityID:   evidence.ID,
		NewValue:   fmt.Sprintf("Text evidence created for case %d", caseID),
	}
	s.evidenceRepo.CreateAuditLog(auditLog)

	return evidence, nil
}

func (s *EvidenceService) CreateImageEvidence(caseID, userID uint, file *multipart.FileHeader, remarks string) (*model.Evidence, error) {
	// Check if case exists
	_, err := s.caseRepo.GetByID(caseID)
	if err != nil {
		return nil, errors.New("case not found")
	}

	// Check if user exists
	_, err = s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate image file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Read the file into a buffer for MIME type detection
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// Reset file pointer
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Detect MIME type
	mtype := mimetype.Detect(buffer)
	if !strings.HasPrefix(mtype.String(), "image/") {
		return nil, errors.New("uploaded file is not an image")
	}

	// Generate unique filename
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("evidence/%d/%s", caseID, filename)

	// Upload file to MinIO
	fileContent, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	_, err = s.minioClient.PutObject(
		context.Background(),
		s.minioBucket,
		objectName,
		bytes.NewReader(fileContent),
		int64(len(fileContent)),
		minio.PutObjectOptions{ContentType: mtype.String()},
	)
	if err != nil {
		return nil, err
	}

	// Create evidence record
	evidence := &model.Evidence{
		CaseID:    caseID,
		Type:      model.EvidenceTypeImage,
		ImagePath: objectName,
		Remarks:   remarks,
		AddedByID: userID,
		IsDeleted: false,
	}

	if err := s.evidenceRepo.Create(evidence); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := &model.AuditLog{
		UserID:     userID,
		Action:     model.ActionCreate,
		EntityType: "evidence",
		EntityID:   evidence.ID,
		NewValue:   fmt.Sprintf("Image evidence uploaded for case %d", caseID),
	}
	s.evidenceRepo.CreateAuditLog(auditLog)

	return evidence, nil
}

func (s *EvidenceService) GetEvidenceByID(id uint) (*model.Evidence, error) {
	return s.evidenceRepo.GetByID(id)
}

func (s *EvidenceService) GetEvidenceImage(id uint) (io.ReadCloser, int64, string, error) {
	evidence, err := s.evidenceRepo.GetByID(id)
	if err != nil {
		return nil, 0, "", err
	}

	if evidence.Type != model.EvidenceTypeImage {
		return nil, 0, "", errors.New("evidence is not an image")
	}

	obj, err := s.minioClient.GetObject(
		context.Background(),
		s.minioBucket,
		evidence.ImagePath,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, 0, "", err
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, 0, "", err
	}

	return obj, stat.Size, stat.ContentType, nil
}

func (s *EvidenceService) UpdateEvidence(id, userID uint, remarks string) (*model.Evidence, error) {
	evidence, err := s.evidenceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Store old value for audit log
	oldValue := evidence.Remarks

	// Update remarks only
	evidence.Remarks = remarks

	if err := s.evidenceRepo.Update(evidence); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := &model.AuditLog{
		UserID:     userID,
		Action:     model.ActionUpdate,
		EntityType: "evidence",
		EntityID:   evidence.ID,
		OldValue:   oldValue,
		NewValue:   remarks,
	}
	s.evidenceRepo.CreateAuditLog(auditLog)

	return evidence, nil
}

func (s *EvidenceService) SoftDeleteEvidence(id, userID uint) error {
	evidence, err := s.evidenceRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.evidenceRepo.SoftDelete(id); err != nil {
		return err
	}

	// Create audit log
	auditLog := &model.AuditLog{
		UserID:     userID,
		Action:     model.ActionDelete,
		EntityType: "evidence",
		EntityID:   evidence.ID,
		OldValue:   fmt.Sprintf("Evidence %d soft deleted", id),
	}
	s.evidenceRepo.CreateAuditLog(auditLog)

	return nil
}

func (s *EvidenceService) HardDeleteEvidence(id, userID uint) error {
	evidence, err := s.evidenceRepo.GetByID(id)
	if err != nil {
		return err
	}

	// If it's an image, delete it from MinIO
	if evidence.Type == model.EvidenceTypeImage {
		err = s.minioClient.RemoveObject(
			context.Background(),
			s.minioBucket,
			evidence.ImagePath,
			minio.RemoveObjectOptions{},
		)
		if err != nil {
			return err
		}
	}

	if err := s.evidenceRepo.HardDelete(id); err != nil {
		return err
	}

	// Create audit log
	auditLog := &model.AuditLog{
		UserID:     userID,
		Action:     model.ActionDelete,
		EntityType: "evidence",
		EntityID:   evidence.ID,
		OldValue:   fmt.Sprintf("Evidence %d hard deleted", id),
	}
	s.evidenceRepo.CreateAuditLog(auditLog)

	return nil
}

func (s *EvidenceService) ListEvidenceByCaseID(caseID uint) ([]model.Evidence, error) {
	return s.evidenceRepo.ListByCaseID(caseID)
}

func (s *EvidenceService) GetEvidenceAuditLogs(evidenceID uint) ([]model.AuditLog, error) {
	return s.evidenceRepo.GetAuditLogsForEvidence(evidenceID)
}
