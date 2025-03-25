package service

import (
	"errors"

	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
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
	// Validate user
	user, err := s.userRepo.GetByID(caseData.CreatedByID)
	if err != nil {
		return nil, err
	}

	// Only investigators can create cases
	if user.Role != model.RoleInvestigator {
		return nil, errors.New("only investigators can create cases")
	}

	// Save case
	err = s.caseRepo.Create(caseData)
	if err != nil {
		return nil, err
	}

	return caseData, nil
}

func (s *CaseService) UpdateCase(caseData *model.Case) (*model.Case, error) {
	// Validate user
	user, err := s.userRepo.GetByID(caseData.CreatedByID)
	if err != nil {
		return nil, err
	}

	// Only investigators can update cases
	if user.Role != model.RoleInvestigator {
		return nil, errors.New("only investigators can update cases")
	}

	// Update case
	err = s.caseRepo.Update(caseData)
	if err != nil {
		return nil, err
	}

	return caseData, nil
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
	// Validate user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Check if the user can be assigned to the case
	caseData, err := s.caseRepo.GetByID(caseID)
	if err != nil {
		return err
	}

	if user.ClearanceLevel < caseData.AuthorizationLevel {
		return errors.New("user does not have the required clearance level")
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
