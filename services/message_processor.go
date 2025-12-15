package services

import (
	"encoding/json"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/ws"
	"strconv"
)

// type MessageProcessor struct {
// 	hub *ws.Hub
// }

type MessageProcessor struct {
	hub        *ws.Hub
	msgRepo    repository.MessageRepository
	receiptSvc *ReceiptAnalysisService
}

// func NewMessageProcessor(hub *ws.Hub) *MessageProcessor {
// 	return &MessageProcessor{hub: hub}
// }

func NewMessageProcessor(
	hub *ws.Hub,
	receiptSvc *ReceiptAnalysisService,
) *MessageProcessor {
	return &MessageProcessor{
		hub:        hub,
		msgRepo:    repository.MessageRepository{},
		receiptSvc: receiptSvc,
	}
}

// func (p *MessageProcessor) ProcessIncomingMessage(input models.IncomingMessage) (*models.Message, error) {

// 	repo := repository.MessageRepository{}

// 	newMessage, err := repo.CreateMessage(input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// WS broadcast (aquí vive TODO lo de tiempo real)
// 	if newMessage.CaseID != 0 && p.hub != nil {
// 		payload, _ := json.Marshal(map[string]interface{}{
// 			"type":    "new_message",
// 			"case_id": newMessage.CaseID,
// 			"data":    newMessage,
// 		})

// 		channel := "case:" + strconv.Itoa(int(newMessage.CaseID))
// 		p.hub.BroadcastJSON(channel, payload)

// 		if newMessage.AgentID != nil {
// 			p.hub.BroadcastJSON(
// 				"agent:"+strconv.Itoa(int(*newMessage.AgentID)),
// 				payload,
// 			)
// 		}
// 	}

// 	return newMessage, nil
// }

func (p *MessageProcessor) ProcessIncomingMessage(
	input models.IncomingMessage,
) (*models.Message, error) {

	// 1️⃣ Guardar mensaje (ÚNICO lugar)
	newMessage, err := p.msgRepo.CreateMessage(input)
	if err != nil {
		return nil, err
	}

	// 2️⃣ Emitir WS (ÚNICO lugar)
	p.broadcast(newMessage)

	// 3️⃣ Side-effects async
	switch input.MessageType {
	case "image":
		go p.processImage(input, newMessage)
	case "audio":
		go p.processAudio(input, newMessage)
	}

	return newMessage, nil
}

func (p *MessageProcessor) broadcast(msg *models.Message) {
	if p.hub == nil {
		return
	}

	// Debe existir un caso para emitir
	if msg.CaseID == 0 {
		return
	}

	payload, err := json.Marshal(map[string]interface{}{
		"type":    "new_message",
		"case_id": msg.CaseID,
		"data":    msg,
	})
	if err != nil {
		return
	}

	// Canal por caso
	p.hub.BroadcastJSON(
		"case:"+strconv.Itoa(int(msg.CaseID)),
		payload,
	)

	// Canal por agente (si aplica)
	if msg.AgentID != nil {
		p.hub.BroadcastJSON(
			"agent:"+strconv.Itoa(int(*msg.AgentID)),
			payload,
		)
	}
}
