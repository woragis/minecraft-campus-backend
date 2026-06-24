package service

import (
	"context"
	"errors"
	"strings"

	"github.com/woragis/minecraft-campus-backend/server/internal/affiliation/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"gorm.io/gorm"
)

type CatalogLabels struct {
	UniversityName string
	UniversityHex  string
	FacultyName    string
	FacultyAbbr    string
	FacultyHex     string
	CourseName     string
	CourseAbbr     string
	CourseHex      string
}

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListUniversities(ctx context.Context) ([]models.University, error) {
	rows, err := s.repo.ListUniversities(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeCatalogListV1ServiceFailed, apperrors.MsgCatalogListV1ServiceFailed, err)
	}
	return rows, nil
}

func (s *Service) ListFaculties(ctx context.Context, universitySlug string) ([]models.Faculty, error) {
	rows, err := s.repo.ListFaculties(ctx, strings.TrimSpace(universitySlug))
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeCatalogListV1ServiceFailed, apperrors.MsgCatalogListV1ServiceFailed, err)
	}
	return rows, nil
}

func (s *Service) ListCourses(ctx context.Context, facultySlug string) ([]models.Course, error) {
	rows, err := s.repo.ListCourses(ctx, strings.TrimSpace(facultySlug))
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeCatalogListV1ServiceFailed, apperrors.MsgCatalogListV1ServiceFailed, err)
	}
	return rows, nil
}

func (s *Service) ResolveLabels(ctx context.Context, player *models.Player) CatalogLabels {
	if player == nil {
		return CatalogLabels{}
	}
	out := CatalogLabels{}
	if player.UniversitySlug != nil && *player.UniversitySlug != "" {
		if uni, err := s.repo.FindUniversity(ctx, *player.UniversitySlug); err == nil {
			out.UniversityName = uni.Name
			out.UniversityHex = uni.ColorHex
		}
	}
	if player.FacultySlug != nil && *player.FacultySlug != "" {
		if fac, err := s.repo.FindFaculty(ctx, *player.FacultySlug); err == nil {
			out.FacultyName = fac.Name
			out.FacultyAbbr = fac.ShortAbbr
			out.FacultyHex = fac.ColorHex
		}
	}
	if player.CourseSlug != nil && *player.CourseSlug != "" {
		if course, err := s.repo.FindCourse(ctx, *player.CourseSlug); err == nil {
			out.CourseName = course.Name
			out.CourseAbbr = course.ShortAbbr
			out.CourseHex = course.ColorHex
		}
	}
	return out
}

type AffiliationInput struct {
	AffiliationType string
	UniversitySlug  *string
	FacultySlug     *string
	CourseSlug      *string
}

func (s *Service) ValidateAndNormalize(ctx context.Context, in AffiliationInput) (AffiliationInput, error) {
	affType := strings.TrimSpace(strings.ToLower(in.AffiliationType))
	if affType == "" {
		affType = models.AffiliationStudent
	}
	if !models.ValidAffiliationType(affType) {
		return AffiliationInput{}, apperrors.Invalid(apperrors.CodeAffiliationPatchV1ServiceTypeInvalid, apperrors.MsgAffiliationPatchV1ServiceTypeInvalid)
	}

	normalized := AffiliationInput{AffiliationType: affType}
	if affType == models.AffiliationGuest {
		return normalized, nil
	}

	uni := trimPtr(in.UniversitySlug)
	fac := trimPtr(in.FacultySlug)
	course := trimPtr(in.CourseSlug)

	if affType == models.AffiliationStudent || affType == models.AffiliationAlumni {
		if uni == nil || fac == nil || course == nil {
			return AffiliationInput{}, apperrors.Invalid(apperrors.CodeAffiliationPatchV1ServiceCatalogIncomplete, apperrors.MsgAffiliationPatchV1ServiceCatalogIncomplete)
		}
	}

	if uni != nil {
		if _, err := s.repo.FindUniversity(ctx, *uni); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return AffiliationInput{}, apperrors.NotFound(apperrors.CodeAffiliationPatchV1ServiceUniversityNotFound, apperrors.MsgAffiliationPatchV1ServiceUniversityNotFound)
			}
			return AffiliationInput{}, apperrors.InternalCause(apperrors.CodeAffiliationPatchV1ServiceFailed, apperrors.MsgAffiliationPatchV1ServiceFailed, err)
		}
		normalized.UniversitySlug = uni
	}
	if fac != nil {
		faculty, err := s.repo.FindFaculty(ctx, *fac)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return AffiliationInput{}, apperrors.NotFound(apperrors.CodeAffiliationPatchV1ServiceFacultyNotFound, apperrors.MsgAffiliationPatchV1ServiceFacultyNotFound)
			}
			return AffiliationInput{}, apperrors.InternalCause(apperrors.CodeAffiliationPatchV1ServiceFailed, apperrors.MsgAffiliationPatchV1ServiceFailed, err)
		}
		if uni != nil && faculty.UniversitySlug != *uni {
			return AffiliationInput{}, apperrors.Invalid(apperrors.CodeAffiliationPatchV1ServiceFacultyMismatch, apperrors.MsgAffiliationPatchV1ServiceFacultyMismatch)
		}
		normalized.FacultySlug = fac
	}
	if course != nil {
		courseRow, err := s.repo.FindCourse(ctx, *course)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return AffiliationInput{}, apperrors.NotFound(apperrors.CodeAffiliationPatchV1ServiceCourseNotFound, apperrors.MsgAffiliationPatchV1ServiceCourseNotFound)
			}
			return AffiliationInput{}, apperrors.InternalCause(apperrors.CodeAffiliationPatchV1ServiceFailed, apperrors.MsgAffiliationPatchV1ServiceFailed, err)
		}
		if fac != nil && courseRow.FacultySlug != *fac {
			return AffiliationInput{}, apperrors.Invalid(apperrors.CodeAffiliationPatchV1ServiceCourseMismatch, apperrors.MsgAffiliationPatchV1ServiceCourseMismatch)
		}
		normalized.CourseSlug = course
	}

	return normalized, nil
}

func trimPtr(v *string) *string {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	return &s
}

func NormalizeInviteAffiliationType(raw string) string {
	t := strings.TrimSpace(strings.ToLower(raw))
	if t == "" {
		return models.AffiliationStudent
	}
	if models.ValidAffiliationType(t) {
		return t
	}
	return models.AffiliationStudent
}
