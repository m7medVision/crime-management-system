package service

import (
	"bytes"
	"fmt"
	"time"

	"codeberg.org/go-pdf/fpdf"
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
	caseRepo *repository.CaseRepository
}

func NewReportService(caseRepo *repository.CaseRepository) (*ReportService, error) {
	return &ReportService{
		caseRepo: caseRepo,
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

	// Generate PDF directly using fpdf
	pdf, err := GeneratePDFFromData(reportData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdf, nil
}

func GeneratePDFFromData(data CaseReportData) ([]byte, error) {
	// Create a new PDF document
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 25)
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.Cell(0, 10, "Case Report")
	pdf.Ln(6)

	// Subtitle
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Case #%d - %s", data.Case.ID, data.Case.Name))
	pdf.Ln(6)

	// Date
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated on %s", data.CreatedAt))
	pdf.Ln(15)

	// Case Information Section
	addSectionTitle(pdf, "Case Information")

	// Case details
	addTableRow(pdf, "Case Number:", fmt.Sprintf("%d", data.Case.ID))
	addTableRow(pdf, "Status:", string(data.Case.Status))
	addTableRow(pdf, "Area/City:", data.Case.Area)
	addTableRow(pdf, "Case Type:", data.Case.CaseType)
	if data.Case.CreatedBy.FullName != "" {
		addTableRow(pdf, "Created By:", data.Case.CreatedBy.FullName)
	}
	addTableRow(pdf, "Created At:", data.Case.CreatedAt.Format("January 2, 2006"))
	addTableRow(pdf, "Authorization Level:", string(data.Case.AuthorizationLevel))
	pdf.Ln(10)

	// Case Description
	addSectionTitle(pdf, "Case Description")
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, data.Case.Description, "", "", false)
	pdf.Ln(10)

	// Evidence Section
	addSectionTitle(pdf, "Evidence")
	if len(data.Evidence) > 0 {
		for i, e := range data.Evidence {
			addItemTitle(pdf, fmt.Sprintf("Evidence #%d", i))
			addItemDetail(pdf, "Type:", string(e.Type))
			if e.AddedBy.FullName != "" {
				addItemDetail(pdf, "Added By:", e.AddedBy.FullName)
			}
			addItemDetail(pdf, "Added On:", e.CreatedAt.Format("January 2, 2006"))

			if e.Type == "text" {
				addItemDetail(pdf, "Content:", e.Content)
			} else if e.Type == "image" {
				addItemDetail(pdf, "Image Path:", e.ImagePath)
			}

			if e.Remarks != "" {
				addItemDetail(pdf, "Remarks:", e.Remarks)
			}
			pdf.Ln(5)
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 10, "No evidence has been recorded for this case.")
		pdf.Ln(10)
	}

	// Suspects Section
	addSectionTitle(pdf, "Suspects")
	if len(data.Suspects) > 0 {
		for i, s := range data.Suspects {
			addItemTitle(pdf, fmt.Sprintf("Suspect #%d", i))
			addItemDetail(pdf, "Name:", fmt.Sprintf("%s %s", s.FirstName, s.LastName))
			addItemDetail(pdf, "Age:", fmt.Sprintf("%d", s.Age))
			addItemDetail(pdf, "Gender:", s.Gender)
			addItemDetail(pdf, "Address:", s.Address)

			isArrested := "No"
			if s.IsArrested {
				isArrested = "Yes"
			}
			addItemDetail(pdf, "Is Arrested:", isArrested)

			if s.Description != "" {
				addItemDetail(pdf, "Description:", s.Description)
			}
			if s.Notes != "" {
				addItemDetail(pdf, "Notes:", s.Notes)
			}
			pdf.Ln(5)
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 10, "No suspects have been recorded for this case.")
		pdf.Ln(10)
	}

	// Victims Section
	addSectionTitle(pdf, "Victims")
	if len(data.Victims) > 0 {
		for i, v := range data.Victims {
			addItemTitle(pdf, fmt.Sprintf("Victim #%d", i))
			addItemDetail(pdf, "Name:", fmt.Sprintf("%s %s", v.FirstName, v.LastName))
			addItemDetail(pdf, "Age:", fmt.Sprintf("%d", v.Age))
			addItemDetail(pdf, "Gender:", v.Gender)
			addItemDetail(pdf, "Address:", v.Address)

			if v.InjuryDescription != "" {
				addItemDetail(pdf, "Injury Description:", v.InjuryDescription)
			}
			if v.Notes != "" {
				addItemDetail(pdf, "Notes:", v.Notes)
			}
			pdf.Ln(5)
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 10, "No victims have been recorded for this case.")
		pdf.Ln(10)
	}

	// Witnesses Section
	addSectionTitle(pdf, "Witnesses")
	if len(data.Witnesses) > 0 {
		for i, w := range data.Witnesses {
			addItemTitle(pdf, fmt.Sprintf("Witness #%d", i))
			addItemDetail(pdf, "Name:", fmt.Sprintf("%s %s", w.FirstName, w.LastName))
			addItemDetail(pdf, "Age:", fmt.Sprintf("%d", w.Age))
			addItemDetail(pdf, "Gender:", w.Gender)
			addItemDetail(pdf, "Address:", w.Address)

			if w.Statement != "" {
				addItemDetail(pdf, "Statement:", w.Statement)
			}
			if w.Notes != "" {
				addItemDetail(pdf, "Notes:", w.Notes)
			}
			pdf.Ln(5)
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 10, "No witnesses have been recorded for this case.")
		pdf.Ln(10)
	}

	// Footer
	pdf.SetY(-25)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 10, "This is an official report generated by the District Core Crime Management System.")
	pdf.Ln(5)
	pdf.Cell(0, 10, "Confidential document. Do not distribute without authorization.")

	// Output the PDF as bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// Helper functions for PDF generation
func addSectionTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, title)
	pdf.Ln(8)
	pdf.SetDrawColor(204, 204, 204)
	pdf.Line(pdf.GetX(), pdf.GetY(), pdf.GetX()+170, pdf.GetY())
	pdf.Ln(5)
}

func addTableRow(pdf *fpdf.Fpdf, label, value string) {
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(50, 6, label)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, value)
	pdf.Ln(6)
}

func addItemTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(0, 123, 255) // Blue color
	pdf.Cell(0, 6, title)
	pdf.Ln(6)
	pdf.SetTextColor(0, 0, 0) // Reset to black
}

func addItemDetail(pdf *fpdf.Fpdf, label, value string) {
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(35, 5, label)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 5, value, "", "", false)
}
