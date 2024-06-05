package middlewares

import (
	"context"
	"log"
	"net/http"

	"golang-abac-demo/internal/config"
	"golang-abac-demo/internal/controllers"
	"golang-abac-demo/internal/models"

	v1 "github.com/Permify/permify-go/generated/base/v1"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/types/known/structpb"
)

func ABACMiddleware(permission string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(controllers.UserKey).(*models.Claims)
			username := claims.Username
			user, err := models.GetUserByUsername(username)

			if err != nil {
				log.Printf("Failed to fetch user: %v", err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			vars := mux.Vars(r)
			documentID := vars["id"]

			// check if document exists
			_, docErr := models.GetDocumentByID(documentID)

			if docErr != nil {
				log.Printf("Failed to fetch document: %v", docErr)
				http.Error(w, "Document not found", http.StatusNotFound)
				return
			}

			data := map[string]interface{}{
				"dept": user.Department,
			}

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

			log.Printf("Permission denied")
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
