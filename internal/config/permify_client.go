package config

import (
	"context"
	"log"

	v1 "github.com/Permify/permify-go/generated/base/v1"
	permify "github.com/Permify/permify-go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
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

func ConvertStringToAny(s string) *anypb.Any {
	anyValue, err := anypb.New(&v1.StringValue{Data: s})
	if err != nil {
		log.Fatalf("Failed to create Any from string: %v", err)
	}
	return anyValue
}

func WritePermifySchemaAndRelationships() {
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

	// Write relationships
	tuples := []*v1.Tuple{
		{
			Entity: &v1.Entity{
				Type: "document",
				Id:   "1",
			},
			Relation: "owner",
			Subject: &v1.Subject{
				Type: "user",
				Id:   "1",
			},
		},
		{
			Entity: &v1.Entity{
				Type: "document",
				Id:   "2",
			},
			Relation: "owner",
			Subject: &v1.Subject{
				Type: "user",
				Id:   "2",
			},
		},
		{
			Entity: &v1.Entity{
				Type: "document",
				Id:   "3",
			},
			Relation: "owner",
			Subject: &v1.Subject{
				Type: "user",
				Id:   "3",
			},
		},
	}

	// Write attributes
	attributes := []*v1.Attribute{
		{
			Entity: &v1.Entity{
				Type: "document",
				Id:   "1",
			},
			Attribute: "classification",
			Value:     ConvertStringToAny("public"),
		},
		{
			Entity: &v1.Entity{
				Type: "document",
				Id:   "1",
			},
			Attribute: "department",
			Value:     ConvertStringToAny("IT"),
		},
		{
			Entity: &v1.Entity{
				Type: "document",
				Id:   "2",
			},
			Attribute: "classification",
			Value:     ConvertStringToAny("internal"),
		},
		{
			Entity: &v1.Entity{
				Type: "document",
				Id:   "2",
			},
			Attribute: "department",
			Value:     ConvertStringToAny("HR"),
		},
		{
			Entity: &v1.Entity{
				Type: "document",
				Id:   "3",
			},
			Attribute: "classification",
			Value:     ConvertStringToAny("confidential"),
		},
		{
			Entity: &v1.Entity{
				Type: "document",
				Id:   "3",
			},
			Attribute: "department",
			Value:     ConvertStringToAny("Sales"),
		},
	}

	rr, err := PermifyClient.Data.Write(context.Background(), &v1.DataWriteRequest{
		TenantId: "t1",
		Metadata: &v1.DataWriteRequestMetadata{
			SchemaVersion: sr.SchemaVersion,
		},
		Tuples:     tuples,
		Attributes: attributes,
	})

	if err != nil {
		log.Fatalf("Failed to write relationships and attributes: %v", err)
	}

	SnapToken = rr.SnapToken

	log.Printf("Data tuples written successfully\nSnap token: %s", SnapToken)
}
