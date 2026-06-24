package httpserver

import (
	"net/http"
	"strings"

	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	affiliationsvc "github.com/woragis/minecraft-campus-backend/server/internal/affiliation/service"
)

type catalogHandler struct {
	svc *affiliationsvc.Service
}

func newCatalogHandler(svc *affiliationsvc.Service) *catalogHandler {
	return &catalogHandler{svc: svc}
}

func (h *catalogHandler) listUniversities(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListUniversities(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"universities": rows})
}

func (h *catalogHandler) listFaculties(w http.ResponseWriter, r *http.Request) {
	universitySlug := strings.TrimSpace(r.URL.Query().Get("universitySlug"))
	if universitySlug == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeAffiliationPatchV1HandlerBodyInvalid, "universitySlug query parameter is required."))
		return
	}
	rows, err := h.svc.ListFaculties(r.Context(), universitySlug)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"faculties": rows})
}

func (h *catalogHandler) listCourses(w http.ResponseWriter, r *http.Request) {
	facultySlug := strings.TrimSpace(r.URL.Query().Get("facultySlug"))
	if facultySlug == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeAffiliationPatchV1HandlerBodyInvalid, "facultySlug query parameter is required."))
		return
	}
	rows, err := h.svc.ListCourses(r.Context(), facultySlug)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"courses": rows})
}
