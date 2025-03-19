package repository

import (
	"errors"

	"github.com/m7medVision/crime-management-system/internal/model"
	"gorm.io/gorm"
)

type CaseRepository struct {
	db *gorm.DB
}

func NewCaseRepository(db *gorm.DB) *CaseRepository {
	return &CaseRepository{db: db}
}

func (r *CaseRepository) Create(cas *model.Case) error {
	return r.db.Create(cas).Error
}

func (r *CaseRepository) GetByID(id uint) (*model.Case, error) {
	var cas model.Case
	result := r.db.Preload("CreatedBy").Preload("ReportedBy").First(&cas, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cas, nil
}

func (r *CaseRepository) GetByCaseNumber(caseNumber string) (*model.Case, error) {
	var cas model.Case
	result := r.db.Preload("CreatedBy").Preload("ReportedBy").Where("case_number = ?", caseNumber).First(&cas)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cas, nil
}

func (r *CaseRepository) Update(cas *model.Case) error {
	return r.db.Save(cas).Error
}

func (r *CaseRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Case{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("case not found")
	}
	return nil
}

func (r *CaseRepository) List(offset, limit int, search string) ([]model.Case, int64, error) {
	var cases []model.Case
	var count int64

	query := r.db.Model(&model.Case{})

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("CreatedBy").Offset(offset).Limit(limit).Find(&cases).Error
	return cases, count, err
}

func (r *CaseRepository) GetAssignees(caseID uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Model(&model.Case{Model: gorm.Model{ID: caseID}}).Association("Assignees").Find(&users)
	return users, err
}

func (r *CaseRepository) AddAssignee(caseID, userID uint) error {
	return r.db.Model(&model.Case{Model: gorm.Model{ID: caseID}}).Association("Assignees").Append(&model.User{Model: gorm.Model{ID: userID}})
}

func (r *CaseRepository) RemoveAssignee(caseID, userID uint) error {
	return r.db.Model(&model.Case{Model: gorm.Model{ID: caseID}}).Association("Assignees").Delete(&model.User{Model: gorm.Model{ID: userID}})
}
