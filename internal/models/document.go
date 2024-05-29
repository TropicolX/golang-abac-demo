package models

import "errors"

type Document struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Classification string `json:"classification"`
	OwnerID        string `json:"owner_id"`
}

var documents = []Document{}

func AddDocument(doc Document) {
	documents = append(documents, doc)
}

func GetDocumentByID(id string) (Document, error) {
	for _, doc := range documents {
		if doc.ID == id {
			return doc, nil
		}
	}
	return Document{}, errors.New("document not found")
}

func UpdateDocument(updatedDoc Document) {
	for i, doc := range documents {
		if doc.ID == updatedDoc.ID {
			documents[i] = updatedDoc
			return
		}
	}
}

func DeleteDocument(id string) {
	for i, doc := range documents {
		if doc.ID == id {
			documents = append(documents[:i], documents[i+1:]...)
			return
		}
	}
}
