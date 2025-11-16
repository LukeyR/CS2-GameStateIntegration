package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/LukeyR/CS2-GameStateIntegration/cmd/BombTimerOnStreamDeck/events"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type eventHander func(connection *StreamDeckConnection, event []byte)

type StreamDeckConnection struct {
	port      string
	conn      *websocket.Conn
	callbacks map[string][]eventHander
	Contexts  map[string]struct{}
}

func setupWebsocket(port string, registerEvent *events.RegisterEvent) *StreamDeckConnection {
	log.Info().Msgf("registerEvent: %+v", registerEvent)

	u := url.URL{Scheme: "ws", Host: fmt.Sprint("localhost:", port), Path: "/"}
	log.Info().Msgf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	sendJson, _ := json.Marshal(registerEvent)
	log.Info().Msgf("connected to %s, sending register JSON %s", u.String(), sendJson)

	err = c.WriteJSON(registerEvent)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	err = c.WriteJSON(&events.LogMessage{Event: "logMessage", Payload: events.LogMessagePayload{Message: "tstst"}})
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msgf("connected to %s, register JSON sent", u.String())

	conn := &StreamDeckConnection{
		port:      port,
		conn:      c,
		callbacks: make(map[string][]eventHander),
		Contexts:  make(map[string]struct{}),
	}

	go func() {
		for {
			var eventType events.GenericReceiveEvent
			var rawMsg bytes.Buffer
			log.Info().Msg("Waiting on message")

			_, r, err := c.NextReader()
			eventReader := io.TeeReader(r, &rawMsg)
			if err != nil {
				log.Error().Err(err).Msg("Error reading message")
				continue
			}
			err = json.NewDecoder(eventReader).Decode(&eventType)
			if err != nil {
				log.Error().Err(err).Msg("Error unmarshalling whilst reading event type")
				continue
			}
			log.Info().Str("EventType", eventType.Event).Send()
			data, err := io.ReadAll(&rawMsg)
			callbacks := conn.callbacks[eventType.Event]
			for _, cb := range callbacks {
				go cb(conn, data)
			}
		}
	}()

	return conn
}

func (c *StreamDeckConnection) CloseWebsocket() {

	log.Info().Msg("interrupt for WS received, shutting down")

	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Error().Err(err).Msg("Error sending close message")
	}
	err = c.conn.Close()
	if err != nil {
		log.Error().Msg("Error closing connection")
	}
}

func (c *StreamDeckConnection) AddEventCallback(eventName string, callback eventHander) {
	// TODO: validate event name?
	c.callbacks[eventName] = append(c.callbacks[eventName], callback)
}

// TODO: should generalise below

func (c *StreamDeckConnection) SetTitle(title string) error {
	baseEvent := &events.SetTitle{Event: "setTitle", Payload: events.SetTitlePayload{Title: title, Target: 0}}
	for k := range c.Contexts {
		baseEvent.Context = k
		sendJson, _ := json.Marshal(baseEvent)
		log.Info().Msgf("sending JSON %s", sendJson)
		err := c.conn.WriteJSON(baseEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Couldn't send setTitle `%s`", title)
		}
		log.Info().Msgf("sent JSON %s", sendJson)
	}
	return nil
}

func (c *StreamDeckConnection) SetImage(imgPath string) error {
	baseEvent := &events.SetImage{Event: "setImage", Payload: events.SetImagePayload{Image: imgPath, Target: 0}}
	for k := range c.Contexts {
		baseEvent.Context = k
		sendJson, _ := json.Marshal(baseEvent)
		log.Info().Msgf("sending JSON %s", sendJson)
		err := c.conn.WriteJSON(baseEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Couldn't send setImage `%s`", imgPath)
		}
		log.Info().Msgf("sent JSON %s", sendJson)
	}
	return nil
}

func (c *StreamDeckConnection) SetState(state int) error {
	baseEvent := &events.SetState{Event: "setState", Payload: events.SetStatePayload{State: state}}
	for k := range c.Contexts {
		baseEvent.Context = k
		sendJson, _ := json.Marshal(baseEvent)
		log.Info().Msgf("sending JSON %s", sendJson)
		err := c.conn.WriteJSON(baseEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Couldn't send setState `%d`", state)
		}
		log.Info().Msgf("sent JSON %s", sendJson)
	}
	return nil
}

func (c *StreamDeckConnection) SendLog(msg string) error {
	baseEvent := &events.LogMessage{Event: "logMessage", Payload: events.LogMessagePayload{Message: msg}}
	sendJson, _ := json.Marshal(baseEvent)
	log.Info().Msgf("sending JSON %s", sendJson)
	err := c.conn.WriteJSON(baseEvent)
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't send logMessage `%s`", msg)
	}
	return nil
}

func (c *StreamDeckConnection) ShowAlert() error {
	baseEvent := &events.ShowAlert{Event: "showAlert"}
	for k := range c.Contexts {
		baseEvent.Context = k
		sendJson, _ := json.Marshal(baseEvent)
		log.Info().Msgf("sending JSON %s", sendJson)
		err := c.conn.WriteJSON(baseEvent)
		if err != nil {
			log.Error().Err(err).Msg("Couldn't show alert")
		}
	}
	return nil
}
