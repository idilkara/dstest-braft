package scheduling

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"math/rand"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	//"log"
)

type POS struct { 
	Scheduler
	Config                   *config.Config
	RequestQuota             int
	FaultQuota             int
	FaultProbability 		float64
	NumClientTypes           int
	ClientRequestProbability float64
	Rand                     *rand.Rand
	ConflictAnalysis		 bool //true or false get from config
	Priorities				map[uint64]int
	PreviousMessageID		uint64 //atomic clock

}

// TO - DO : delay based on conflicts. - but first make simple pos

var _ Scheduler = &POS{}


func (s *POS) Init(config *config.Config) {
	s.Config = config
	s.RequestQuota = config.SchedulerConfig.ClientRequests
	s.NumClientTypes = len(config.ProcessConfig.ClientScripts)
	s.ClientRequestProbability = config.SchedulerConfig.Params["client_request_probability"].(float64)
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
	s.Priorities = make (map[uint64]int) //
	s.FaultQuota = config.SchedulerConfig.FaultQuota
	s.FaultProbability = config.SchedulerConfig.Params["fault_probability"].(float64)
}


func (s *POS) NextIteration() { //get ready for the next iteration 

	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests // thats required.

	s.FaultQuota = s.Config.SchedulerConfig.FaultQuota



}


func (s *POS) Reset() {

	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests //set quota
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed))) //reset seed 

	s.FaultQuota = s.Config.SchedulerConfig.FaultQuota


}
func (s *POS) Shutdown() {

}


// - make it distinct - 

func (s *POS) distinctRandomInteger() int {
	for {
		i := s.Rand.Int()
		if !s.contains(i, s.Priorities) {
			return i
		}
	}
}


func (s *POS) contains(i int, priorities map[uint64]int) bool {
	for _, val := range priorities {
		if val == i {
			return true
		}
	}
	return false
}


// Returns index of the msg with highest priority
func (s *POS) Next(messages []*network.Message, faults []*faults.Fault, context faults.FaultContext) SchedulerDecision {

	
	//try for a faults
	faultRequestIndex := s.GetFaultRequest(len(faults))
	if faultRequestIndex >= 0 {
		return SchedulerDecision{
			DecisionType: InjectFault,
			Index:        faultRequestIndex,
		}
	}

	// local var go
	var mostPrioritizedMsgID uint64 
	mostPrioritizedMsgID = 0// change that to something that wont occur but it should be changed after anyways if there are msges and otherwise action type will be noop so - 
	mostPriority := -1
	mostPriorIndex := -1

	for i, msg := range messages { // for each message in the messages. - i is the inxed already - 

		if _, exists := s.Priorities[msg.MessageId]; !exists { // check if it exists - if it does not exist - apparently there is a syntax like this.

			s.Priorities[msg.MessageId] = s.distinctRandomInteger() // random priority  - must be different and random so check that. 
		}

		if mostPriority <= s.Priorities[msg.MessageId]{

			mostPrioritizedMsgID = msg.MessageId 
			mostPriority = s.Priorities[mostPrioritizedMsgID]
			mostPriorIndex = i
		}	
	}
		// the max priority is known. 

	if len(messages) > 0 {	
		s.PreviousMessageID = mostPrioritizedMsgID
		return SchedulerDecision{
			DecisionType: SendMessage,
			Index:        mostPriorIndex, 
		}


	} 

	//log.Printf("[*POS SCHEDULER*]NO OPERATION \n")
	return SchedulerDecision{
		DecisionType: NoOp,
	}
	
}


// returns messageIDX to be dropped
func (s *POS) NextToDrop(messages []*network.Message) int {

	// local var go
	var mostPrioritizedMsgID uint64 
	mostPrioritizedMsgID = 0// change that to something that wont occur but it should be changed after anyways if there are msges and otherwise action type will be noop so - 
	mostPriority := -1
	mostPriorIndex := -1

	for i, msg := range messages { // for each message in the messages. - i is the inxed already - 

		if _, exists := s.Priorities[msg.MessageId]; !exists { 
			s.Priorities[msg.MessageId] = s.distinctRandomInteger() // random priority  - must be different and random so check that. 
		}
		if mostPriority <= s.Priorities[msg.MessageId]{
			mostPrioritizedMsgID = msg.MessageId 
			mostPriority = s.Priorities[mostPrioritizedMsgID]
			mostPriorIndex = i
		}	
	}
		// the max priority is known. 

	if len(messages) > 0 {	
		//s.PreviousMessageID = mostPrioritizedMsgID //That's dropped
		return  mostPriorIndex
	} 
	return -1
	

}



func (s *POS) GetClientRequest() int {


	if s.RequestQuota > 0 {
		r := s.Rand.Float64()
		if r <= s.ClientRequestProbability || s.ClientRequestProbability == 1.0 {

			s.RequestQuota--
			return s.Rand.Intn(s.NumClientTypes)

		}
	}

	return -1
}


//returns a fault request based on a probability
func (s *POS) GetFaultRequest(numberOfFaultsAvailable int) int {
	if s.FaultQuota > 0 {

		r := s.Rand.Float64()
		if r <= s.FaultProbability || s.FaultProbability == 1.0 {

			s.FaultQuota--
			return s.Rand.Intn(numberOfFaultsAvailable)

		}
		return -1
	}

	return -1
}

