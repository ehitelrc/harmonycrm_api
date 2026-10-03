package services

import (
	"testing"
)

func TestIsCandidateReceipt(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			name:     "Empty string",
			text:     "",
			expected: false,
		},
		{
			name:     "Short text",
			text:     "Hola",
			expected: false,
		},
		{
			name:     "Text without digits",
			text:     "Buenos días, deseo consultar el estado de mi solicitud por favor.",
			expected: false,
		},
		{
			name: "BAC Sinpe Movil real",
			text: `BAC
Notificación de transferencia SINPE
Móvil
Hola,
Le informamos que ANA MARIA realizó una transferencia por medio
de SINPE Móvil al teléfono N° 88504000
Referencia: 2026021610283000788483335
Monto: 14,396.00`,
			expected: true,
		},
		{
			name: "BCR Sinpe Movil real",
			text: `BCR
Transacción SINPE MÓVIL
Número de referencia: 2026020615283008820212999
Teléfono Destino: 88504000
Monto: 100.00`,
			expected: true,
		},
		{
			name: "BNCR Transferencia real",
			text: `Banco Nacional
Comprobante de transaccion
Número de comprobante: 29213687
Monto debitado: 5,600.00 Colones`,
			expected: true,
		},
		{
			name: "Voucher Casio / Boleta física",
			text: `No 182761
045
CASIO
¢5,000`,
			expected: true,
		},
		{
			name: "IBAN Transfer",
			text: `Transferencia a cuenta CR33015109920010324425
Monto: 25000`,
			expected: true,
		},
		{
			name:     "Meme / greeting image",
			text:     "Feliz cumpleaños que Dios te bendiga mucho hoy y siempre",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCandidateReceipt(tt.text)
			if got != tt.expected {
				t.Errorf("IsCandidateReceipt() for %q = %v, expected %v", tt.name, got, tt.expected)
			}
		})
	}
}
