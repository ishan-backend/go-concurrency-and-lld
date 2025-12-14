package main

import (
	"fmt"
	"time"

	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/scheduler"
)

func main() {
	fmt.Println("=== CPU Scheduler Demo ===\n")

	// Example 1: Priority Scheduling (from Statement 1)
	demonstratePriorityScheduling()

	// Example 2: FCFS Scheduling
	demonstrateFCFSScheduling()

	// Example 3: Round Robin Scheduling
	demonstrateRoundRobinScheduling()

	// Example 4: Runtime Strategy Switching
	demonstrateStrategySwitching()
}

func demonstratePriorityScheduling() {
	fmt.Println("--- Example 1: Priority Scheduling ---")

	// Create a priority scheduler
	priorityScheduler := scheduler.NewPriorityScheduler()
	cpuScheduler := NewCPUScheduler(priorityScheduler)

	fmt.Printf("Using: %s\n\n", cpuScheduler.GetSchedulerName())

	// Add jobs as per Statement 1 example
	cpuScheduler.AddJob("Job1", 5)
	cpuScheduler.AddJob("Job2", 3)
	cpuScheduler.AddJob("Job3", 5)

	fmt.Println("Added jobs:")
	fmt.Println("  - Job1 with priority 5")
	fmt.Println("  - Job2 with priority 3")
	fmt.Println("  - Job3 with priority 5")
	fmt.Println()

	// Execute jobs
	fmt.Println("Executing jobs:")
	for cpuScheduler.HasJobs() {
		jobID, err := cpuScheduler.ExecuteNext()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		fmt.Printf("  ✓ Executed: %s\n", jobID)
	}

	fmt.Println()
}

func demonstrateFCFSScheduling() {
	fmt.Println("--- Example 2: First Come First Served (FCFS) ---")

	fcfsScheduler := scheduler.NewFCFSScheduler()
	cpuScheduler := NewCPUScheduler(fcfsScheduler)

	fmt.Printf("Using: %s\n\n", cpuScheduler.GetSchedulerName())

	// Add jobs
	cpuScheduler.AddJob("JobA", 1)
	time.Sleep(10 * time.Millisecond) // Small delay to ensure different submission times
	cpuScheduler.AddJob("JobB", 5)    // Higher priority, but FCFS ignores it
	time.Sleep(10 * time.Millisecond)
	cpuScheduler.AddJob("JobC", 3)

	fmt.Println("Added jobs (in order): JobA, JobB, JobC")
	fmt.Println("Note: FCFS ignores priority, executes in arrival order\n")

	fmt.Println("Executing jobs:")
	for cpuScheduler.HasJobs() {
		jobID, err := cpuScheduler.ExecuteNext()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		fmt.Printf("  ✓ Executed: %s\n", jobID)
	}

	fmt.Println()
}

func demonstrateRoundRobinScheduling() {
	fmt.Println("--- Example 3: Round Robin Scheduling ---")

	// Create Round Robin scheduler with time quantum of 2
	rrScheduler := scheduler.NewRoundRobinScheduler(2)
	cpuScheduler := NewCPUScheduler(rrScheduler)

	fmt.Printf("Using: %s (Time Quantum: 2)\n\n", cpuScheduler.GetSchedulerName())

	// Add jobs with burst times
	cpuScheduler.AddJobWithBurstTime("JobX", 1, 5) // Needs 5 time units
	cpuScheduler.AddJobWithBurstTime("JobY", 2, 3) // Needs 3 time units
	cpuScheduler.AddJobWithBurstTime("JobZ", 1, 4) // Needs 4 time units

	fmt.Println("Added jobs with burst times:")
	fmt.Println("  - JobX: 5 time units")
	fmt.Println("  - JobY: 3 time units")
	fmt.Println("  - JobZ: 4 time units")
	fmt.Println()

	fmt.Println("Executing jobs (each gets 2 time units per turn):")
	executionCount := 0
	for cpuScheduler.HasJobs() {
		executionCount++
		jobID, err := cpuScheduler.ExecuteNext()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		fmt.Printf("  [Turn %d] Executed: %s\n", executionCount, jobID)
	}

	fmt.Printf("\nTotal execution turns: %d\n", executionCount)
	fmt.Println()
}

func demonstrateStrategySwitching() {
	fmt.Println("--- Example 4: Runtime Strategy Switching ---")

	// Start with Priority Scheduling
	priorityScheduler := scheduler.NewPriorityScheduler()
	cpuScheduler := NewCPUScheduler(priorityScheduler)

	fmt.Printf("Initial Strategy: %s\n", cpuScheduler.GetSchedulerName())

	// Add some jobs
	cpuScheduler.AddJob("Task1", 3)
	cpuScheduler.AddJob("Task2", 1)
	cpuScheduler.AddJob("Task3", 5)

	fmt.Println("Added 3 jobs with priorities 3, 1, 5")

	// Execute one job with Priority
	jobID, _ := cpuScheduler.ExecuteNext()
	fmt.Printf("Executed with Priority: %s (highest priority)\n", jobID)

	// Execute remaining with Priority
	fmt.Println("\nExecuting remaining jobs with Priority:")
	for cpuScheduler.HasJobs() {
		jobID, err := cpuScheduler.ExecuteNext()
		if err != nil {
			break
		}
		fmt.Printf("  ✓ Executed: %s\n", jobID)
	}

	// Switch to FCFS at runtime
	fmt.Println("\nSwitching to FCFS strategy...")
	fcfsScheduler := scheduler.NewFCFSScheduler()
	cpuScheduler.SetStrategy(fcfsScheduler)
	fmt.Printf("New Strategy: %s\n", cpuScheduler.GetSchedulerName())

	// Add new jobs to demonstrate FCFS
	cpuScheduler.AddJob("TaskA", 5) // High priority, but FCFS ignores it
	cpuScheduler.AddJob("TaskB", 1) // Low priority
	cpuScheduler.AddJob("TaskC", 3) // Medium priority

	fmt.Println("Added 3 new jobs with priorities 5, 1, 3")
	fmt.Println("Note: FCFS ignores priority, executes in arrival order")

	// Execute jobs with FCFS
	fmt.Println("\nExecuting jobs with FCFS:")
	for cpuScheduler.HasJobs() {
		jobID, err := cpuScheduler.ExecuteNext()
		if err != nil {
			break
		}
		fmt.Printf("  ✓ Executed: %s\n", jobID)
	}

	fmt.Println("\nThis demonstrates the Strategy Pattern - algorithms can be swapped at runtime!")
	fmt.Println()
}
