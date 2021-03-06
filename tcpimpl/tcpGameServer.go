package tcpimpl

import (
	"fmt"
	"net"

	"github.com/google/uuid"
)

type PendingGame struct {
	players []PlayerConnection
}

func (pending *PendingGame) addPlayer(connection net.Conn) error {
	if len(pending.players) < 2 {
		uuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}

		playerConnection := PlayerConnection{
			Connection: connection,
			Id:         uuid,
		}
		pending.players = append(pending.players, playerConnection)
	}

	return nil
}

func (pending PendingGame) IsFull() bool {
	return len(pending.players) == 2
}

func (pending PendingGame) stillAcceptingPlayers() bool {
	return len(pending.players) < 2
}

func message(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Errorf("failed to send message: %s", err.Error())
	}
}

func (pendingGame *PendingGame) AddPlayer(playerConnection net.Conn) {
	if pendingGame.stillAcceptingPlayers() {
		addPlayerErr := pendingGame.addPlayer(playerConnection)
		if addPlayerErr != nil {
			defer playerConnection.Close()
			message(playerConnection, "failed to generate an id")
			fmt.Errorf("failed to add the player to a pending game: %s", addPlayerErr.Error())
		}
	} else {
		defer playerConnection.Close()
		message(playerConnection, "Maximum players reached")
	}
}

func NewPendingGame() PendingGame {
	var players []PlayerConnection
	return PendingGame{
		players: players,
	}
}
