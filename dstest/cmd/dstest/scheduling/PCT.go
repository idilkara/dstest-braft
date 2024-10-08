package scheduling

import (
	"math/rand"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type PCT struct {
	Scheduler
	Config                   *config.Config
	RequestQuota             int
	FaultQuota             int
	NumClientTypes           int
	ClientRequestProbability float64
	FaultProbability float64
	NetworkManager           *network.Manager
	Rand                     *rand.Rand
	Depth                    int
	InitialPriorities        []int
	PriorityChangePoints     []int
	NumPriorityChange        int
	Step                     int
}

// assert PCT implements the Scheduler interface
var _ Scheduler = &PCT{}

func (s *PCT) Init(config *config.Config) {
	s.Config = config
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
	s.RequestQuota = config.SchedulerConfig.ClientRequests
	s.NumClientTypes = len(config.ProcessConfig.ClientScripts)
	s.ClientRequestProbability = config.SchedulerConfig.Params["client_request_probability"].(float64)
	s.Depth = config.SchedulerConfig.Params["d"].(int)
	s.NumPriorityChange = 0
	s.Step = 0
	s.PriorityChangePoints = make([]int, 0)
	s.InitialPriorities = make([]int, 0)
	for i := 0; i < s.Depth; i++ {
		s.PriorityChangePoints = append(s.PriorityChangePoints, s.distinctRandomInteger(s.Config.SchedulerConfig.Steps, s.PriorityChangePoints))
	}
	for i := 0; i < s.Config.ProcessConfig.NumReplicas; i++ {
		s.InitialPriorities = append(s.InitialPriorities, s.distinctRandomInteger(s.Config.ProcessConfig.NumReplicas, s.InitialPriorities))
	}
	s.FaultQuota = config.SchedulerConfig.FaultQuota
	s.FaultProbability = config.SchedulerConfig.Params["fault_probability"].(float64)
}

func (s *PCT) distinctRandomInteger(max int, arr []int) int {
	for {
		i := s.Rand.Intn(max)
		if !contains(i, arr) {
			return i
		}
	}
}

func contains(i int, s []int) bool {
	if len(s) == 0 {
		return false
	}

	for _, val := range s {
		if val == i {
			return true
		}
	}
	return false
}

func (s *PCT) NextIteration() {
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
	s.NumPriorityChange = 0
	s.Step = 0
	s.PriorityChangePoints = make([]int, 0)
	s.InitialPriorities = make([]int, 0)
	for i := 0; i < s.Depth; i++ {
		s.PriorityChangePoints = append(s.PriorityChangePoints, s.distinctRandomInteger(s.Config.SchedulerConfig.Steps, s.PriorityChangePoints))
	}
	for i := 0; i < s.Config.ProcessConfig.NumReplicas; i++ {
		s.InitialPriorities = append(s.InitialPriorities, s.distinctRandomInteger(s.Config.ProcessConfig.NumReplicas, s.InitialPriorities))
	}
	s.FaultQuota = s.Config.SchedulerConfig.FaultQuota
}

func (s *PCT) Reset() {
	s.FaultQuota = s.Config.SchedulerConfig.FaultQuota
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
}
func (s *PCT) Shutdown() {

}

// Returns a random index from available messages
func (s *PCT) Next(messages []*network.Message, faults []*faults.Fault, context faults.FaultContext) SchedulerDecision {
	
	//give fault a chance
	numberOfFaultsAvailable := len(faults)

	faultRequestIndex := s.GetFaultRequest(numberOfFaultsAvailable)
	// returns fault
	if faultRequestIndex >= 0 {
		return SchedulerDecision{
			DecisionType: InjectFault,
			Index:        faultRequestIndex,
		}
	}
	
	
	if s.NumPriorityChange < s.Depth {
		if s.Step == s.PriorityChangePoints[s.NumPriorityChange] {
			s.NumPriorityChange++
		}
	}

	if len(messages) > 0 {
		s.Step++
		decision := -1
		for i, message := range messages {
			if message.Sender == s.InitialPriorities[s.NumPriorityChange] {
				decision = i
			}
		}

		if decision < 0 {
			return SchedulerDecision{
				DecisionType: NoOp,
			}
		}

		return SchedulerDecision{
			DecisionType: SendMessage,
			Index:        decision,
		}
	} else {
		return SchedulerDecision{
			DecisionType: NoOp,
		}
	}
}

func (s *PCT) GetClientRequest() int {
	if s.RequestQuota > 0 {
		r := s.Rand.Float64()
		if r <= s.ClientRequestProbability || s.ClientRequestProbability == 1.0 {
			s.RequestQuota--
			return s.Rand.Intn(s.NumClientTypes)
		}
	}
	return -1
}


//FOR DROP MESSAGE FAULT
//returns a fault request based on a probability - INDEPENDENT OF THE SCHEDULERs TYPE
func (s *PCT) GetFaultRequest(numberOfFaultsAvailable int) int {

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

// returns messageIDX to be dropped
func (s *PCT) NextToDrop(messages []*network.Message) int {

	//priority change
	if s.NumPriorityChange < s.Depth {
		if s.Step == s.PriorityChangePoints[s.NumPriorityChange] {
			s.NumPriorityChange++
		}
	}

	
	if len(messages) > 0 {

		s.Step++
		decision := -1 //decision being -1 

		for i, message := range messages {
			if message.Sender == s.InitialPriorities[s.NumPriorityChange] {
				decision = i
			}
		}

		return decision //invalid decision is already <0 

	} else { 
		return -1
	}

}
