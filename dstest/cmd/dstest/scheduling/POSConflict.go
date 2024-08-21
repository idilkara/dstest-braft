package scheduling

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"math/rand"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	//"github.com/egeberkaygulcan/dstest/cmd/dstest/engine" //import cycle 
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	//"log"
)



type POSConflict struct { 
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
	PreviousMessage			*network.Message 
	
	//NetworkManager			*network.Manager
}


var _ Scheduler = &POSConflict{}

func (s *POSConflict) Init(config *config.Config) {
	s.Config = config
	s.RequestQuota = config.SchedulerConfig.ClientRequests
	s.NumClientTypes = len(config.ProcessConfig.ClientScripts)
	s.ClientRequestProbability = config.SchedulerConfig.Params["client_request_probability"].(float64)
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
	//s.NetworkManager = config.SchedulerConfig.Params["networkMan"].(*network.Manager)  
	s.Priorities = make (map[uint64]int) 
	s.FaultQuota = config.SchedulerConfig.FaultQuota
	s.FaultProbability = config.SchedulerConfig.Params["fault_probability"].(float64)
}


func (s *POSConflict) NextIteration() { //get ready for the next iteration 

	s.FaultQuota = s.Config.SchedulerConfig.FaultQuota
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests // thats required.

}


func (s *POSConflict) Reset() {

	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests //set quota
	
	s.FaultQuota = s.Config.SchedulerConfig.FaultQuota

	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed))) //reset seed 
	s.Priorities = make (map[uint64]int) //resets? 
	
}

func (s *POSConflict) Shutdown() {

}


// - make it distinct - 

func (s *POSConflict) distinctRandomInteger() int {
	for {
		i := s.Rand.Int()
		if !s.contains(i, s.Priorities) {
			return i
		}
	}
}


func (s *POSConflict) contains(i int, priorities map[uint64]int) bool {
	for _, val := range priorities {
		if val == i {
			return true
		}
	}
	return false
}


func isIndexDisabled(index int, disabledIndexes []int) bool {
	for _, disabledIndex := range disabledIndexes {
		if disabledIndex == index {
			return true
		}
	}
	return false
}


// Returns a random index from available messages - not
func (s *POSConflict) Next(messages []*network.Message, faults []*faults.Fault, context faults.FaultContext) SchedulerDecision {
	//try for a faults
	faultRequestIndex := s.GetFaultRequest(len(faults))
	if faultRequestIndex >= 0 {
		return SchedulerDecision{
			DecisionType: InjectFault,
			Index:        faultRequestIndex,
		}
	}

	//set priority
	var mostPrioritizedMsgID uint64 
	mostPrioritizedMsgID = 0 // change that to something that wont occur but it should be changed after anyways if there are msges and otherwise action type will be noop so - 
	mostPriority := -1
	mostPriorIndex := -1

	// these are -1 so the random should assign unsigned int

	for i, msg := range messages { // for each message in the messages. - i is the inxed already - 
		
		if _, exists := s.Priorities[msg.MessageId]; !exists { // check if it exists - if it does not exist - apparently there is a syntax like this.

			s.Priorities[msg.MessageId] = s.Rand.Int() // random priority  - must be different and random so check that. 
			//log.Printf("[*pos SCHEDULER*]THERE IS A NEW MSG and priority is set to : %d \n" ,s.Priorities[msg.MessageId] )
		}	
		// get the max while already traversing over the msges list

		if mostPriority <= s.Priorities[msg.MessageId]{

			mostPrioritizedMsgID = msg.MessageId 
			mostPriority = s.Priorities[mostPrioritizedMsgID]
			mostPriorIndex = i
		}

	}


	// the max priority is known. 

	if len(messages) > 0 {

		MaybeConflictingMostPriorIndex :=  mostPriorIndex
		/// WHILE CONFLICTS - SELECT NEW - IF EVERYTING IS CONFLICTING - CONTNUE WITH THE CONFLICTING SELECTION 
		var disabledIndexes []int

		for { 
			if (!s.IsConflicting(messages[mostPriorIndex])) || len(disabledIndexes) == len(messages) {
				break
			}
			// conflicting -> select new one
			disabledIndexes = append(disabledIndexes, mostPriorIndex)
			for i, msg := range messages { 

				if isIndexDisabled(i, disabledIndexes) {
					continue
				}
				if mostPriority <= s.Priorities[msg.MessageId]{
					mostPrioritizedMsgID = msg.MessageId 
					mostPriority = s.Priorities[mostPrioritizedMsgID]
					mostPriorIndex = i
				}
			}
		}
		// after this -  if all conflict

		if (len(disabledIndexes) == len(messages)){
			mostPriorIndex = MaybeConflictingMostPriorIndex
		}

		//log.Printf("[*pos SCHEDULER*]action SELECTED  %d\n", mostPrioritizedMsgID )
		//RETURN START
		s.PreviousMessage = messages[mostPriorIndex] 
		return SchedulerDecision{
			DecisionType: SendMessage,
			Index:        mostPriorIndex, 
		}
		//RETURN END
	} 

	return SchedulerDecision{
		DecisionType: NoOp,
	}
	
}


// returns messageIDX to be dropped
func (s *POSConflict) NextToDrop(messages []*network.Message) int {

	// set priority
	var mostPrioritizedMsgID uint64 
	mostPrioritizedMsgID = 0 // change that to something that wont occur but it should be changed after anyways if there are msges and otherwise action type will be noop so - 
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
		MaybeConflictingMostPriorIndex :=  mostPriorIndex
		// WHILE CONFLICTS - SELECT NEW - IF EVERYTING IS CONFLICTING - CONTNUE WITH THE CONFLICTING SELECTION 
		var disabledIndexes []int

		for { 
			if (!s.IsConflicting(messages[mostPriorIndex])) || len(disabledIndexes) == len(messages) {
				break
			}
			// conflicting -> select new one
			disabledIndexes = append(disabledIndexes, mostPriorIndex)
			for i, msg := range messages { 

				if isIndexDisabled(i, disabledIndexes) {
					continue
				}
				if mostPriority <= s.Priorities[msg.MessageId]{
					mostPrioritizedMsgID = msg.MessageId 
					mostPriority = s.Priorities[mostPrioritizedMsgID]
					mostPriorIndex = i
				}
			}
		}
		// after this -  if all conflict
		if (len(disabledIndexes) == len(messages)){
			mostPriorIndex = MaybeConflictingMostPriorIndex
		}

		//RETURN START
		//s.PreviousMessage = messages[mostPriorIndex] 
		return  mostPriorIndex
		//RETURN END
	} 
	return -1	
}

func (s *POSConflict) GetClientRequest() int {
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
func (s *POSConflict) GetFaultRequest(numberOfFaultsAvailable int) int {
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

// NOT CORRECT * - ✨ VECTOR CLOCK CHECK IS NOT NECESSARY ✨
func (s *POSConflict) IsConflicting(message *network.Message) bool {
	//log.Printf("[*POSConflict SCHEDULER*] CHECKING IF THEY CONFLICT  \n")
	//if PreviousMessage - make

	// THE VECTOR CLOCK IS NOT ATTACHED TO THE MESSAGE - so when msg is the action executed then 
	// the vector clock is updated and then this changes the relation

	if s.PreviousMessage != nil {
		// check receiver

			if message.Receiver ==  s.PreviousMessage.Receiver {
				return true
			}
		
    }
	
	return false
}