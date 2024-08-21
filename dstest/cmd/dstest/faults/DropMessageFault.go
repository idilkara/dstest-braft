package faults

import (
	"fmt"
)

type DropMessageFault struct { // SENDER RECEIVER MESSAGE ID NO NEED ANY PARAMETER 

	BaseFault
}

type DropMessageFaultParams struct {
}

var _ Fault = (*DropMessageFault)(nil) //interface 

func NewDropMessageFault(params map[string]interface{}) (*DropMessageFault, error) {
	
	fmt.Println("Creating a new dropmsg Fault: ")
	
	return &DropMessageFault{

		BaseFault: BaseFault{
			Precondition: &AlwaysEnabledPrecondition{},
			Behavior: &DropMessageBehavior{},
		},
	}, nil
}
