package service

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
)

type CaseReportData struct {
	Case      model.Case
	Evidence  []model.Evidence
	Suspects  []model.Suspect
	Victims   []model.Victim
	Witnesses []model.Witness
	CreatedAt string
}

type ReportService struct {
	caseRepo     *repository.CaseRepository
	pdfTemplate  *template.Template
	templatePath string
}

func NewReportService(caseRepo *repository.CaseRepository, templatePath string) (*ReportService, error) {
	// Ensure template directory exists
	if err := os.MkdirAll(filepath.Dir(templatePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create template directory: %w", err)
	}

	// Parse the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &ReportService{
		caseRepo:     caseRepo,
		pdfTemplate:  tmpl,
		templatePath: templatePath,
	}, nil
}

func (s *ReportService) GenerateCaseReport(caseID uint) ([]byte, error) {
	// Get case details
	caseData, err := s.caseRepo.GetByID(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case: %w", err)
	}

	// Get related data
	evidence, err := s.caseRepo.GetEvidence(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence: %w", err)
	}

	suspects, err := s.caseRepo.GetSuspects(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get suspects: %w", err)
	}

	victims, err := s.caseRepo.GetVictims(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get victims: %w", err)
	}

	witnesses, err := s.caseRepo.GetWitnesses(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get witnesses: %w", err)
	}

	// Prepare data for template
	reportData := CaseReportData{
		Case:      *caseData,
		Evidence:  evidence,
		Suspects:  suspects,
		Victims:   victims,
		Witnesses: witnesses,
		CreatedAt: time.Now().Format("January 2, 2006"),
	}

	// Execute template
	out := bytes.NewBuffer([]byte{})
	if err := s.pdfTemplate.Execute(out, reportData); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Generate PDF from LaTeX
	pdf, err := GeneratePDFFromLatex(out.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdf, nil
}

func GeneratePDFFromLatex(src []byte) ([]byte, error) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("./tmp", "report-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up when done

	// Write LaTeX source to temp file
	texFile := filepath.Join(tmpDir, "report.tex")
	if err := os.WriteFile(texFile, src, 0644); err != nil {
		return nil, fmt.Errorf("failed to write tex file: %w", err)
	}

	// Run xelatex command
	cmd := exec.Command("xelatex", "-interaction=nonstopmode", "-output-directory="+tmpDir, texFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("xelatex error: %v, output: %s", err, output)
	}

	// Read the generated PDF
	pdfFile := filepath.Join(tmpDir, "report.pdf")
	pdf, err := os.ReadFile(pdfFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file: %w", err)
	}

	return pdf, nil
}
