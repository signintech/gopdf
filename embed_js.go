package gopdf

import (
	"bytes"
	"fmt"
)

//EmbedJs  embed JavaScript inside the PDF
type EmbedJs struct {
	buffer                bytes.Buffer
	indexOfEmbedJsContent int
}

func (e *EmbedJs) init(func() *GoPdf) {

}
func (e *EmbedJs) getType() string {
	return "EmbedJs"
}
func (e *EmbedJs) getObjBuff() *bytes.Buffer {
	return &e.buffer
}

func (e *EmbedJs) build(objID int) error {
	e.buffer.WriteString("<<\n")
	e.buffer.WriteString(fmt.Sprintf(" /Names [(EmbeddedJS) %d 0 R] \n", e.indexOfEmbedJsContent+1))
	e.buffer.WriteString(">>\n")
	return nil
}

func (e *EmbedJs) setIndexOfEmbedJsContent(index int) {
	e.indexOfEmbedJsContent = index
}

//EmbedJsContent JavaScript content
type EmbedJsContent struct {
	buffer  bytes.Buffer
	content string
}

func (e *EmbedJsContent) init(func() *GoPdf) {

}
func (e *EmbedJsContent) getType() string {
	return "EmbedJsContent"
}
func (e *EmbedJsContent) getObjBuff() *bytes.Buffer {
	return &e.buffer
}

func (e *EmbedJsContent) build(objID int) error {
	e.buffer.WriteString("<<\n")
	e.buffer.WriteString(" /S /JavaScript\n")
	e.buffer.WriteString(fmt.Sprintf(" /JS (%s)\n", e.content))
	e.buffer.WriteString(">>\n")
	return nil
}

func (e *EmbedJsContent) setJs(content string) {
	e.content = content
}
