package controllers

import (
	"encoding/json"
	"golang-abac-demo/internal/models"
	"golang-abac-demo/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func UploadDocument(w http.ResponseWriter, r *http.Request) {
	var doc models.Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := r.Context().Value(UserKey).(*models.Claims)
	doc.OwnerID = user.Id                            // Set document owner as the user who uploaded it
	doc.ID = string(rune(len(models.Documents) + 1)) // Generate document ID

	// Add document to repository (this would be replaced with actual DB call)
	models.AddDocument(doc)

	utils.InfoLogger.Printf("User '%s' uploaded document %s", user.Username, doc.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Document uploaded successfully"})
}

func ViewDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	docID := params["id"]
	user := r.Context().Value(UserKey).(*models.Claims)

	// Fetch document from repository (this would be replaced with actual DB call)
	doc, err := models.GetDocumentByID(docID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	utils.InfoLogger.Printf("User '%s' viewed document %s", user.Username, doc.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(doc)
}

func EditDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	docID := params["id"]
	user := r.Context().Value(UserKey).(*models.Claims)

	var doc models.Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fetch document from repository (this would be replaced with actual DB call)
	existingDoc, err := models.GetDocumentByID(docID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Update document in repository (this would be replaced with actual DB call)
	existingDoc.Title = doc.Title
	existingDoc.Content = doc.Content
	models.UpdateDocument(existingDoc)

	utils.InfoLogger.Printf("User '%s' edited document %s", user.Username, doc.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Document edited successfully"})
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	docID := params["id"]
	user := r.Context().Value(UserKey).(*models.Claims)

	// Fetch document from repository (this would be replaced with actual DB call)
	doc, err := models.GetDocumentByID(docID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete document from repository (this would be replaced with actual DB call)
	models.DeleteDocument(docID)

	utils.InfoLogger.Printf("User '%s' deleted document %s", user.Username, doc.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Document deleted successfully"})
}
