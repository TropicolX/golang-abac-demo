package controllers

import (
	"encoding/json"
	"golang-abac-demo/internal/models"
	"golang-abac-demo/internal/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func UploadDocument(w http.ResponseWriter, r *http.Request) {
	var doc models.Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(UserKey).(*models.Claims)
	user, err := models.GetUserByUsername(claims.Username)
	if err != nil {
		log.Printf("Failed to fetch user: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	doc.OwnerID = user.ID
	noOfDocs := len(models.Documents)

	if noOfDocs == 0 {
		doc.ID = "1"
	} else {
		lastDocIndex := noOfDocs - 1
		lastDocID, err := strconv.Atoi(models.Documents[lastDocIndex].ID)

		if err != nil {
			log.Printf("Failed to fetch last document ID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		doc.ID = strconv.Itoa(lastDocID + 1)
	}

	// Add document to repository (this would be replaced with actual DB call)
	models.AddDocument(doc)

	utils.InfoLogger.Printf("User '%s' uploaded document %s", user.Username, doc.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Document uploaded successfully"})
}

func ViewDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	docID := params["id"]
	claims := r.Context().Value(UserKey).(*models.Claims)

	// Fetch document from repository (this would be replaced with actual DB call)
	doc, err := models.GetDocumentByID(docID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	utils.InfoLogger.Printf("User '%s' viewed document %s", claims.Username, doc.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(doc)
}

func EditDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	docID := params["id"]
	claims := r.Context().Value(UserKey).(*models.Claims)

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
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Update document in repository (this would be replaced with actual DB call)
	existingDoc.Title = doc.Title
	existingDoc.Content = doc.Content
	models.UpdateDocument(existingDoc)

	utils.InfoLogger.Printf("User '%s' edited document %s", claims.Username, doc.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Document edited successfully"})
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	docID := params["id"]
	claims := r.Context().Value(UserKey).(*models.Claims)

	// Fetch document from repository (this would be replaced with actual DB call)
	doc, err := models.GetDocumentByID(docID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Delete document from repository (this would be replaced with actual DB call)
	models.DeleteDocument(docID)

	utils.InfoLogger.Printf("User '%s' deleted document %s", claims.Username, doc.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Document deleted successfully"})
}
