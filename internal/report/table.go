package report

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/detect"
)

func Table(repository string, days int, suspects []detect.Suspect) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Suspected flaky GitHub Actions jobs in %s, last %d days\n\n", repository, days)
	if len(suspects) == 0 {
		fmt.Fprintf(&b, "No likely flaky workflow or job signals matched the selected scan window.\n")
		return b.String()
	}
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "CONF\tWORKFLOW\tJOB\tSHA\tSIGNAL\tCATEGORY")
	for _, suspect := range suspects {
		job := suspect.JobName
		if job == "" {
			job = "(workflow)"
		}
		fmt.Fprintf(tw, "%.2f\t%s\t%s\t%s\t%s\t%s\n",
			suspect.Confidence,
			suspect.WorkflowName,
			job,
			shortSHA(suspect.HeadSHA),
			suspect.Signal,
			suspect.Category,
		)
	}
	tw.Flush()
	fmt.Fprintf(&b, "\nRun with --format markdown to create a report for an issue.\n")
	return b.String()
}
