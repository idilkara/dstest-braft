package engine

import (
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/process"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/scheduling"
)

type Action struct {
	Sender   int
	Receiver int
	Name     string
}

// FIXME: Repetition of FaultManager interface to avoid cyclic import
// how to avoid this?
type FaultManager interface {
	Init(config *config.Config) error
	GetFaults() []*faults.Fault
	GetEnabledFaults() []*faults.Fault
	PrintFaults()
}

type TestEngine struct {
	Config         *config.Config
	Scheduler      scheduling.Scheduler
	NetworkManager *network.Manager
	ProcessManager *process.ProcessManager
	FaultManager   FaultManager
	Log            *log.Logger

	Experiments   int
	Iterations    int
	Steps         int
	SleepDuration time.Duration
	ReplicaIds    []int
	RecoveryPoint int

	DropNextMessage *bool //for drop message fault 

}

func (te *TestEngine) Init(config *config.Config) error {
	te.Config = config
	te.Experiments = config.TestConfig.Experiments
	te.Iterations = config.TestConfig.Iterations
	te.RecoveryPoint = config.TestConfig.RecoveryPoint
	te.Steps = config.SchedulerConfig.Steps
	te.SleepDuration = time.Duration(config.TestConfig.WaitDuration) * time.Millisecond
	te.ReplicaIds = make([]int, te.Config.ProcessConfig.NumReplicas)
	
	for i := 0; i < te.Config.ProcessConfig.NumReplicas; i++ {
		te.ReplicaIds[i] = i
	}

	scheduler, err := scheduling.NewScheduler(scheduling.SchedulerType(strings.ToLower(config.SchedulerConfig.Type)))
	if err != nil {
		return fmt.Errorf("Error initializing Scheduler: %s", err.Error())
	}

	te.Scheduler = scheduler
	te.NetworkManager = new(network.Manager)
	te.ProcessManager = new(process.ProcessManager)

	// te.ProcessManager.Init(config, te.Iterations)
	te.FaultManager = new(faults.FaultManager)

	if err := te.FaultManager.Init(config); err != nil {
		return fmt.Errorf("Error initializing FaultManager: %s", err.Error())
	}

	te.Log = log.New(os.Stdout, "[TestEngine] ", log.LstdFlags)

	defFalse := false
	te.DropNextMessage = &defFalse //is this safe? - it works 

	return nil
}

