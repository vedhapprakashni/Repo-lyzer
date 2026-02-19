package cmd

import (
	"fmt"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/output"
	"github.com/agnivo988/Repo-lyzer/internal/progress"
	"github.com/spf13/cobra"
)

var analyzePRCmd = &cobra.Command{
	Use:   "analyze-pr owner/repo",
	Short: "Analyze Pull Request metrics for a GitHub repository",
	Long: `Analyze pull request metrics including:
  • Average time to merge
  • Review participation (% of PRs with 2+ reviewers)
  • PR size distribution
  • Abandoned PR ratio
  • First-time contributor friendliness

Note: Each PR requires 2 API calls (details + reviews). With authentication,
you have 5,000 requests/hour. Default limit is 100 PRs (200 requests).
Use --limit 0 for no limit, but be cautious of rate limits.

Examples:
  repo-lyzer analyze-pr golang/go
  repo-lyzer analyze-pr microsoft/vscode --state closed --limit 50
  repo-lyzer analyze-pr octocat/Hello-World --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, _ := cmd.Flags().GetString("state")
		limit, _ := cmd.Flags().GetInt("limit")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		// Validate the repository URL format
		owner, repo, err := validateRepoURL(args[0])
		if err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}

		// Record start time for analysis timing
		startTime := time.Now()

		// Initialize GitHub client
		client := github.NewClient()

		// Create progress spinner
		spinner := progress.NewSpinner()

		// Inform user about fetching
		if !jsonOutput {
			spinner.Start(fmt.Sprintf("🔍 Fetching pull requests for %s/%s (state: %s)...", owner, repo, state))
		}

		// Fetch pull requests
		var prs []github.PullRequest
		if limit > 0 {
			prs, err = client.GetPullRequestsWithLimit(owner, repo, state, limit)
		} else {
			prs, err = client.GetPullRequests(owner, repo, state)
		}
		if err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to fetch pull requests: %w", err)
		}

		if len(prs) == 0 {
			spinner.Stop()
			if !jsonOutput {
				fmt.Printf("No pull requests found for %s/%s with state '%s'\n", owner, repo, state)
			}
			return nil
		}

		if !jsonOutput {
			spinner.StopWithMessage(fmt.Sprintf("Found %d pull requests", len(prs)))
		}

		if !jsonOutput {
			spinner.Start("🔄 Fetching PR details and reviews concurrently...")
		}

		// Fetch PR details and reviews concurrently with worker pool
		type prResult struct {
			pr      *github.PullRequest
			reviews []github.Review
			index   int
			err     error
		}

		workers := 10 // Concurrent workers
		semaphore := make(chan struct{}, workers)
		results := make(chan prResult, len(prs))

		// Launch goroutines for each PR
		for i, pr := range prs {
			go func(prNumber, index int) {
				semaphore <- struct{}{}        // Acquire
				defer func() { <-semaphore }() // Release

				// Fetch detailed PR info (includes additions, deletions, changed_files)
				prDetails, err := client.GetPullRequestDetails(owner, repo, prNumber)
				if err != nil {
					results <- prResult{index: index, err: fmt.Errorf("PR #%d details: %w", prNumber, err)}
					return
				}

				// Fetch reviews
				prReviews, err := client.GetPullRequestReviews(owner, repo, prNumber)
				if err != nil {
					results <- prResult{index: index, err: fmt.Errorf("PR #%d reviews: %w", prNumber, err)}
					return
				}

				results <- prResult{
					pr:      prDetails,
					reviews: prReviews,
					index:   index,
				}
			}(pr.Number, i)
		}

		// Collect results
		updatedPRs := make([]*github.PullRequest, len(prs))
		reviews := make(map[int][]github.Review)
		errorCount := 0

		for i := 0; i < len(prs); i++ {
			result := <-results

			if !jsonOutput {
				spinner.Update(fmt.Sprintf("🔄 Fetching PR details and reviews... %d/%d", i+1, len(prs)))
			}

			if result.err != nil {
				errorCount++
				continue
			}

			updatedPRs[result.index] = result.pr
			reviews[result.pr.Number] = result.reviews
		}

		// Filter out nil entries (failed fetches)
		var finalPRs []github.PullRequest
		for _, pr := range updatedPRs {
			if pr != nil {
				finalPRs = append(finalPRs, *pr)
			}
		}

		if !jsonOutput {
			spinner.StopWithMessage(fmt.Sprintf("Fetched %d PRs (%d errors)", len(finalPRs), errorCount))
		}

		if len(finalPRs) == 0 {
			return fmt.Errorf("no PRs could be fetched successfully")
		}

		// Use finalPRs instead of prs
		prs = finalPRs

		// Analyze pull requests
		if !jsonOutput {
			spinner.Start("📊 Analyzing pull request metrics...")
		}
		analytics := analyzer.AnalyzePullRequests(prs, reviews)
		if !jsonOutput {
			spinner.StopWithMessage("Pull request analysis complete")
		}

		// Output results
		if jsonOutput {
			jsonStr, err := output.FormatPRAnalyticsJSON(analytics)
			if err != nil {
				return fmt.Errorf("failed to format JSON: %w", err)
			}
			fmt.Println(jsonStr)
		} else {
			output.PrintPRAnalytics(analytics)

			// Display analysis time
			duration := time.Since(startTime)
			fmt.Printf("⏱️  Analysis completed in %.2f seconds\n", duration.Seconds())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzePRCmd)
	analyzePRCmd.Flags().String("state", "all", "Filter PRs by state: open, closed, or all")
	analyzePRCmd.Flags().Int("limit", 100, "Limit number of PRs to analyze (0 = no limit, use with caution)")
	analyzePRCmd.Flags().Bool("json", false, "Output results as JSON")
}
