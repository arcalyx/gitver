package cmd

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/constants"
	"github.com/arcalyx/gitver/internal/gitops"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var (
	// Version bump flags
	autoFlag  bool
	majorFlag bool
	minorFlag bool
	patchFlag bool

	// Pre-release flags
	prereleaseFlag        string
	prereleaseBumpFlag    bool
	prereleasePromoteFlag bool
	prereleaseClearFlag   bool
	prereleaseVersionFlag int

	// Build metadata flags
	buildMetaFlag      string
	buildMetaClearFlag bool

	// Git operation flags
	commitFlag bool
	tagFlag    bool
	pushFlag   bool
	amend      bool

	// Other flags
	updatePackages bool
	dryRunFlag     bool
	changelogFlag  bool
	syncFlag       bool
	verbose        bool
)

const (
	PriorityNone = iota
	PriorityFix
	PriorityFeat
	PriorityBreakingChange
)

const (
	message0001 = "Please provide a valid version bump flag: --auto, --major, --minor, or --patch"
	message0002 = "Version bumped: %v -> %v"
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump the version of the project",
	Long: `Bump the version of the project and perform optional actions.
	
This command allows you to increment the version number according to semantic versioning rules,
set pre-release versions, add build metadata, and perform git operations in one step.

Examples:
  gitver bump --major                # Bump major version (x.0.0)
  gitver bump --minor                # Bump minor version (0.x.0)
  gitver bump --patch                # Bump patch version (0.0.x)
  gitver bump --auto                 # Automatically determine version bump based on commit messages
  gitver bump --minor --prerelease beta  # Set version to 0.x.0-beta.1
  gitver bump --build-meta build.123     # Add build metadata
  gitver bump --minor --tag --push       # Bump, tag and push changes
  gitver bump --minor --changelog --sync # Bump, generate changelog and sync package versions
  gitver bump --minor --dry-run          # Show what would happen without making changes`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set the updatePackages flag in the version package
		version.SetUpdatePackages(updatePackages)

		// Validate flags
		if dryRunFlag {
			log.Println("Running in dry-run mode - no changes will be made")
		}

		// Load configuration
		if !dryRunFlag {
			loadConfig()
		}

		// Process version bump if requested
		if majorFlag || minorFlag || patchFlag || autoFlag {
			if !dryRunFlag {
				// Prepare git operations if needed
				if tagFlag || pushFlag {
					err := prepareGitOperation()
					if err != nil {
						log.Fatal(err)
					}
				}

				// Execute the appropriate bump
				switch {
				case majorFlag:
					executeMajorMode()
				case minorFlag:
					executeMinorMode()
				case patchFlag:
					executePatchMode()
				case autoFlag:
					executeAutoMode()
				}

				log.Printf(message0002, version.GetLastVersion(), version.ToString())
			} else {
				// In dry-run mode, just show what would happen
				bumpType := "unknown"
				if majorFlag {
					bumpType = "major"
				} else if minorFlag {
					bumpType = "minor"
				} else if patchFlag {
					bumpType = "patch"
				} else if autoFlag {
					bumpType = "auto (determined from commits)"
				}
				log.Printf("Would bump version (%s)", bumpType)
			}
		}

		// Process pre-release if requested
		if prereleaseFlag != "" || prereleaseBumpFlag || prereleasePromoteFlag || prereleaseClearFlag {
			if !dryRunFlag {
				var err error

				switch {
				case prereleaseFlag != "":
					if cmd.Flags().Changed("prerelease-version") {
						err = version.SetPreRelease(prereleaseFlag, prereleaseVersionFlag)
					} else {
						err = version.SetPreRelease(prereleaseFlag, 1)
					}
				case prereleaseBumpFlag:
					err = version.BumpPreRelease()
				case prereleasePromoteFlag:
					err = version.PromotePreRelease()
				case prereleaseClearFlag:
					err = version.ClearPreRelease()
				}

				if err != nil {
					log.Fatalf("Error processing pre-release: %v", err)
				}

				log.Printf("Pre-release updated: %s", version.ToString())
			} else {
				// In dry-run mode
				action := "unknown"
				if prereleaseFlag != "" {
					action = "set to " + prereleaseFlag
				} else if prereleaseBumpFlag {
					action = "bump"
				} else if prereleasePromoteFlag {
					action = "promote"
				} else if prereleaseClearFlag {
					action = "clear"
				}
				log.Printf("Would update pre-release (%s)", action)
			}
		}

		// Process build metadata if requested
		if buildMetaFlag != "" || buildMetaClearFlag {
			if !dryRunFlag {
				var err error

				if buildMetaFlag != "" {
					err = version.SetBuildMetadata(buildMetaFlag)
					if err != nil {
						log.Fatalf("Failed to set build metadata: %v", err)
					}
					log.Printf("Build metadata set to '%s'", buildMetaFlag)
				} else if buildMetaClearFlag {
					err = version.ClearBuildMetadata()
					if err != nil {
						log.Fatalf("Failed to clear build metadata: %v", err)
					}
					log.Printf("Build metadata cleared")
				}

				log.Printf("Version with build metadata: %s", version.ToString())
			} else {
				// In dry-run mode
				action := "unknown"
				if buildMetaFlag != "" {
					action = "set to " + buildMetaFlag
				} else if buildMetaClearFlag {
					action = "clear"
				}
				log.Printf("Would update build metadata (%s)", action)
			}
		}

		// Generate changelog if requested
		if changelogFlag && !dryRunFlag {
			log.Println("Generating changelog...")
			// TODO: Implement changelog generation
		} else if changelogFlag {
			log.Println("Would generate changelog")
		}

		// Sync package versions if requested
		if syncFlag && !dryRunFlag {
			log.Println("Syncing package versions...")
			// TODO: Implement sync functionality
		} else if syncFlag {
			log.Println("Would sync package versions")
		}

		// Execute git operations if requested
		if (commitFlag || tagFlag || pushFlag) && !dryRunFlag {
			executeGitOperations()
		} else if commitFlag || tagFlag || pushFlag {
			actions := []string{}
			if commitFlag {
				actions = append(actions, "commit")
			}
			if tagFlag {
				actions = append(actions, "tag")
			}
			if pushFlag {
				actions = append(actions, "push")
			}
			log.Printf("Would perform git operations: %s", strings.Join(actions, ", "))
		}
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)

	// Version bump flags
	bumpCmd.Flags().BoolVarP(&autoFlag, "auto", "a", false, "Automatically determine version bump based on commit messages")
	bumpCmd.Flags().BoolVarP(&majorFlag, "major", "M", false, "Bump major version")
	bumpCmd.Flags().BoolVarP(&minorFlag, "minor", "m", false, "Bump minor version")
	bumpCmd.Flags().BoolVarP(&patchFlag, "patch", "p", false, "Bump patch version")

	// Pre-release flags
	bumpCmd.Flags().StringVar(&prereleaseFlag, "prerelease", "", "Set pre-release version (alpha, beta, rc)")
	bumpCmd.Flags().BoolVar(&prereleaseBumpFlag, "prerelease-bump", false, "Bump the current pre-release version")
	bumpCmd.Flags().BoolVar(&prereleasePromoteFlag, "prerelease-promote", false, "Promote to next pre-release level")
	bumpCmd.Flags().BoolVar(&prereleaseClearFlag, "prerelease-clear", false, "Clear pre-release version")
	bumpCmd.Flags().IntVar(&prereleaseVersionFlag, "prerelease-version", 1, "Specify pre-release version number")

	// Build metadata flags
	bumpCmd.Flags().StringVar(&buildMetaFlag, "build-meta", "", "Set build metadata")
	bumpCmd.Flags().BoolVar(&buildMetaClearFlag, "build-meta-clear", false, "Clear build metadata")

	// Git operation flags
	bumpCmd.Flags().BoolVarP(&commitFlag, "commit", "c", false, "Create a commit with the version change")
	bumpCmd.Flags().BoolVarP(&tagFlag, "tag", "t", false, "Create a tag for the new version")
	bumpCmd.Flags().BoolVarP(&pushFlag, "push", "", false, "Push changes to remote repository")
	bumpCmd.Flags().BoolVarP(&amend, "amend", "", false, "Amend last commit")

	// Other flags
	bumpCmd.Flags().BoolVarP(&updatePackages, "update-packages", "u", true, "Update version in package manager files")
	bumpCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Show what would happen without making changes")
	bumpCmd.Flags().BoolVar(&changelogFlag, "changelog", false, "Generate changelog")
	bumpCmd.Flags().BoolVar(&syncFlag, "sync", false, "Sync versions in package files")

	// Set the flag in the version package
	version.SetUpdatePackages(updatePackages)
}

