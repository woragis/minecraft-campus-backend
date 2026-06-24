package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListUniversities(ctx context.Context) ([]models.University, error) {
	var rows []models.University
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list universities: %w", err)
	}
	return rows, nil
}

func (r *Repository) ListFaculties(ctx context.Context, universitySlug string) ([]models.Faculty, error) {
	var rows []models.Faculty
	q := r.db.WithContext(ctx).Order("name ASC")
	if universitySlug != "" {
		q = q.Where("university_slug = ?", universitySlug)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list faculties: %w", err)
	}
	return rows, nil
}

func (r *Repository) ListCourses(ctx context.Context, facultySlug string) ([]models.Course, error) {
	var rows []models.Course
	q := r.db.WithContext(ctx).Order("name ASC")
	if facultySlug != "" {
		q = q.Where("faculty_slug = ?", facultySlug)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	return rows, nil
}

func (r *Repository) FindUniversity(ctx context.Context, slug string) (*models.University, error) {
	var row models.University
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find university: %w", err)
	}
	return &row, nil
}

func (r *Repository) FindFaculty(ctx context.Context, slug string) (*models.Faculty, error) {
	var row models.Faculty
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find faculty: %w", err)
	}
	return &row, nil
}

func (r *Repository) FindCourse(ctx context.Context, slug string) (*models.Course, error) {
	var row models.Course
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find course: %w", err)
	}
	return &row, nil
}
