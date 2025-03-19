package service

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
)

type EvidenceService struct {
	evidenceRepo *repository.EvidenceRepository
}

func NewEvidenceService(evidenceRepo *repository.EvidenceRepository) *EvidenceService {
	return &EvidenceService{evidenceRepo: evidenceRepo}
}

func (s *EvidenceService) RecordEvidence(evidence *model.Evidence) (*model.Evidence, error) {
	err := s.evidenceRepo.Create(evidence)
	if err != nil {
		return nil, err
	}
	return evidence, nil
}

func (s *EvidenceService) GetEvidenceByID(evidenceID uint) (*model.Evidence, error) {
	return s.evidenceRepo.GetByID(evidenceID)
}

func (s *EvidenceService) GetEvidenceImageByID(evidenceID uint) (*model.Evidence, error) {
	evidence, err := s.evidenceRepo.GetByID(evidenceID)
	if err != nil {
		return nil, err
	}

	if evidence.Type != model.EvidenceTypeImage {
		return nil, errors.New("evidence is not an image")
	}

	return evidence, nil
}

func (s *EvidenceService) UpdateEvidence(evidence *model.Evidence) (*model.Evidence, error) {
	err := s.evidenceRepo.Update(evidence)
	if err != nil {
		return nil, err
	}
	return evidence, nil
}

func (s *EvidenceService) SoftDeleteEvidence(evidenceID uint) error {
	evidence, err := s.evidenceRepo.GetByID(evidenceID)
	if err != nil {
		return err
	}

	evidence.IsDeleted = true
	return s.evidenceRepo.Update(evidence)
}

func (s *EvidenceService) SaveImage(file *multipart.FileHeader, dest string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return filepath.Base(dest), nil
}
