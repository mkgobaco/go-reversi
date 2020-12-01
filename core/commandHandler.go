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

func (policy UninitializedGameCommandPolicy) processCommand(command Command, rejectHandler CommandRejectHandler) {
	if command.commandType != INITIALIZE {
		rejectHandler.InvalidCommand(command)
	}

	policy.eventConsumer.SendEvent(NewInitializedEvent())
}

type InProgressCommandPolicy struct {
	eventConsumer EventConsumer
}

func (policy InProgressCommandPolicy) processCommand(command Command, rejectHandler CommandRejectHandler) {
	if command.commandType == INITIALIZE {
		rejectHandler.InvalidCommand(command)
	}
}

func (commandHandler *CommandHandler) AttemptCommand(command Command, rejectHandler CommandRejectHandler) {
	commandHandler.commandPolicy.processCommand(command, rejectHandler)
}
func (commandHandler *CommandHandler) StateUpdated(gameState GameState) {
	commandHandler.commandPolicy = InProgressCommandPolicy{
		eventConsumer: commandHandler.eventConsumer,
	}
}

func NewCommandHandler(eventConsumer EventConsumer, notifier StateUpdateSource) CommandHandler {
	commandHandler := CommandHandler{
		eventConsumer: eventConsumer,
		commandPolicy: UninitializedGameCommandPolicy{eventConsumer: eventConsumer},
	}

	notifier.Register(&commandHandler)

	return commandHandler
}
