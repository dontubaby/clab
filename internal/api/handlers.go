package server

import (
	"cyber/internal/game"
	"cyber/internal/models"
	"encoding/json"
	"fmt"
	"github.com/lxzan/gws"
	"log"
)

// WebSocketHandler реализует интерфейс gws.EventHandler.
type WebSocketHandler struct {
	actionHandlers map[string]ActionHandler
}

// OnMessage срабатывает при установке соединения.
func (h *WebSocketHandler) OnOpen(socket *gws.Conn) {
	log.Println("WebSocket connection open sucess")
}

// OnMessage срабатывает при получении сообщения по websocket.
func (h *WebSocketHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	log.Printf("Received from WebScoket: %s", message.Data.String())

	action, err := h.parseMessage(message.Bytes())
	if err != nil {
		log.Printf("cant parse frontend message to action type:%v\n", err)
		return
	}

	if result , err := h.handleAction(action); err != nil {
		log.Printf("cant handle action:%v\n", err)
		return
	}

	  // Отправляем результат клиенту
	  if err := h.sendResponse(socket, result); err != nil {
        log.Printf("failed to send response: %v", err)
    }
}

// OnClose вызывается при закрытии по websocket соединения.
func (h *WebSocketHandler) OnClose(socket *gws.Conn, err error) {
	log.Println("WebSocket connection closed")
}
//метод отправки реузльтатов обработки сервером сообщения, полученного по websocket.
func (h *WebSocketHandler) sendResponse(socket *gws.Conn, result interface{}) error {
	response,err:=json.Marshal(result)
		if err!=nil{
			fmt.Errorf("marshal to JSON error: %v\n",err)
		}
	}
	return socket.WriteMessage(gws.OpcodeText,response)
}

//parseMessage преобразует сообщение сообщение полученное по websocket в *models.Action
func (h *WebSocketHandler) parseMessage(data []byte) (*models.Action, error) {
	var action models.Action
	if err := json.Unmarshal(data, &action); err != nil {
		return nil, fmt.Errorf("frontend message unmarshal error", err)
	}
	return &action, nil
}

//основной обработчик, запускающий обработчик соответствующий полученному с фронта действию
func (h *WebSocketHandler) handleAction(action *models.Action) (interface{},error) {
	handler, ok := h.actionHandlers[action.ActionType]
	if !ok {
		return nil,fmt.Errorf("unknown action type: %s", action.ActionType)
	}
	
	result , err := handler.Handle(action)
	if err!=nil{
		return nil,err
	}

	return result,nil
}

//Конструктор для WebSocketHandler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		actionHandlers: map[string]ActionHandler{
			"move":    &MoveActionHandler{},
			"harvest": &HarvestActionHandler{},
			"build":   &BuildActionHandler{},
			"attack":  &AttackActionHandler{},
		},
	}
}

// ActionHandler представляет интерфейс для обработки действий.
type ActionHandler interface {
	Handle(action *models.Action) error
}

// Функция преобразующая байтовый срез в тип данных удовлетворяющий троебованию [T any]
//В UnmarshalCharacteristics используются дженерики в связи с тем, характеристики
// разнятся в зависимости от типа действия, для этого и используется параметризация типа
func UnmarshalCharacteristics[T any](data []byte) (*T, error) {
	var characteristics T
	if err:=json.Unmarshal(data,&characteristics);err!=nil{
		return nil, fmt.Errorf("unmarshal charachteristics error: %v\n",err)
	}
	return &characteristics, nil
}

// MoveActionHandler обрабатывает действия типа "move".
type MoveActionHandler struct{}

func (mh *MoveActionHandler) Handle(action *Action) error {
	//подаем на вход (action.Charachteristics) являющийся []byte 
	charachteristics,err:=UnmarshalCharacteristics[MoveActionCharacteristics](action.Charachteristics)
	if err!=nil{
		return return nil, fmt.Errorf("unmarshal action charachteristics error: %v\n",err)
	}
	// Теперь characteristics имеет тип *MoveActionCharacteristics и запускает метод поиска пути
	game.AStar(moveCharacteristics.From, moveCharacteristics.To, action.AreaId)
	return nil
}

// HarvestActionHandler обрабатывает действия типа "harvest".
type HarvestActionHandler struct{}

func (hh *HarvestActionHandler) Handle(action *Action) error {
	characteristics, err := UnmarshalCharacteristics[HarvestActionCharacteristics](action.Characteristics)
    if err != nil {
        return fmt.Errorf("unmarshal characteristics error: %v", err)
    }

	// TODO:Логика для harvest
	return nil
}

// BuildActionHandler обрабатывает действия типа "build".
type BuildActionHandler struct{}

func (bh *BuildActionHandler) Handle(action *Action) error {
	characteristics,err=UnmarshalCharacteristics[BuildActionCharacteristics](action.Charachteristics)
	if err != nil {
        return fmt.Errorf("unmarshal characteristics error: %v", err)
    }

	// TODO:Логика для build
	return nil
}

// AttackActionHandler обрабатывает действия типа "attack".
type AttackActionHandler struct{}

func (ah *AttackActionHandler) Handle(action *Action) error {
	characteristics,err=UnmarshalCharacteristics[AttackActionCharacteristics](action.Charachteristics)
	if err != nil {
        return fmt.Errorf("unmarshal characteristics error: %v", err)
    }

	// TODO:Логика для attack
	return nil
}
