package middlewares

import (
	"context"
	"log"
	"net/http"

	"golang-abac-demo/internal/config"
	"golang-abac-demo/internal/models"

	v1 "github.com/Permify/permify-go/generated/base/v1"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/types/known/structpb"
)

func ABACMiddleware(permission string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get claims from context
			claims := r.Context().Value(UserKey).(*models.Claims)

			log.Printf("claims: %v", claims)
			username := claims.Username
			user, err := models.GetUserByUsername(username)

			if err != nil {
				// log error message
				log.Printf("Failed to get user by ID: %v", err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			vars := mux.Vars(r)
			documentID := vars["id"]

			data := map[string]interface{}{
				"dept": user.Department,
			}

			// TODO: remove this
			log.Printf("userId: %s, documentId: %s", username, documentID)

			// Convert map[string]interface{} to *structpb.Struct
			structData, err := structpb.NewStruct(data)
			if err != nil {
				log.Fatalf("Failed to create protobuf struct: %v", err)
			}

			cr, err := config.PermifyClient.Permission.Check(context.Background(), &v1.PermissionCheckRequest{
				TenantId: "t1",
				Metadata: &v1.PermissionCheckRequestMetadata{
					SnapToken: "v1",
					Depth:     50,
				},
				Entity: &v1.Entity{
					Type: "document",
					Id:   documentID,
				},
				Permission: permission,
				Subject: &v1.Subject{
					Type: "user",
					Id:   user.ID,
				},
				Context: &v1.Context{
					Data: structData,
				},
			})

			if err != nil {
				log.Printf("Failed to check permission: %v", err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if cr.Can == v1.CheckResult_CHECK_RESULT_ALLOWED {
				next.ServeHTTP(w, r)
				return
			}

			// descriptive error message
			log.Printf("Permission denied")
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
