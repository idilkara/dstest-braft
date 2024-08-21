package faults

import "fmt"

type DropMessageBehavior struct {
	
}

var _ Behavior = (*DropMessageBehavior)(nil) //interface 

func NewDropMessageBehavior() *DropMessageBehavior {
	return &DropMessageBehavior{
	
	}
}

func (fb *DropMessageBehavior) Apply(context *FaultContext) error {

	//CALL NETWORK MANAGER USING CONTEXT - messageId uint64 of NEXT ACTION HOW to get the message ID 
	// get messages - remove it . 
		
	//Need to write a function that returns next messageID ->tryingt to get the scheduler causes cyclic import
	//decisionMessageID := (*context).GetScheduler().Next_to_drop(nm.GetActions())

	//do the dropping message with test engine 
	dropNext := true
	*(*context).GetDropNextMessage() = dropNext

	return nil
}

func (fb *DropMessageBehavior) String() string {
	return fmt.Sprintf("dropped msg")
}
