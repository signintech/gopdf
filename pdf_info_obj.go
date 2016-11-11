package gopdf

import "time"

//PdfInfo Document Information Dictionary
type PdfInfo struct {
	Title        string    //The document’s title
	Author       string    //The name of the person who created the document.
	Subject      string    //The subject of the document.
	Creator      string    // If the document was converted to PDF from another format, the name of the application original document from which it was converted.
	Producer     string    //If the document was converted to PDF from another format, the name of the application (for example, Acrobat Distiller) that converted it to PDF.
	CreationDate time.Time //The date and time the document was created, in human-readable form
}
