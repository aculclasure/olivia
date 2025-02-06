package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/olivia-ai/olivia/locales"
	"github.com/olivia-ai/olivia/modules"

	"github.com/olivia-ai/olivia/modules/start"

	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	"github.com/olivia-ai/olivia/analysis"
	"github.com/olivia-ai/olivia/user"
	"github.com/olivia-ai/olivia/util"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// RequestMessage is the structure that uses entry connections to chat with the websocket
type RequestMessage struct {
	Type        int              `json:"type"` // 0 for handshakes and 1 for messages
	Content     string           `json:"content"`
	Token       string           `json:"user_token"`
	Locale      string           `json:"locale"`
	Information user.Information `json:"information"`
}

// ResponseMessage is the structure used to reply to the user through the websocket
type ResponseMessage struct {
	Content     string           `json:"content"`
	Tag         string           `json:"tag"`
	Information user.Information `json:"information"`
}

type webSocketSession struct {
	conn *websocket.Conn
}

func (ws *webSocketSession) Get() ([]byte, error) {
	msgType, msgBytes, err := ws.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if msgType != websocket.TextMessage {
		return nil, errors.New("websocket message type must be text message")
	}
	return msgBytes, nil
}

func (ws *webSocketSession) Send(msg []byte) error {
	log.Printf("sending message %q over the wire\n", string(msg))
	err := ws.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Printf("error in web socket session: %s\n", err)
		return err
	}
	return nil
}

// SocketHandle manages the entry connections and reply with the neural network
func SocketHandle(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	fmt.Println(color.FgGreen.Render("A new connection has been opened"))

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Unmarshal the json content of the message
		var request RequestMessage
		if err = json.Unmarshal(msg, &request); err != nil {
			continue
		}

		// Set the information from the client into the cache
		if reflect.DeepEqual(user.GetUserInformation(request.Token), user.Information{}) {
			user.SetUserInformation(request.Token, request.Information)
		}

		// If the type of requests is a handshake then execute the start modules
		if request.Type == 0 {
			start.ExecuteModules(request.Token, request.Locale)

			message := start.GetMessage()
			if message != "" {
				// Generate the response to send to the user
				response := ResponseMessage{
					Content:     message,
					Tag:         "start module",
					Information: user.GetUserInformation(request.Token),
				}

				bytes, err := json.Marshal(response)
				if err != nil {
					panic(err)
				}

				if err = conn.WriteMessage(msgType, bytes); err != nil {
					continue
				}
			}

			continue
		}

		// Write message back to browser
		// response := Reply(request)
		// if err = conn.WriteMessage(msgType, response); err != nil {
		// 	continue
		// }
		if err = process(request, conn); err != nil {
			fmt.Printf("error processing request %+v: %s", request, err)
			continue
		}
	}
}

func process(req RequestMessage, conn *websocket.Conn) error {
	if len(req.Content) > 500 {
		resp := ResponseMessage{
			Content:     util.GetMessage(req.Locale, "too long"),
			Tag:         "too long",
			Information: user.GetUserInformation(req.Token),
		}
		bytes, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
			return err
		}
		return nil
	}

	locale := req.Locale
	if !locales.Exists(locale) {
		locale = "en"
	}
	respTag := analysis.NewSentence(locale, req.Content).PredictTagWithCache(*cache, neuralNetworks[locale], req.Token)
	err := modules.RunModule(req.Content, locale, respTag, req.Token, &webSocketSession{conn})
	if err != nil {
		return err
	}
	return nil
}

// Reply takes the entry message and returns an array of bytes for the answer
func Reply(request RequestMessage) []byte {
	var responseSentence, responseTag string

	// Send a message from res/datasets/messages.json if it is too long
	if len(request.Content) > 500 {
		responseTag = "too long"
		responseSentence = util.GetMessage(request.Locale, responseTag)
	} else {
		// If the given locale is not supported yet, set english
		locale := request.Locale
		if !locales.Exists(locale) {
			locale = "en"
		}

		responseTag, responseSentence = analysis.NewSentence(
			locale, request.Content,
		).Calculate(*cache, neuralNetworks[locale], request.Token)
	}

	// Marshall the response in json
	response := ResponseMessage{
		Content:     responseSentence,
		Tag:         responseTag,
		Information: user.GetUserInformation(request.Token),
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return bytes
}
