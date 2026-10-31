package gopdf

import (
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

// PdfInfo Document Information Dictionary
type PdfInfo struct {
	Title        string    // The document’s title
	Author       string    // The name of the person who created the document
	Subject      string    // The subject of the document
	Creator      string    // If the document was converted to PDF from another format, the name of the application which created the original document
	Producer     string    // If the document was converted to PDF from another format, the name of the application that converted the original document to PDF
	CreationDate time.Time // The date and time the document was created, in human-readable form
}

type infoObj struct {
	getRoot func() *GoPdf
}

func (i *infoObj) init(funcGetRoot func() *GoPdf) {
	i.getRoot = funcGetRoot
}

func (i *infoObj) getType() string {
	return "Info"
}

func (i *infoObj) write(w io.Writer, objID int) error {
	var zerotime time.Time
	root := i.getRoot()
	info := root.info
	protection := root.protection()
	io.WriteString(w, "<<\n")
	if info != nil {
		if info.Author != "" {
			if err := writeInfoText(w, "Author", info.Author, objID, protection); err != nil {
				return err
			}
		}

		if info.Title != "" {
			if err := writeInfoText(w, "Title", info.Title, objID, protection); err != nil {
				return err
			}
		}

		if info.Subject != "" {
			if err := writeInfoText(w, "Subject", info.Subject, objID, protection); err != nil {
				return err
			}
		}

		if info.Creator != "" {
			if err := writeInfoText(w, "Creator", info.Creator, objID, protection); err != nil {
				return err
			}
		}

		if info.Producer != "" {
			if err := writeInfoText(w, "Producer", info.Producer, objID, protection); err != nil {
				return err
			}
		}

		if !zerotime.Equal(info.CreationDate) {
			dateStr := "D:" + infodate(info.CreationDate)
			if protection != nil {
				enc, err := rc4Cip(protection.objectkey(objID), []byte(dateStr))
				if err != nil {
					return err
				}
				fmt.Fprintf(w, "/CreationDate <%X>\n", enc)
			} else {
				fmt.Fprintf(w, "/CreationDate(%s)\n", dateStr)
			}
		}
	}
	io.WriteString(w, ">>\n")
	return nil
}

func writeInfoText(w io.Writer, key, val string, objID int, protection *PDFProtection) error {
	hexStr := "FEFF" + encodeUtf8(val)
	if protection != nil {
		raw, err := hex.DecodeString(hexStr)
		if err != nil {
			return err
		}
		enc, err := rc4Cip(protection.objectkey(objID), raw)
		if err != nil {
			return err
		}
		hexStr = fmt.Sprintf("%X", enc)
	}
	_, err := fmt.Fprintf(w, "/%s <%s>\n", key, hexStr)
	return err
}
