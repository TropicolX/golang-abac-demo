package config

import (
	"context"
	"golang-abac-demo/internal/models"
	"golang-abac-demo/internal/utils"
	"log"

	v1 "github.com/Permify/permify-go/generated/base/v1"
	permify "github.com/Permify/permify-go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var PermifyClient *permify.Client
var SchemaVersion string
var SnapToken string

func InitPermifyClient() {
	client, err := permify.NewClient(
		permify.Config{
			Endpoint: "localhost:3478",
		},
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Permify client: %v", err)
	} else {
		PermifyClient = client
		log.Println("Permify client initialized successfully")
	}
}

func WritePermifySchema() {
	// Write schema
	schema := `
		entity user {}

		entity document {
			relation owner @user
			attribute classification string
			attribute department string
				
			permission view = is_public(classification) or (is_internal(classification) and in_same_department(department, request.dept)) or (is_confidential(classification) and owner)
			permission edit = owner or (is_internal(classification) and in_same_department(department, request.dept))
			permission delete = owner
		}

		rule is_public(classification string) {
			classification == 'public'
		}

		rule is_internal(classification string) {
			classification == 'internal'
		}

		rule is_confidential(classification string) {
			classification == 'confidential'
		}

		rule in_same_department(department string, dept string) {
			department == dept
		}
  `

	sr, err := PermifyClient.Schema.Write(context.Background(), &v1.SchemaWriteRequest{
		TenantId: "t1",
		Schema:   schema,
	})

	if err != nil {
		log.Fatalf("Failed to write schema: %v", err)
	}

	SchemaVersion = sr.SchemaVersion
	log.Printf("Schema version %s written successfully", SchemaVersion)
}

func SyncPermify() {
	// Read current relationships from Permify
	rr, err := PermifyClient.Data.ReadRelationships(context.Background(), &v1.RelationshipReadRequest{
		TenantId: "t1",
		Metadata: &v1.RelationshipReadRequestMetadata{
			SnapToken: SnapToken,
		},
		Filter: &v1.TupleFilter{
			Entity: &v1.EntityFilter{
				Type: "document",
			},
			Subject: &v1.SubjectFilter{
				Type: "user",
			},
		},
	})

	if err != nil {
		log.Fatalf("Failed to read relationships from Permify: %v", err)
	}

	// Map of existing document IDs in Permify
	existingDocumentIDs := make([]string, 0)
	nonExistingDocumentIDs := make([]string, 0)

	for _, tuple := range rr.Tuples {
		if tuple.Entity.Type == "document" {
			_, err := models.GetDocumentByID(tuple.Entity.Id)

			if err != nil {
				nonExistingDocumentIDs = append(nonExistingDocumentIDs, tuple.Entity.Id)
			} else {
				existingDocumentIDs = append(existingDocumentIDs, tuple.Entity.Id)
			}
		}
	}

	// Delete documents that don't exist in the database
	if len(nonExistingDocumentIDs) > 0 {
		rr, err := PermifyClient.Data.Delete(context.Background(), &v1.DataDeleteRequest{
			TenantId: "t1",
			TupleFilter: &v1.TupleFilter{
				Entity: &v1.EntityFilter{
					Type: "document",
					Ids:  nonExistingDocumentIDs,
				},
			},
			AttributeFilter: &v1.AttributeFilter{
				Entity: &v1.EntityFilter{
					Type: "document",
					Ids:  nonExistingDocumentIDs,
				},
				Attributes: []string{"classification", "department"},
			},
		})

		if err != nil {
			log.Fatalf("Failed to delete orphaned documents from Permify: %v", err)
		}

		SnapToken = rr.SnapToken
		log.Printf("Orphaned documents deleted from Permify successfully\nSnap token: %s", SnapToken)

	} else {
		log.Println("No orphaned documents to delete from Permify")
	}

	// Add missing documents to Permify
	var tuples []*v1.Tuple
	var attributes []*v1.Attribute

	for _, doc := range models.Documents {
		if !utils.ContainsString(existingDocumentIDs, doc.ID) {
			tuples = append(tuples, &v1.Tuple{
				Entity: &v1.Entity{
					Type: "document",
					Id:   doc.ID,
				},
				Relation: "owner",
				Subject: &v1.Subject{
					Type: "user",
					Id:   doc.OwnerID,
				},
			})

			user, err := models.GetUserByID(doc.OwnerID)

			if err != nil {
				log.Fatalf("Failed to fetch user by ID: %v", err)
			}

			attributes = append(attributes, []*v1.Attribute{
				{
					Entity: &v1.Entity{
						Type: "document",
						Id:   doc.ID,
					},
					Attribute: "classification",
					Value:     utils.ConvertStringToAny(doc.Classification),
				},
				{
					Entity: &v1.Entity{
						Type: "document",
						Id:   doc.ID,
					},
					Attribute: "department",
					Value:     utils.ConvertStringToAny(user.Department),
				},
			}...)
		}
	}

	// Write missing documents to Permify
	if len(tuples) > 0 {
		rr, err := PermifyClient.Data.Write(context.Background(), &v1.DataWriteRequest{
			TenantId: "t1",
			Metadata: &v1.DataWriteRequestMetadata{
				SchemaVersion: SchemaVersion,
			},
			Tuples:     tuples,
			Attributes: attributes,
		})

		if err != nil {
			log.Fatalf("Failed to write missing documents to Permify: %v", err)
		}

		SnapToken = rr.SnapToken
		log.Printf("Missing documents written successfully\nSnap token: %s", SnapToken)
	} else {
		log.Println("No missing documents to write to Permify")
	}

	log.Println("Permify synced successfully")
}