func executeMajorMode() {
	log.Println("bump major version")
	err := version.BumpMajor()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("bump major version success")
}

func executeMinorMode() {
	log.Println("bump minor version")
	err := version.BumpMinor()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("bump minor version success")
}

func executePatchMode() {
	log.Println("bump patch version success")
	err := version.BumpPatch()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("bump patch version success")
}

func executeAutoMode() {
	log.Println("start auto mode")
	bumpFunc, err := detectAutoBump()
	if err != nil {
		log.Fatal(err)
	}
	err = bumpFunc()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("start auto mode success")
}

func detectAutoBump() (func() error, error) {
	log.Println("start detect bump function for auto mode")
	tag, err := gitops.GetLastTag()
	if err != nil {
		return nil, err
	}

	if tag == fmt.Sprintf(constants.ReleaseTag, version.ToString()) {
		return analyzeCommits(tag)
	} else if tag == fmt.Sprintf(constants.VersionTag, version.ToString()) {
		return analyzeAndCompareCommits(tag, fmt.Sprintf(constants.ReleaseTag, version.GetLastVersion()))
	} else if tag == "" {
		return analyzeCommits(tag)
	} else {
		return nil, fmt.Errorf("no new version required")
	}
}

func analyzeCommits(tag string) (func() error, error) {
	log.Println("analyze commits from head to", tag)
	commits, err := gitops.GetCommits(tag)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(len(commits), "commits found")

	priority := findCommitPriority(commits)
	log.Println("bump priority are:", priority)
	switch priority {
	case PriorityBreakingChange:
		return version.BumpMajor, nil
	case PriorityFeat:
		return version.BumpMinor, nil
	case PriorityFix:
		return version.BumpPatch, nil
	default:
		return nil, fmt.Errorf("no new version required")
	}
}

