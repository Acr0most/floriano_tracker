package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"tracker/helpers"
	"tracker/tracker"
)

var startCmd = &cobra.Command{
	Use:   "start [project]",
	Short: "Start monitoring time for the given project.",
	Long: `Start monitoring time for the given project.

  You can add tags indicating
  more specifically what you are working on with '+tag'.

  Example :

  $ tracker start test +foo +bar
  Starting project test [foo bar] at 21:38`,
	Run: start,
}

func init() {
	RootCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Printf("Error: %s\n", helpers.PrintRed("No project given."))
		return
	}

	project := args[0]
	tags := make([]string, 0)

	for i, arg := range args {
		if i == 0 {
			continue
		}

		if arg[0] == '+' {
			tags = append(tags, arg[1:])
		}
	}

	newFrame := tracker.NewFrame(project, tags)
	frames := tracker.GetFrames()

	for _, frame := range frames {
		if frame.End.IsZero() {
			fmt.Printf("Error: %s\n", helpers.PrintRed("Project "+frame.Project+" is already started"))
			return
		}
	}

	frames = append(frames, newFrame)
	frames.Persist()

	fmt.Printf("Starting project %s %s at %s\n", newFrame.FormattedProject(), newFrame.FormattedTags(), newFrame.FormattedStartTime())
}