func (te *TestEngine) Run() error {

	for i := 0; i < te.Experiments; i++ {
		te.Log.Printf("Starting experiment %d...\n", i+1)

		te.Scheduler.Init(te.Config)
		for j := 0; j < te.Iterations; j++ { // 
			te.Log.Printf("Starting iteration %d\n", j+1)

			// Initialize NetworkManager
			err := te.NetworkManager.Init(te.Config, te.ReplicaIds)
			if err != nil {
				return fmt.Errorf("Error initializing NetworkManager: %s", err.Error())
			}

			// Initialize FaultManager
			err = te.FaultManager.Init(te.Config)
			if err != nil {
				return fmt.Errorf("Error initializing FaultManager: %s", err.Error())
			}
			te.FaultManager.PrintFaults()

			te.ProcessManager.Init(te.Config, te.ReplicaIds, j)
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				te.NetworkManager.Run()
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				te.ProcessManager.Run()
				wg.Done()
			}()

			time.Sleep(time.Duration(te.Config.TestConfig.StartupDuration) * time.Second)

			schedule := make([]Action, 0)
			for s := 0; s < te.Steps; {
				if te.ProcessManager.BugCandidate {
					break
				}
				actions := te.NetworkManager.GetActions()
				sc := te.Scheduler.GetClientRequest()
				if sc >= 0 {
					te.ProcessManager.RunClient(sc)
					schedule = append(schedule, Action{
						Sender:   -1,
						Receiver: -1,
						Name:     fmt.Sprintf("ClientRequest_%d_%d", s, sc),
					})
				}
				// TODO - Get fault from scheduler
				var faultContext faults.FaultContext = NewEngineFaultContext(te)
				decision := te.Scheduler.Next(actions, te.FaultManager.GetFaults(), faultContext)

				if decision.DecisionType == scheduling.SendMessage {
					action := decision.Index
					te.NetworkManager.SendMessage(actions[action].MessageId)
					//get the vector clock info from network manager, add it to name as string

					// Get the vector clock info from the network manager
					senderVC := te.NetworkManager.VectorClocks[actions[action].Sender]
					receiverVC := te.NetworkManager.VectorClocks[actions[action].Receiver]
				
					// Format the vector clocks into a string
					vectorClockInfo := fmt.Sprintf("SenderVC: %v ReceiverVC: %v", senderVC, receiverVC)
					actionNameWithVC := actions[action].Name + " " + vectorClockInfo			

					schedule = append(schedule, Action{
						Sender:   actions[action].Sender,
						Receiver: actions[action].Receiver,
						Name:     actionNameWithVC,
					})
					s++
				}

				if s < te.RecoveryPoint { //can't inject faults after this point
						
					if decision.DecisionType == scheduling.InjectFault {

						fault := te.FaultManager.GetFaults()[decision.Index]
						te.Log.Printf("Applying fault: %+v\n", fault)
						err := (*fault).ApplyBehaviorIfPreconditionMet(&faultContext)
						if err != nil {
							te.Log.Printf("Error applying fault: %s\n", err)
						}
						// TODO - Append fault to schedule

						schedule = append(schedule, Action{
							Sender:   -1,
							Receiver: -1,
							Name:     fmt.Sprintf("Drop Message Fault"),
						})

						//This part is specific for drop message faults

						if te.DropNextMessage != nil && *te.DropNextMessage {//  IF DROP MESSAGE FAULT -> drop the NEXT message
							
							dropMessageIDX := -1
							counterDropMsg := 0

							//don't let it pass if there are no messages in the queues.
							for (len(actions) == 0 || dropMessageIDX == -1 ) && (counterDropMsg != 10){
								counterDropMsg = counterDropMsg +1
								time.Sleep(te.SleepDuration)
								actions = te.NetworkManager.GetActions()

								if len(actions) != 0 {
									dropMessageIDX = te.Scheduler.NextToDrop(actions) 
								}
								//check for infinite loops 

							}
							
						
							if (dropMessageIDX != -1) {
								te.Log.Printf("dropping message ID %d\n",actions[dropMessageIDX].MessageId)
								senderVC := te.NetworkManager.VectorClocks[actions[dropMessageIDX].Sender]
								receiverVC := te.NetworkManager.VectorClocks[actions[dropMessageIDX].Receiver]
								// Format the vector clocks into a string
								vectorClockInfo := fmt.Sprintf("SenderVC: %v ReceiverVC: %v", senderVC, receiverVC)
								droppedActionNameWithVC := actions[dropMessageIDX].Name + " is dropped " + vectorClockInfo	

								schedule = append(schedule, Action{
									Sender:   actions[dropMessageIDX].Sender,
									Receiver: actions[dropMessageIDX].Receiver,
									Name:    droppedActionNameWithVC , 
			
								})

								te.NetworkManager.DropMessage(actions[dropMessageIDX].MessageId)
							}else {
								te.Log.Printf("No message to drop \n")

							}
							dropNext := false
							te.DropNextMessage = &dropNext

						}

					}
				}

				time.Sleep(te.SleepDuration)
			}
			// te.Schedules = append(te.Schedules, schedule)
			te.Log.Println("Shutting down ProcessManager...")
			te.ProcessManager.Shutdown()
			te.Log.Println("Shutting down NetworkManager...")
			te.NetworkManager.Shutdown()
			te.Log.Println("Shutdown complete.")
			wg.Wait()

			te.Log.Println("Checking for bugs...")
			//if te.ProcessManager.BugCandidate {
				outputFile, err := os.OpenFile(filepath.Join(te.ProcessManager.Basedir, "schedule.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					te.Log.Printf("Could not create schedule file.\n Err: %s\n", err)
				}

				for _, action := range schedule {
					fmt.Fprintln(outputFile, action)
				}
				outputFile.Close()
			//}
			te.Log.Println("Iteration complete.")
			te.Scheduler.NextIteration()
		}
		te.Scheduler.Reset()
	}
	return nil
}

type EngineFaultContext struct {
	engine *TestEngine
}

func NewEngineFaultContext(engine *TestEngine) *EngineFaultContext {
	return &EngineFaultContext{engine: engine}
}

func (efc *EngineFaultContext) GetConfig() *config.Config {
	return efc.engine.Config
}

func (efc *EngineFaultContext) GetNetworkManager() *network.Manager {
	return efc.engine.NetworkManager
}

func (efc *EngineFaultContext) GetProcessManager() *process.ProcessManager {
	return efc.engine.ProcessManager
}

//for drop message fault
func (efc *EngineFaultContext) GetDropNextMessage() *bool {
	return efc.engine.DropNextMessage
}

// confirm that EngineFaultContext implements FaultContext
var _ faults.FaultContext = (*EngineFaultContext)(nil)
