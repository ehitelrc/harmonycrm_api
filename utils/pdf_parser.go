package utils

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// ConvertHTMLToPDFChromedp recibe html en []byte y devuelve base64 del PDF
func ConvertHTMLToPDFChromedp(html []byte) (*string, error) {

	encodedHTML := base64.StdEncoding.EncodeToString(html)
	dataURL := "data:text/html;base64," + encodedHTML

	// Crear contexto Chrome
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Timeout para evitar cuelgues
	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var pdfBytes []byte

	// Ejecutar Chrome
	err := chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.Sleep(300*time.Millisecond),

		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0.5).
				WithMarginLeft(0.5).
				WithMarginRight(0.5).
				WithMarginBottom(0.5).
				Do(ctx)

			if err != nil {
				return err
			}

			pdfBytes = buf
			return nil
		}),
	)

	if err != nil {
		return nil, err
	}

	// Convertimos a base64 porque tu estructura lo usa así
	b64 := base64.StdEncoding.EncodeToString(pdfBytes)

	return &b64, nil
}
