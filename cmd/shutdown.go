package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"syscall"
	"time"
)

var pid int

var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shut down PufferPanel",
	Run:   executeShutdown,
}

func init() {
	shutdownCmd.Flags().IntVar(&pid, "pid", 0, "PID to send SIGINT to")
	_ = shutdownCmd.MarkFlagRequired("pid")
}

func executeShutdown(cmd *cobra.Command, args []string) {
	if pid == 0 {
		fmt.Printf("PID required")
		os.Exit(1)
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil || proc == nil {
		if err == nil {
			err = errors.New("no process found")
		}
		fmt.Printf("Error finding process: %s\n", err.Error())
		os.Exit(1)
		return
	}

	err = proc.Signal(syscall.Signal(15))
	if err != nil {
		fmt.Printf("Error sending signal: %s\n", err.Error())
		os.Exit(1)
		return
	}

	wait := make(chan error)

	waitForProcess(proc, wait)

	err = <-wait

	if err != nil {
		fmt.Printf("Error shutting down: %s\n", err.Error())
		os.Exit(1)
		return
	}

	err = proc.Release()
	if err != nil {
		fmt.Printf("Error shutting down: %s\n", err.Error())
		os.Exit(1)
		return
	}
}

func waitForProcess(process *os.Process, c chan error) {
	var err error
	timer := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range timer.C {
			err = process.Signal(syscall.Signal(0))
			if err != nil {
				if err.Error() == "os: process already finished" {
					c <- nil
				} else {
					c <- err
				}

				timer.Stop()
			}
		}
	}()
}
