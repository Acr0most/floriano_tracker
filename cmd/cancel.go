package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"tracker/helpers"
	"tracker/tracker"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel the last call to the start command.",
	Long:  `Cancel the last call to the start command. The time will not be recorded.`,
	Run:   cancelFrame,
}

func init() {
	RootCmd.AddCommand(cancelCmd)
}

func cancelFrame(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()
	newFrames := make(tracker.Frames, 0, len(frames)-1)

	for _, frame := range frames {
		if !frame.InProgress() {
			newFrames = append(newFrames, frame)
		} else {
			fmt.Printf("Canceling the timer for project %s %s\n", frame.FormattedProject(), frame.FormattedTags())
		}
	}

	if len(newFrames) == len(frames) {
		fmt.Printf("Error: %s\n", helpers.PrintRed("No project started."))
		return
	}

	newFrames.Persist()

}
