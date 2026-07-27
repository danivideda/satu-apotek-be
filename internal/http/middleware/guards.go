package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/go-chi/chi/v5"
)

const pharmacyDetailCtx = "PharmacyDetailCtx"

type pharmacyDetail struct {
	dbsqlc.Pharmacy
}
// This is to guard resources from getting accessed from unauthorized entity
// e.g owner getting details from other Pharmacies that's not owned is not allowed

func (m *AppMiddleware) GuardPharmacyDetailByOwner(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authOwner, err := AuthOwnerFromCtx(ctx)
		if err != nil {
			json.ResponseInternalServerError(w, r, err)
			return
		}
		
		appID := chi.URLParam(r, "appID")
		pharmacy, err := m.repo.Pharmacies.GetByAppIDForOwner(ctx, appID, authOwner.ID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				json.ResponseNotFound(w, r, err)
			} else {
				json.ResponseInternalServerError(w, r, err)
			}
			return
		}

		pharmacyDetail := pharmacyDetail{
			Pharmacy: *pharmacy,
		}

		ctx = context.WithValue(ctx, pharmacyDetailCtx, pharmacyDetail)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	
	return http.HandlerFunc(fn)
}

func PharmacyDetailFromCtx(ctx context.Context) (*pharmacyDetail, error) {
	pharmacyDetail, ok := ctx.Value(pharmacyDetailCtx).(pharmacyDetail)
	if !ok {
		return nil, errors.New("PharmacyDetailCtx type assertion missmatch") 
	}

	return &pharmacyDetail, nil
}