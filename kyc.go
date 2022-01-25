package mango

import (
	"encoding/base64"
	"errors"
)

type DocumentType string

const (
	IdentityProof          DocumentType = "IDENTITY_PROOF"
	RegistrationProof      DocumentType = "REGISTRATION_PROOF"
	ArticlesOfAssociation  DocumentType = "ARTICLES_OF_ASSOCIATION"
	ShareholderDeclaration DocumentType = "SHAREHOLDER_DECLARATION"
	AddressProof           DocumentType = "ADDRESS_PROOF"
)

type DocumentStatus string

const (
	DocumentStatusCreated         DocumentStatus = "CREATED"
	DocumentStatusValidationAsked DocumentStatus = "VALIDATION_ASKED"
	DocumentStatusValidated       DocumentStatus = "VALIDATED"
	DocumentStatusRefused         DocumentStatus = "REFUSED"
	DocumentStatusOutOfDate       DocumentStatus = "OUT_OF_DATE"
)

type DocumentRefusedReasonType string

const (
	DocumentRefusedReasonTypeUnreadable          DocumentRefusedReasonType = "DOCUMENT_UNREADABLE"
	DocumentRefusedReasonTypeNotAccepted         DocumentRefusedReasonType = "DOCUMENT_NOT_ACCEPTED"
	DocumentRefusedReasonTypeHasExpired          DocumentRefusedReasonType = "DOCUMENT_HAS_EXPIRED"
	DocumentRefusedReasonTypeIncomplete          DocumentRefusedReasonType = "DOCUMENT_INCOMPLETE"
	DocumentRefusedReasonTypeNotMatchUserData    DocumentRefusedReasonType = "DOCUMENT_DO_NOT_MATCH_USER_DATA"
	DocumentRefusedReasonTypeNotMatchAccountData DocumentRefusedReasonType = "DOCUMENT_DO_NOT_MATCH_ACCOUNT_DATA"
	DocumentRefusedReasonTypeFalsified           DocumentRefusedReasonType = "DOCUMENT_FALSIFIED"
	DocumentRefusedReasonTypeUnderagePerson      DocumentRefusedReasonType = "UNDERAGE PERSON"
	DocumentRefusedReasonTypeSpecificCase        DocumentRefusedReasonType = "SPECIFIC_CASE"
)

func (m *MangoPay) Document(id string) (*Document, error) {
	any, err := m.anyRequest(new(Document), actionFetchKYCDocument, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return any.(*Document), nil
}

func (m *MangoPay) NewDocument(user Consumer, docType DocumentType, tag string) (*Document, error) {
	id := consumerId(user)
	if id == "" {
		return nil, errors.New("user has empty Id")
	}
	data := JsonObject{
		"UserId": id,
		"Type":   docType,
	}
	if len(tag) > 0 {
		data["Tag"] = tag
	}

	doc, err := m.anyRequest(new(Document), actionCreateKYCDocument, data)
	if err != nil {
		return nil, err
	}
	casted := doc.(*Document)
	casted.service = m
	return casted, nil
}

func (m *MangoPay) Documents(user Consumer, status DocumentStatus) (DocumentList, error) {
	data := JsonObject{}
	action := actionFetchAllKYCDocuments
	if user != nil {
		id := consumerId(user)
		if id == "" {
			return nil, errors.New("user has empty Id")
		}
		data["UserId"] = id
		data["Status"] = status
		action = actionFetchUserKYCDocuments
	}

	list, err := m.anyRequest(new(DocumentList), action, data)
	if err != nil {
		return nil, err
	}
	casted := *(list.(*DocumentList))
	for _, doc := range casted {
		doc.service = m
	}
	return casted, nil
}

type DocumentList []*Document

type Document struct {
	ProcessIdent
	UserId               string
	Status               DocumentStatus
	Type                 DocumentType
	RefusedReasonMessage string
	RefusedReasonType    DocumentRefusedReasonType

	service *MangoPay
}

func (d *Document) Submit(status DocumentStatus, tag string) error {
	data := JsonObject{
		"Id":     d.Id,
		"UserId": d.UserId,
		"Status": status,
	}
	if len(tag) > 0 {
		data["Tag"] = tag
	}

	doc, err := d.service.anyRequest(new(Document), actionSubmitKYCDocument, data)
	if err != nil {
		return err
	}
	casted := doc.(*Document)
	casted.service = d.service
	*d = *casted
	return nil
}

func (d *Document) CreatePage(file []byte) error {
	data := JsonObject{
		"Id":     d.Id,
		"UserId": d.UserId,
		"File":   base64.StdEncoding.EncodeToString(file),
	}

	_, err := d.service.anyRequest(new(JsonObject), actionCreateKYCPage, data)
	return err
}
