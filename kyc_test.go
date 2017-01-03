package mango

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestKYC(test *testing.T) {
	serv := newTestService(test)
	user := createTestUser(serv)
	if err := user.Save(); err != nil {
		test.Fatal("Unable to store user:", err)
	}
	doc, err := serv.NewDocument(user, IdentityProof, "Tag1")
	if err != nil {
		test.Fatal("Unable to create identity proof doc:", err)
	}
	if err := doc.CreatePage(newPngImageFile()); err != nil {
		test.Fatal("Unable to create document's page:", err)
	}
	if err := doc.Submit(DocumentStatusValidationAsked, "Tag2"); err != nil {
		test.Fatal("Unable to submit document")
	}
}

func newPngImageFile() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	var buffer bytes.Buffer
	err := png.Encode(&buffer, img)
	if err != nil {
		panic("Unable to create image file: " + err.Error())
	}
	return buffer.Bytes()
}
