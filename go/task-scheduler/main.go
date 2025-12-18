package main

import (
	"fmt"
	"sync"
	"time"
	"vcs/models"
)

func main() {
	fmt.Println("=== Task Scheduler Demo ===\n")

	// Create a scheduler with 3 worker threads
	scheduler := NewTaskScheduler(3)
	
	// Track executed tasks
	var executedTasks sync.Map
	
	// Start the scheduler
	fmt.Println("--- Starting scheduler ---")
	err := scheduler.Start()
	if err != nil {
		fmt.Printf("Error starting scheduler: %v\n", err)
		return
	}
	fmt.Println("Scheduler started with 3 worker threads\n")

	// Demo 1: One-time task scheduled for future
	fmt.Println("--- Demo 1: One-time task scheduled for 2 seconds ---")
	task1 := models.NewOneTimeTask(
		"task-1",
		"Print Hello",
		time.Now().Add(2*time.Second),
		func() {
			fmt.Println("  [Task-1] Hello from one-time task!")
			executedTasks.Store("task-1", true)
		},
	)
	scheduler.ScheduleTask(task1)
	fmt.Printf("Scheduled task '%s' to run in 2 seconds\n\n", task1.Name)

	// Demo 2: Delayed task
	fmt.Println("--- Demo 2: Delayed task (3 seconds delay) ---")
	task2 := models.NewDelayedTask(
		"task-2",
		"Delayed Task",
		3*time.Second,
		func() {
			fmt.Println("  [Task-2] This task was delayed by 3 seconds!")
			executedTasks.Store("task-2", true)
		},
	)
	scheduler.ScheduleTask(task2)
	fmt.Printf("Scheduled task '%s' with 3 second delay\n\n", task2.Name)

	// Demo 3: Recurring task
	fmt.Println("--- Demo 3: Recurring task (every 2 seconds) ---")
	recurringCount := 0
	task3 := models.NewRecurringTask(
		"task-3",
		"Recurring Task",
		time.Now().Add(1*time.Second),
		2*time.Second,
		func() {
			recurringCount++
			fmt.Printf("  [Task-3] Recurring task execution #%d\n", recurringCount)
			executedTasks.Store(fmt.Sprintf("task-3-%d", recurringCount), true)
			
			// Stop after 5 executions
			if recurringCount >= 5 {
				scheduler.CancelTask("task-3")
				fmt.Println("  [Task-3] Recurring task cancelled after 5 executions")
			}
		},
	)
	scheduler.ScheduleTask(task3)
	fmt.Printf("Scheduled recurring task '%s' (every 2 seconds, starting in 1 second)\n\n", task3.Name)

	// Demo 4: Multiple concurrent tasks
	fmt.Println("--- Demo 4: Multiple concurrent tasks ---")
	for i := 1; i <= 5; i++ {
		taskID := fmt.Sprintf("concurrent-%d", i)
		task := models.NewDelayedTask(
			taskID,
			fmt.Sprintf("Concurrent Task %d", i),
			time.Duration(i)*500*time.Millisecond,
			func() {
				fmt.Printf("  [%s] Executing concurrently\n", taskID)
				time.Sleep(1 * time.Second) // Simulate work
				executedTasks.Store(taskID, true)
			},
		)
		scheduler.ScheduleTask(task)
	}
	fmt.Println("Scheduled 5 concurrent tasks\n")

	// Demo 5: Task status check
	fmt.Println("--- Demo 5: Checking task status ---")
	time.Sleep(500 * time.Millisecond)
	status, err := scheduler.GetTaskStatus("task-1")
	if err == nil {
		fmt.Printf("Task-1 Status: RunCount=%d, ExecuteAt=%s\n", status.RunCount, status.ExecuteAt.Format("15:04:05"))
	}
	fmt.Println()

	// Wait for tasks to execute
	fmt.Println("--- Waiting for tasks to execute ---")
	time.Sleep(15 * time.Second)

	// Demo 6: Cancel a task
	fmt.Println("\n--- Demo 6: Cancelling a task ---")
	cancelTask := models.NewDelayedTask(
		"cancel-me",
		"Will be cancelled",
		5*time.Second,
		func() {
			fmt.Println("  [cancel-me] This should not execute!")
		},
	)
	scheduler.ScheduleTask(cancelTask)
	fmt.Println("Scheduled task 'cancel-me'")
	time.Sleep(1 * time.Second)
	err = scheduler.CancelTask("cancel-me")
	if err == nil {
		fmt.Println("Cancelled task 'cancel-me' successfully")
	}
	time.Sleep(5 * time.Second)
	fmt.Println()

	// Demo 7: Long-running task (shouldn't block other tasks)
	fmt.Println("--- Demo 7: Long-running task (non-blocking) ---")
	longTask := models.NewDelayedTask(
		"long-task",
		"Long Running Task",
		1*time.Second,
		func() {
			fmt.Println("  [long-task] Starting long-running task (5 seconds)...")
			time.Sleep(5 * time.Second)
			fmt.Println("  [long-task] Long-running task completed!")
			executedTasks.Store("long-task", true)
		},
	)
	scheduler.ScheduleTask(longTask)
	
	// Schedule another task that should run during the long task
	quickTask := models.NewDelayedTask(
		"quick-task",
		"Quick Task",
		2*time.Second,
		func() {
			fmt.Println("  [quick-task] This task ran while long-task was executing!")
			executedTasks.Store("quick-task", true)
		},
	)
	scheduler.ScheduleTask(quickTask)
	fmt.Println("Scheduled long-running task and quick task\n")

	// Wait a bit more
	time.Sleep(8 * time.Second)

	// Show pending tasks count
	fmt.Println("\n--- Scheduler Statistics ---")
	fmt.Printf("Pending tasks: %d\n", scheduler.GetPendingTasksCount())
	fmt.Printf("Scheduler running: %v\n", scheduler.IsRunning())

	// Stop the scheduler
	fmt.Println("\n--- Stopping scheduler ---")
	scheduler.Stop()
	fmt.Println("Scheduler stopped gracefully")

	fmt.Println("\n=== Demo Complete ===")
}

// Example: Scheduled backup task
func demonstrateBackupTask(scheduler *TaskScheduler) {
	backupTask := models.NewRecurringTask(
		"backup-task",
		"Daily Backup",
		time.Now().Add(24*time.Hour), // Start tomorrow
		24*time.Hour,                  // Run daily
		func() {
			fmt.Println("  [Backup] Performing daily backup...")
			// Backup logic here
			fmt.Println("  [Backup] Backup completed!")
		},
	)
	scheduler.ScheduleTask(backupTask)
}

// Example: Scheduled cleanup task
func demonstrateCleanupTask(scheduler *TaskScheduler) {
	cleanupTask := models.NewRecurringTaskWithDelay(
		"cleanup-task",
		"Cleanup Old Files",
		1*time.Hour,    // Start in 1 hour
		6*time.Hour,    // Run every 6 hours
		func() {
			fmt.Println("  [Cleanup] Cleaning up old files...")
			// Cleanup logic here
			fmt.Println("  [Cleanup] Cleanup completed!")
		},
	)
	scheduler.ScheduleTask(cleanupTask)
}

