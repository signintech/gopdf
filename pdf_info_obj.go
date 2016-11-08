package gopdf

import "time"

type pdfInfoObj struct {
	Author       string
	Creator      string
	Producer     string
	CreationDate time.Time
}
