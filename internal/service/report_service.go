package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/m7medVision/crime-management-system/internal/dto"
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
)

type ReportService struct {
	reportRepo *repository.ReportRepository
	caseRepo   *repository.CaseRepository
}

func NewReportService(reportRepo *repository.ReportRepository, caseRepo *repository.CaseRepository) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
		caseRepo:   caseRepo,
	}
}

func (s *ReportService) GenerateReport(reportDTO *dto.ReportDTO) (*model.Report, error) {
	// Validate user
	user, err := s.userRepo.GetByID(reportDTO.ReportedByID)
	if err != nil {
		return nil, err
	}

	// Only citizens, admins, and investigators can create reports
	if user.Role != model.RoleCitizen && user.Role != model.RoleAdmin && user.Role != model.RoleInvestigator {
		return nil, errors.New("only citizens, admins, and investigators can create reports")
	}

	// Create report
	report := &model.Report{
		ReportID:     fmt.Sprintf("RPT-%d", time.Now().UnixNano()),
		Title:        reportDTO.Title,
		Description:  reportDTO.Description,
		Location:     reportDTO.Location,
		ReportedByID: reportDTO.ReportedByID,
	}

	err = s.reportRepo.Create(report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) GetReportStatus(reportID uint) (*dto.ReportStatusDTO, error) {
	report, err := s.reportRepo.GetByID(reportID)
	if err != nil {
		return nil, err
	}

	// Get associated cases
	cases, err := s.caseRepo.GetByReportID(reportID)
	if err != nil {
		return nil, err
	}

	// Calculate status
	status := "pending"
	for _, cas := range cases {
		if cas.Status == model.StatusOngoing {
			status = "ongoing"
			break
		} else if cas.Status == model.StatusClosed {
			status = "closed"
		}
	}

	return &dto.ReportStatusDTO{
		ReportID: report.ReportID,
		Status:   status,
	}, nil
}
