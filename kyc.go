package mango

import (
	"encoding/base64"
	"errors"
)

func (m *MangoPay) Document(id string) (*Document, error) {
	any, err := m.anyRequest(new(Document), actionFetchKYCDocument, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return any.(*Document), nil
}

func (m *MangoPay) NewDocument(user Consumer, docType, tag string) (*Document, error) {
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

func (m *MangoPay) Documents(user Consumer) (DocumentList, error) {
	data := JsonObject{}
	action := actionFetchAllKYCDocuments
	if user != nil {
		id := consumerId(user)
		if id == "" {
			return nil, errors.New("user has empty Id")
		}
		data["UserId"] = id
		action = actionFetchUserKYCDocuments
	}

	list, err := m.anyRequest(new(DocumentList), action, data)
	if err != nil {
		return nil, err
	}
	return list.(DocumentList), nil
}

type DocumentList []*Document

type Document struct {
	ProcessIdent
	UserId               string
	Status               string
	Type                 string
	RefusedReasonMessage string
	RefusedReasonType    string

	service *MangoPay
}

func (d *Document) Submit(status, tag string) error {
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

	_, err := d.service.anyRequest(new(JsonObject), actionSubmitKYCDocument, data)
	return err
}
