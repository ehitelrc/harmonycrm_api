package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

//
// ────────────────────────────────────────────────────────────────
//                      DATE NORMALIZATION
// ────────────────────────────────────────────────────────────────
//

// normalizeDate convierte fechas como:
//   - "04 diciembre 2025"
//   - "4 dic 2025"
//   - "04/12/2025"
//   - "2025-12-04"
//   - etc.
//
// a formato estándar:  yyyy/MM/dd
func NormalizeDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return ""
	}

	// 1) Probar formatos numéricos directamente
	numericLayouts := []string{
		"02/01/2006",
		"2/1/2006",
		"02-01-2006",
		"2-1-2006",
		"2006-01-02",
		"2006/01/02",
	}

	for _, layout := range numericLayouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("2006/01/02")
		}
	}

	// 2) Probar fechas con nombre de mes en español
	//    Ej: "04 diciembre 2025", "4 Dic 2025", etc.
	re := regexp.MustCompile(`(?i)(\d{1,2})\s+([A-Za-zÁÉÍÓÚÜáéíóúüñÑ]+)\s+(\d{4})`)
	match := re.FindStringSubmatch(dateStr)
	if len(match) == 4 {
		day := strings.TrimSpace(match[1])
		monthName := NormalizeMonth(match[2])
		year := strings.TrimSpace(match[3])

		monthMap := map[string]string{
			"enero": "01", "febrero": "02", "marzo": "03",
			"abril": "04", "mayo": "05", "junio": "06",
			"julio": "07", "agosto": "08", "septiembre": "09",
			"setiembre": "09", "octubre": "10",
			"noviembre": "11", "diciembre": "12",
		}

		if mm, ok := monthMap[monthName]; ok {
			if len(day) == 1 {
				day = "0" + day
			}
			return fmt.Sprintf("%s/%s/%s", year, mm, day)
		}
	}

	// 3) Fallback (si ningún formato coincidió)
	return dateStr
}

// NormalizeMonth homogeneiza y limpia nombres de meses
func NormalizeMonth(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	// Opcional: remover acentos
	replacer := strings.NewReplacer(
		"á", "a",
		"é", "e",
		"í", "i",
		"ó", "o",
		"ú", "u",
		"ü", "u",
	)
	s = replacer.Replace(s)

	return s
}

//
// ────────────────────────────────────────────────────────────────
//                       TIME NORMALIZATION
// ────────────────────────────────────────────────────────────────
//

// NormalizeTime recibe valores como:
//   - "8:07 PM"
//   - "08:07 PM"
//   - "8:7 PM"   (OCR defect)
//   - "20:07"
//   - "20:07:59"
//
// y devuelve:  HH:mm  (24 horas)
//
// Además, convierte la hora a zona horaria de Costa Rica (UTC-6).
func NormalizeTime(timeStr string) string {
	timeStr = strings.TrimSpace(timeStr)
	if timeStr == "" {
		return ""
	}

	// Limpiar espacios y variaciones del OCR
	timeStr = regexp.MustCompile(`\s+`).ReplaceAllString(timeStr, " ")

	// Intentar parsear con layouts conocidos
	layouts := []string{
		"3:04 PM",
		"03:04 PM",
		"3:04 pm",
		"03:04 pm",
		"3:04PM",
		"03:04PM",
		"3:04pm",
		"03:04pm",
		"15:04",
		"15:04:05",
	}

	// Intento directo con layouts
	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			// Convertir a zona horaria de Costa Rica (UTC-6)
			loc, _ := time.LoadLocation("America/Costa_Rica")
			tCR := t.In(loc)
			return tCR.Format("15:04")
		}
	}

	// Caso específico: OCR deja "8:7 PM" -> lo corregimos
	re := regexp.MustCompile(`(?i)^(\d{1,2}):(\d{1})\s?(AM|PM|am|pm)$`)
	match := re.FindStringSubmatch(timeStr)
	if len(match) == 4 {
		h := match[1]
		m := "0" + match[2]
		ampm := match[3]
		fixed := fmt.Sprintf("%s:%s %s", h, m, ampm)

		for _, layout := range layouts {
			if t, err := time.Parse(layout, fixed); err == nil {
				loc, _ := time.LoadLocation("America/Costa_Rica")
				tCR := t.In(loc)
				return tCR.Format("15:04")
			}
		}
	}

	// Fallback
	return timeStr
}
