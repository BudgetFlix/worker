package job

type State string

const (
	StateReceived   State = "RECEIVED"
	StateValidated  State = "VALIDATED"
	StateProcessing State = "PROCESSING"
	StateCompleted  State = "COMPLETED"
	StateFailed     State = "FAILED"
)