func analyzeAndCompareCommits(starttag, endtag string) (func() error, error) {
	log.Println("analyze commits from", starttag, "to", endtag)
	oldCommits, err := gitops.GetCommitsBetweenTags(starttag, endtag)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(oldCommits), "commits found")

	log.Println("analyze commits from head to", starttag)
	newCommits, err := gitops.GetCommits(starttag)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(newCommits), "commits found")

	priority := comparePriority(findCommitPriority(oldCommits), findCommitPriority(newCommits))
	log.Println("bump priority are:", priority)
	switch priority {
	case PriorityBreakingChange:
		return version.BumpMajor, nil
	case PriorityFeat:
		return version.BumpMinor, nil
	case PriorityFix:
		return version.BumpPatch, nil
	default:
		return nil, fmt.Errorf("cannot detect bump function")
	}
}

func comparePriority(first, second int) int {
	if second > first {
		return second
	}
	return 0
}

func findCommitPriority(commits []*object.Commit) int {
	highestPriority := PriorityNone

	for _, commit := range commits {
		message := commit.Message
		if strings.Contains(message, "BREAKING CHANGE") || strings.Contains(message, "#breaking") {
			return PriorityBreakingChange
		} else if strings.Contains(message, "feat") || strings.Contains(message, "feature") {
			highestPriority = comparePriority(highestPriority, PriorityFeat)
		} else if strings.Contains(message, "fix") || strings.Contains(message, "bugfix") {
			highestPriority = comparePriority(highestPriority, PriorityFix)
		}
	}

	return highestPriority
}

func executeGitOperations() {
	if commitFlag {
		if _, err := gitops.Add(); err != nil {
			log.Fatal(err)
		}

		if err := gitops.Commit(fmt.Sprintf(constants.CommitMessage, version.GetLastVersion(), version.ToString()), amend); err != nil {
			log.Fatal(err)
		}
	}

	if tagFlag {
		if err := gitops.CreateTag(fmt.Sprintf(constants.VersionTag, version.ToString()), constants.TagMessage); err != nil {
			log.Fatal(err)
		}
	}

	if pushFlag {
		if err := gitops.Push(); err != nil {
			log.Fatal(err)
		}
	}
}
