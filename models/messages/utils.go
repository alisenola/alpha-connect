package messages

import "gitlab.com/alphaticks/alpha-connect/models"

func (o *NewOrder) IsForceMaker() bool {
	for _, x := range o.ExecutionInstructions {
		if x == models.ExecutionInstruction_ParticipateDoNotInitiate {
			return true
		}
	}
	return false
}
