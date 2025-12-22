package gui

import (
	"fmt"
	"kituxlsx/domain"
	"kituxlsx/pdfproc"
	"log"
)

func (g *gui) pdf(out string, codes []*domain.Pallete) error {
	pdf, err := pdfproc.New(g)
	if err != nil {
		log.Fatalf("Error pdf new: %v", err)
	}
	err = pdf.BuildMaroto()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = pdf.BuildPages(codes)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = pdf.DocumentGenerate()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	fn := out + ".pdf"
	err = pdf.PdfDocumentSave(fn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (g *gui) initPdf() error {
	return nil
}
