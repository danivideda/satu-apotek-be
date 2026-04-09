package handler

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

type pharmacyHandler struct {
	repo repository.Repository
}

type pharmacyJSON struct {
	AppID     any `json:"app_id"`
	Name      any `json:"name"`
	CreatedAt any `json:"created_at,omitempty"`
	UpdatedAt any `json:"updated_at,omitempty"`
}

func (h *pharmacyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	var payload struct {
		ApotekName string `json:"apotek_name"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	apotek, err := h.repo.Pharmacies.Create(ctx, authOwner.ID, payload.ApotekName)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := pharmacyJSON{
		AppID:     apotek.AppID,
		Name:      apotek.Name,
		CreatedAt: apotek.CreatedAt,
		UpdatedAt: apotek.UpdatedAt,
	}
	if err := json.ResponseCreated(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *pharmacyHandler) GetByOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	pharmacies, err := h.repo.Pharmacies.GetByOwnerID(ctx, authOwner.ID)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := []pharmacyJSON{}
	for _, pharmacy := range *pharmacies {
		item := pharmacyJSON{
			AppID:     pharmacy.AppID,
			Name:      pharmacy.Name,
			CreatedAt: pharmacy.CreatedAt,
			UpdatedAt: pharmacy.UpdatedAt,
		}
		res = append(res, item)
	}
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *pharmacyHandler) Connect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO:
	// Connect the Pharmacy to the browser

	// 1. Client send a request to connect a pharmacy with a code
	var payload struct {
		Code string `json:"code"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// 2. Find the code associated in pharmacy_codes table
	pharmacyCode, err := h.repo.Pharmacies.GetCodeByCode(ctx, payload.Code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			json.ResponseBadRequest(w, r, fmt.Errorf("%w: no pharmacy_code matched", err))
		} else {
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}
	// --check if expired
	if time.Now().After(pharmacyCode.ExpiresAt.Time) {
		json.ResponseBadRequest(w, r, ErrPharmacyCodeExpired)
		return
	}

	// 3. Create pharmacy_sessions and send the pharmacy_session cookies to the user
	pharmacySession, err := h.repo.PharmacySessions.Create(ctx, pharmacyCode.ApotekID, time.Now().Add(5 * time.Minute))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			json.ResponseBadRequest(w, r, err)
		} else {
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}

	service.SetPharmacyCookies(w, pharmacySession.ID.String(), time.Now().Add(5 * time.Minute))
	if err := json.ResponseNoContent(w); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *pharmacyHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. client send which pharmacy to connect: app_id
	var payload struct {
		AppID string `json:"app_id"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// 2. check if app_id is indeed the Owner's
	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	pharmacy, err := h.repo.Pharmacies.GetByAppID(ctx, payload.AppID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			json.ResponseBadRequest(w, r, err)
		} else {
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}

	if pharmacy.OwnerID != authOwner.ID {
		json.ResponseForbidden(w, r, ErrAppIDNotAllowed)
		return
	}

	// 3. Create new code for connection in pharmacy_codes table
	exec := func() (*dbsqlc.PharmacyCode, error) {
		newCode, err := h.generateApotekCode()
		if err != nil {
			return nil, err
		}
		return h.repo.Pharmacies.UpsertCode(ctx, pharmacy.ID, newCode)
	}

	// Retry upsert if pharmacy_code is duplicate (unique constraint sql)
	retry := 5
	var pharmacyCode *dbsqlc.PharmacyCode
	for retry > 0 {
		pharmacyCode, err = exec()
		if err != nil {
			if errors.Is(err, repository.ErrDuplicateValue) {
				// if it's duplicate err, reduce the retry
				retry -= 1
			} else {
				// if it's not a duplicate err, don't retry.
				// Handle the error immediately after the for loop
				retry = 0
			}
		} else {
			// if there's no error, it means it's success
			retry = 0
		}
	}
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := struct {
		pharmacyJSON
		Code      string    `json:"code"`
		ExpiredAt time.Time `json:"expired_at"`
	}{
		pharmacyJSON: pharmacyJSON{
			AppID: pharmacy.AppID,
			Name:  pharmacy.Name,
		},
		Code:      pharmacyCode.Code,
		ExpiredAt: pharmacyCode.ExpiresAt.Time,
	}
	if err := json.ResponseCreated(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *pharmacyHandler) generateApotekCode() (string, error) {
	b := make([]byte, 3) // 3 bytes for 6-character hex string code
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	encoded := hex.EncodeToString(b)
	return encoded, nil
}
