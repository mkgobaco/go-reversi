package tcpimpl

import (
	"encoding/json"
	"fmt"
	"reversi/core"

	"github.com/google/uuid"
)

type Coordinate struct {
	X int
	Y int
}

type GameStateDAO struct {
	Player core.Player
	AvailableMoves []Coordinate
	WhiteCells []Coordinate
	BlackCells []Coordinate
}

type tcpResponder struct {
	infraResponder infraResponder
}

func (responder *tcpResponder) MoveSuccess(result core.MoveSuccessResult) {
	nextPlayInfo := result.NextPlayerInfo()

	fmt.Printf("%s can make a play\n", nextPlayInfo.NextPlayer())
}

type MoveList struct {
	Data []core.Coordinate
}

func (responder *tcpResponder) GameInitialized(nextPlay core.NextPlayInfo) {
	responder.infraResponder.RespondAll("Game initialized!")
	responder.infraResponder.NotifyActivePlayer(fmt.Sprintf("%s -> \n", nextPlay.NextPlayer()), nextPlay.NextPlayer())

	// moves := nextPlay.Moves()
	// playerTurn := nextPlay.NextPlayer()
	// whiteCells := nextPlay.WhiteCells()
	// blackCells := nextPlay.BlackCells()
	
	data, err := json.Marshal(MoveList{
		Data: nextPlay.Moves(),
	})
	if err != nil {
		fmt.Println("Got an erorr :(")
	}

	responder.infraResponder.NotifyActivePlayer(string(data), nextPlay.NextPlayer())
	responder.infraResponder.NotifyInactivePlayer(fmt.Sprintf("it is %s players turn", nextPlay.NextPlayer()), nextPlay.NextPlayer())
}
func (responder *tcpResponder) MoveFailure() {
	fmt.Println("Move failed")

	responder.infraResponder.Respond("Move Failed")
}

func newTcpResponder(infraResponder infraResponder) core.ResultHandler {
	return &tcpResponder{infraResponder: infraResponder}
}

type infraResponder struct {
	players    []*ActivePlayer
	responseId uuid.UUID
}

func (responder infraResponder) RespondAll(message string) {
	for _, player := range responder.players {
		player.Notify(message)
	}
}

func (responder infraResponder) Respond(message string) {
	for _, player := range responder.players {
		if player.ResponseId == responder.responseId {
			player.Notify(message + " -> " + responder.responseId.String())
		}
	}
}

func (responder infraResponder) NotifyActivePlayer(message string, activePlayerSide core.Player) {
	for _, player := range responder.players {
		if player.side == activePlayerSide {
			player.Notify(message + " -> " + "for player of side black")
		}
	}
}

func (responder infraResponder) NotifyInactivePlayer(message string, activePlayerSide core.Player) {
	for _, player := range responder.players {
		if player.side != activePlayerSide {
			player.Notify(message + " -> " + "for player of side white")
		}
	}
}

type responderFactory struct {
	players []*ActivePlayer
}

func (factory responderFactory) getInstance(responseId uuid.UUID) core.ResultHandler {
	infraResponder := infraResponder{
		players:    factory.players,
		responseId: responseId,
	}

	return newTcpResponder(infraResponder)
}
