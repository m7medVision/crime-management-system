package service

import (
	"errors"

	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
	"github.com/m7medVision/crime-management-system/internal/util"
)

type CaseService struct {
	caseRepo *repository.CaseRepository
	userRepo *repository.UserRepository
}

func NewCaseService(caseRepo *repository.CaseRepository, userRepo *repository.UserRepository) *CaseService {
	return &CaseService{
		caseRepo: caseRepo,
		userRepo: userRepo,
	}
}

func (s *CaseService) CreateCase(caseData *model.Case) (*model.Case, error) {
	return nil, s.caseRepo.Create(caseData)
}

func (s *CaseService) UpdateCase(caseData *model.Case) (*model.Case, error) {
	return nil, s.caseRepo.Update(caseData)
}

func (s *CaseService) GetCaseByID(caseID uint) (*model.Case, error) {
	return s.caseRepo.GetByID(caseID)
}

func (s *CaseService) ListCases(offset, limit int, search string) ([]model.Case, int64, error) {
	return s.caseRepo.List(offset, limit, search)
}

func (s *CaseService) GetAssignees(caseID uint) ([]model.User, error) {
	return s.caseRepo.GetAssignees(caseID)
}

func (s *CaseService) AddAssignee(caseID, userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	caseData, err := s.caseRepo.GetByID(caseID)
	if err != nil {
		return err
	}

	if !util.IsClearnceLevelHigherOrEqual(user.ClearanceLevel, caseData.AuthorizationLevel) {
		return errors.New("insufficient clearance level")
	}

	return s.caseRepo.AddAssignee(caseID, userID)
}

func (s *CaseService) RemoveAssignee(caseID, userID uint) error {
	return s.caseRepo.RemoveAssignee(caseID, userID)
}

func (s *CaseService) SubmitCrimeReport(report *model.Report) (*model.Report, error) {
	if err := s.caseRepo.CreateReport(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *CaseService) GetEvidence(caseID uint) ([]model.Evidence, error) {
	return s.caseRepo.GetEvidence(caseID)
}

func (s *CaseService) GetSuspects(caseID uint) ([]model.Suspect, error) {
	return s.caseRepo.GetSuspects(caseID)
}

func (s *CaseService) GetVictims(caseID uint) ([]model.Victim, error) {
	return s.caseRepo.GetVictims(caseID)
}

func (s *CaseService) GetWitnesses(caseID uint) ([]model.Witness, error) {
	return s.caseRepo.GetWitnesses(caseID)
}

func (s *CaseService) ExtractLinksFromCase(caseID uint) ([]string, error) {
	return s.caseRepo.ExtractLinks(caseID)
}

func (s *CaseService) GetCaseStatusByReportID(reportID string) (model.CaseStatus, error) {
	status, err := s.caseRepo.GetStatusByReportID(reportID)
	if err != nil {
		return "", err
	}
	return status, nil
}
