package scheduling

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"math/rand"
	
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type RandomScheduler struct {
	Scheduler
	Config                   *config.Config
	RequestQuota             int
	NumClientTypes           int
	ClientRequestProbability float64
	Rand                     *rand.Rand
	FaultProbability		float64
	FaultQuota				int
}

// assert RandomScheduler implements the Scheduler interface
var _ Scheduler = &RandomScheduler{}

func (s *RandomScheduler) Init(config *config.Config) {
	s.Config = config
	s.RequestQuota = config.SchedulerConfig.ClientRequests
	s.NumClientTypes = len(config.ProcessConfig.ClientScripts)
	s.ClientRequestProbability = config.SchedulerConfig.Params["client_request_probability"].(float64)
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
	s.FaultProbability = config.SchedulerConfig.Params["fault_probability"].(float64) // add fault quota
	s.FaultQuota = config.SchedulerConfig.FaultQuota

}

func (s *RandomScheduler) NextIteration() {
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
	s.FaultQuota = s.Config.SchedulerConfig.FaultQuota
}

func (s *RandomScheduler) Reset() {
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
	s.FaultQuota = s.Config.SchedulerConfig.FaultQuota
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
}
func (s *RandomScheduler) Shutdown() {

}

// Returns a random index from available messages
func (s *RandomScheduler) Next(messages []*network.Message, faults []*faults.Fault, context faults.FaultContext) SchedulerDecision {
	
	//try for a faults
	faultRequestIndex := s.GetFaultRequest(len(faults))
	if faultRequestIndex >= 0 {
		return SchedulerDecision{
			DecisionType: InjectFault,
			Index:        faultRequestIndex,
		}
	}

	// return message 
	if len(messages) > 0 {
		return SchedulerDecision{
			DecisionType: SendMessage,
			Index:        s.Rand.Intn(len(messages)),
		}
	} else {
		return SchedulerDecision{
			DecisionType: NoOp,
		}
	}
}

func (s *RandomScheduler) GetClientRequest() int {
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
func (s *RandomScheduler) GetFaultRequest(numberOfFaultsAvailable int) int {
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
func (s *RandomScheduler) NextToDrop(messages []*network.Message) int {

	if len(messages) > 0 {
		messageIDX := s.Rand.Intn(len(messages))
		//log.Printf("msg TO DROP:",messageIDX)
		return messageIDX

	} else { // If nothing to drop return -1 ?? 
		return -1
	}

}