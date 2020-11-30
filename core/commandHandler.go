package core

type CommandType string

const (
	INITIALIZE CommandType = "INITIALIZE"
	MOVE       CommandType = "MOVE"
	CONCEDE    CommandType = "CONCEDE"
)

type Command struct {
	commandType CommandType
	data        interface{}
}

func NewInitializeCommand() Command {
	return Command{
		commandType: INITIALIZE,
		data:        nil,
	}
}

type CommandRejectHandler interface {
	InvalidCommand(command Command)
}

type CommandPolicy interface {
	processCommand(command Command, rejectHandler CommandRejectHandler)
}
type CommandHandler struct {
	eventConsumer EventConsumer
	commandPolicy CommandPolicy
}

type UninitializedGameCommandPolicy struct {
	eventConsumer EventConsumer
}

func (uninitializedGameCommandPolicy UninitializedGameCommandPolicy) processCommand(command Command, rejectHandler CommandRejectHandler) {
	if command.commandType != INITIALIZE {
		rejectHandler.InvalidCommand(command)
	}

	uninitializedGameCommandPolicy.eventConsumer.SendEvent(NewInitializedEvent())
}

func (commandHandler *CommandHandler) AttemptCommand(command Command, rejectHandler CommandRejectHandler) {
	commandHandler.commandPolicy.processCommand(command, rejectHandler)
}
func (commandHandler *CommandHandler) StateUpdated(gameState GameState) {
	// commandHandler.commandPolicy = NewInProgressCommandPolicy(
	// 	gameState.activePlayer,
	// 	gameState.availableMoves,
	// )

}

func NewCommandHandler(eventConsumer EventConsumer, notifier StateUpdateSource) CommandHandler {
	commandHandler := CommandHandler{
		eventConsumer: eventConsumer,
		commandPolicy: UninitializedGameCommandPolicy{eventConsumer: eventConsumer},
	}

	notifier.Register(&commandHandler)

	return commandHandler
}
