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

type PharmacyJSON struct {
	AppID     any `json:"app_id"`
	Name      any `json:"name"`
	Address   any `json:"address"`
	CreatedAt any `json:"created_at,omitempty"`
	UpdatedAt any `json:"updated_at,omitempty"`
}

type UserJSON struct {
	ID       *int64 `json:"id,omitempty"`
	Username string `json:"username"`
}

type ApotekCodeJSON struct {
	Code      string `json:"code"`
	ExpiresAt any    `json:"expires_at"`
}

func (h *pharmacyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	var payload struct {
		ApotekName string `json:"apotek_name" validate:"required,min=1,max=100,alphanumspace"`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
		return
	}

	apotek, err := h.repo.Pharmacies.Create(ctx, authOwner.ID, payload.ApotekName)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := PharmacyJSON{
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

	res := []PharmacyJSON{}
	for _, pharmacy := range *pharmacies {
		item := PharmacyJSON{
			AppID:   pharmacy.AppID,
			Name:    pharmacy.Name,
			Address: pharmacy.Address,
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
	// 1. Client send a request to connect a pharmacy with a code
	var payload struct {
		Code string `json:"code" validate:"required,hexadecimal,len=6"`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
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

	// 3. Delete pharmacy code after used
	if _, err := h.repo.Pharmacies.DeleteCode(ctx, pharmacyCode.Code); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	// 4. Create pharmacy_sessions
	pharmacySession, err := h.repo.PharmacySessions.Create(ctx, pharmacyCode.ApotekID, time.Now().Add(5*time.Minute))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			json.ResponseBadRequest(w, r, err)
		} else {
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}

	// 5. Set pharmacy cache
	users, err := service.GetUsersFromPharmacyID(ctx, h.repo, pharmacySession.PharmacyID)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	service.SetPharmacyCookies(w, pharmacySession.ID.String(), time.Now().Add(5*time.Minute))
	h.repo.CacheStore.PharmacySessions.SetDefault(pharmacySession.ID.String(), repository.PharmacySessionCacheValue{
		PharmacyID: pharmacySession.PharmacyID,
		Users:      *users,
	})

	if err := json.ResponseNoContent(w); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *pharmacyHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. get pharmacy that's already parsed through `middleware.GuardPharmacyDetailByOwner`
	pharmacy, err := middleware.PharmacyDetailFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	// 2. Create or Update (Upsert) pharmacy code in pharmacy_codes table
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
		Code      string    `json:"code"`
		ExpiresAt time.Time `json:"expires_at"`
	}{
		Code:      pharmacyCode.Code,
		ExpiresAt: pharmacyCode.ExpiresAt.Time,
	}
	if err := json.ResponseCreated(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *pharmacyHandler) GetCodeByPharmacy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pharmacy, err := middleware.PharmacyDetailFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	code, err := h.repo.Pharmacies.GetCodeByID(ctx, pharmacy.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			json.ResponseNotFound(w, r, err)
		} else {
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}

	response := struct {
		Code      string `json:"code"`
		ExpiresAt any    `json:"expires_at"`
	}{
		Code:      code.Code,
		ExpiresAt: code.ExpiresAt,
	}
	if err := json.ResponseOK(w, response); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *pharmacyHandler) GetLanding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authPharmacy, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	usersJSON := []UserJSON{}
	for _, user := range authPharmacy.Users {
		item := UserJSON{
			ID:       &user.ID,
			Username: user.Username,
		}
		usersJSON = append(usersJSON, item)
	}

	res := struct {
		Name  string     `json:"name"`
		Users []UserJSON `json:"users"`
	}{
		Name:  authPharmacy.Name,
		Users: usersJSON,
	}
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *pharmacyHandler) GetDetailByAppID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// appID := chi.URLParam(r, "appID")

	// authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	// if err != nil {
	// 	json.ResponseInternalServerError(w, r, err)
	// 	return
	// }

	// pharmacy, err := h.repo.Pharmacies.GetByAppIDForOwner(ctx, appID, authOwner.ID)
	// if err != nil {
	// 	if errors.Is(err, repository.ErrNotFound) {
	// 		json.ResponseNotFound(w, r, err)
	// 	} else {
	// 		json.ResponseInternalServerError(w, r, err)
	// 	}
	// 	return
	// }

	pharmacy, err := middleware.PharmacyDetailFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	users, err := h.repo.Users.GetByPharmacyID(ctx, pharmacy.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			json.ResponseNotFound(w, r, err)
		} else {
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}

	usersJSON := []UserJSON{}
	for _, user := range *users {
		item := UserJSON{
			ID:       &user.ID,
			Username: user.Username,
		}
		usersJSON = append(usersJSON, item)
	}

	pharmacyJSON := PharmacyJSON{
		AppID:     pharmacy.AppID,
		Name:      pharmacy.Name,
		Address:   pharmacy.Address,
		CreatedAt: pharmacy.CreatedAt,
		UpdatedAt: pharmacy.UpdatedAt,
	}

	// fmt.Println(pharmacy.ID)

	// apotekCode, err := h.repo.Pharmacies.GetCodeByID(ctx, pharmacy.ID)
	// if err != nil {
	// 	if errors.Is(err, repository.ErrNotFound) {
	// 		json.ResponseBadRequest(w, r, fmt.Errorf("%w: no active codes matched for this Pharmacy", err))
	// 	} else {
	// 		json.ResponseInternalServerError(w, r, err)
	// 	}
	// 	return
	// }

	res := struct {
		PharmacyJSON
		Users []UserJSON `json:"users"`
		// ApotekCode ApotekCodeJSON `json:"apotek_code"`
	}{
		PharmacyJSON: pharmacyJSON,
		Users:        usersJSON,
		// ApotekCode: ApotekCodeJSON{
		// 	Code:      apotekCode.Code,
		// 	ExpiresAt: apotekCode.ExpiresAt,
		// },
	}

	if err := json.ResponseOK(w, res); err != nil {
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
