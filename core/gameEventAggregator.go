package core

// should be defined in the brain probably
type EventType string

const (
	INITILIZED EventType = "INITIALIZED"
	MOVED      EventType = "MOVED"
)

type Event struct {
	EventType EventType
	Data      interface{}
}

func NewInitializedEvent() Event {
	return Event{
		EventType: INITILIZED,
		Data:      nil,
	}
}

func NewMoveEvent(coordinate Coordinate) Event {
	return Event{
		EventType: MOVED,
		Data:      coordinate,
	}
}

type EventConsumer interface {
	SendEvent(event Event)
}

type StateUpdateConsumer interface {
	StateUpdated(gameState GameState)
}
type StateUpdateSource interface {
	Register(consumer StateUpdateConsumer)
}

func getInitialGameState() GameState {
	initialSide := BLACK

	board := GetInitialBoard()
	used := collectUsed(board)
	edge := collectEdge(board, used)
	possibleMoves := possibleMovesFor(edge, board, initialSide)

	gameState := GameState{
		Board:         board,
		PlayerTurn:    initialSide,
		Used:          used,
		Edge:          edge,
		PossibleMoves: possibleMoves,
	}

	return gameState
}

func possibleMovesFor(edge map[Coordinate]bool, board map[Coordinate]CellClaim, side Player) possibleMoves {
	moves := make(map[Coordinate]bool)
	count := 0
	for e := range edge {
		if IsPossibleMove(side, e, board) {
			moves[e] = true
		}
		count = count + 1
	}

	return possibleMoves{
		side:  side,
		moves: moves,
	}
}

func (player Player) opposite() Player {
	if player == WHITE {
		return BLACK
	}

	return WHITE
}

func applyMove(gameState GameState, side Player, coordinate Coordinate) GameState {
	var owner CellClaim
	if side == BLACK {
		owner = ownedByBlack{}
	} else {
		owner = ownedByWhite{}
	}

	cellsToFlip := getCellsToFlip(gameState.Board, coordinate, side)
	for cell := range cellsToFlip {
		gameState.Board[cell] = owner
	}
	gameState.Board[coordinate] = owner

	gameState.Used[coordinate] = true
	gameState.Edge = updateEdge(gameState.Edge, gameState.Used, coordinate)

	// Try to get moves for the opposite side
	possibleMoves := possibleMovesFor(gameState.Edge, gameState.Board, side.opposite())
	if len(possibleMoves.moves) == 0 {
		// If there are no moves for the opposite side, get moves for the same side
		possibleMoves = possibleMovesFor(gameState.Edge, gameState.Board, side)
	}

	gameState.PossibleMoves = possibleMoves
	gameState.PlayerTurn = possibleMoves.side

	return gameState
}

type gameEventAggregator struct {
	state                GameState
	stateUpdateConsumers []StateUpdateConsumer
}

func (aggregator *gameEventAggregator) SendEvent(event Event) {
	if event.EventType == INITILIZED {
		aggregator.state = getInitialGameState()
	}

	if event.EventType == MOVED {
		side := aggregator.state.PlayerTurn
		coordinate := event.Data.(Coordinate)

		aggregator.state = applyMove(aggregator.state, side, coordinate)
	}

	for _, consumer := range aggregator.stateUpdateConsumers {
		consumer.StateUpdated(aggregator.state)
	}
}

func (aggregator *gameEventAggregator) Register(consumer StateUpdateConsumer) {
	consumers := aggregator.stateUpdateConsumers
	updatedSize := len(consumers) + 1
	updatedConsumers := make([]StateUpdateConsumer, updatedSize)
	for i, consumer := range consumers {
		updatedConsumers[i] = consumer
	}

	updatedConsumers[updatedSize-1] = consumer
	aggregator.stateUpdateConsumers = updatedConsumers
}

type StateUpdaterAndEventConsumer interface {
	StateUpdateSource
	EventConsumer
}

func NewGameEventAggregator() StateUpdaterAndEventConsumer {
	return &gameEventAggregator{stateUpdateConsumers: []StateUpdateConsumer{}}
}
