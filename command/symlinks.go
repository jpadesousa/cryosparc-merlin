package command

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strings"
// )

// const (
// 	Bold  = "\033[1m"
// 	Reset = "\033[0m"
// )

// type Project struct {
// 	UID string
// }

// type Job struct {
// 	UID string
// }

// type Symlink struct {
// 	LinkTarget string
// }

// type CryoCLI interface {
// 	ListProjects() ([]Project, error)
// 	ListJobs(projectUID string) ([]Job, error)
// 	GetJobSymlinks(projectUID, jobUID string) ([]Symlink, error)
// 	JobImportReplaceSymlinks(projectUID, jobUID, prefixCut, prefixNew string) (int, error)
// }

// // Replace this with your implementation.
// var cli CryoCLI

// type JobRef struct {
// 	ProjectUID string
// 	JobUID     string
// }

// func main() {
// 	reader := bufio.NewReader(os.Stdin)

// 	fmt.Println("🧩 This script will loop over CryoSPARC jobs and attempt to repair symlinks.")
// 	fmt.Println("🔁 It replaces path prefixes in imported file links using JobImportReplaceSymlinks.")
// 	fmt.Println("🛑 You will be asked to confirm before any changes are made.")
// 	fmt.Println()

// 	// ------------------------------------------------------------------
// 	// Step 1: Project selection
// 	// ------------------------------------------------------------------

// 	projects, err := cli.ListProjects()
// 	if err != nil {
// 		fmt.Printf("❌ Failed to list projects: %v\n", err)
// 		return
// 	}

// 	projectUIDs := make(map[string]bool)

// 	fmt.Println("Available projects:")
// 	for _, p := range projects {
// 		projectUIDs[p.UID] = true
// 		fmt.Printf("  - %s\n", p.UID)
// 	}

// 	projectChoice := prompt(reader, "\nEnter a project UID to process, or 'all' to process all projects: ")

// 	if projectChoice != "all" && !projectUIDs[projectChoice] {
// 		fmt.Println("❌ Invalid project UID. Exiting.")
// 		return
// 	}

// 	var selectedProjects []string
// 	if projectChoice == "all" {
// 		for _, p := range projects {
// 			selectedProjects = append(selectedProjects, p.UID)
// 		}
// 	} else {
// 		selectedProjects = []string{projectChoice}
// 	}

// 	var selectedJobs []JobRef

// 	// ------------------------------------------------------------------
// 	// Step 2: Job selection
// 	// ------------------------------------------------------------------

// 	if len(selectedProjects) == 1 {

// 		projectUID := selectedProjects[0]

// 		jobs, err := cli.ListJobs(projectUID)
// 		if err != nil {
// 			fmt.Printf("❌ Could not get jobs for project %s: %v\n", projectUID, err)
// 			return
// 		}

// 		jobUIDs := make(map[string]bool)

// 		fmt.Printf("\nJobs for project %s:\n", projectUID)
// 		for _, job := range jobs {
// 			jobUIDs[job.UID] = true
// 			fmt.Printf("  - %s\n", job.UID)
// 		}

// 		jobChoice := prompt(reader, "\nEnter a job UID to process, or 'all' to process all jobs: ")

// 		if jobChoice != "all" && !jobUIDs[jobChoice] {
// 			fmt.Println("❌ Invalid job UID. Exiting.")
// 			return
// 		}

// 		if jobChoice == "all" {
// 			for _, job := range jobs {
// 				selectedJobs = append(selectedJobs, JobRef{
// 					ProjectUID: projectUID,
// 					JobUID:     job.UID,
// 				})
// 			}
// 		} else {
// 			selectedJobs = append(selectedJobs, JobRef{
// 				ProjectUID: projectUID,
// 				JobUID:     jobChoice,
// 			})
// 		}

// 	} else {

// 		for _, projectUID := range selectedProjects {

// 			jobs, err := cli.ListJobs(projectUID)
// 			if err != nil {
// 				fmt.Printf("⚠️ Could not get jobs for project %s: %v\n", projectUID, err)
// 				continue
// 			}

// 			for _, job := range jobs {
// 				selectedJobs = append(selectedJobs, JobRef{
// 					ProjectUID: projectUID,
// 					JobUID:     job.UID,
// 				})
// 			}
// 		}
// 	}

// 	// ------------------------------------------------------------------
// 	// Step 3: Symlink summary
// 	// ------------------------------------------------------------------

// 	fmt.Println("\n🔍 Analyzing current symlink targets...")

// 	// root -> jobs
// 	symlinkRoots := make(map[string]map[JobRef]struct{})

// 	for _, ref := range selectedJobs {

// 		symlinks, err := cli.GetJobSymlinks(ref.ProjectUID, ref.JobUID)
// 		if err != nil {
// 			fmt.Printf("⚠️ Could not get symlinks for %s %s: %v\n",
// 				ref.ProjectUID, ref.JobUID, err)
// 			continue
// 		}

// 		for _, link := range symlinks {

// 			if link.LinkTarget == "" {
// 				continue
// 			}

// 			root := dirname(link.LinkTarget)

// 			if _, ok := symlinkRoots[root]; !ok {
// 				symlinkRoots[root] = make(map[JobRef]struct{})
// 			}

// 			symlinkRoots[root][ref] = struct{}{}
// 		}
// 	}

// 	if len(symlinkRoots) > 0 {

// 		fmt.Println("\n🔗 Current symlink target prefixes and associated jobs:")

// 		for root, jobs := range symlinkRoots {

// 			fmt.Printf("\n%s%s%s\n", Bold, root, Reset)

// 			for job := range jobs {
// 				fmt.Printf("  - %s %s\n", job.ProjectUID, job.JobUID)
// 			}
// 		}

// 	} else {
// 		fmt.Println("⚠️ No symlinks found or accessible in the selected jobs.")
// 	}

// 	// ------------------------------------------------------------------
// 	// Step 4: Prefixes
// 	// ------------------------------------------------------------------

// 	fmt.Println()

// 	prefixCut := prompt(reader, "Enter old prefix to remove (prefix_cut): ")
// 	prefixNew := prompt(reader, "Enter new prefix to insert (prefix_new): ")

// 	if prefixCut == "" || prefixNew == "" {
// 		fmt.Println("❌ Both prefix_cut and prefix_new must be provided. Exiting.")
// 		return
// 	}

// 	cutSlash := strings.HasSuffix(prefixCut, "/")
// 	newSlash := strings.HasSuffix(prefixNew, "/")

// 	if cutSlash != newSlash {
// 		fmt.Println("⚠️ Warning: The trailing slash on prefix_cut and prefix_new differs.")
// 		fmt.Printf("  prefix_cut ends with '/'? %v\n", cutSlash)
// 		fmt.Printf("  prefix_new ends with '/'? %v\n", newSlash)
// 		fmt.Println("  This might cause unexpected mismatches during replacement.")
// 		fmt.Println()
// 	}

// 	fmt.Printf("\n📂 prefix_cut: %s\n", prefixCut)
// 	fmt.Printf("📂 prefix_new: %s\n", prefixNew)

// 	fmt.Printf("\nYou are about to process %d job(s).\n", len(selectedJobs))

// 	confirm := strings.ToLower(prompt(reader, "Are you sure you want to apply these changes? [y/N]: "))

// 	if confirm != "y" {
// 		fmt.Println("❌ Aborted by user.")
// 		return
// 	}

// 	// ------------------------------------------------------------------
// 	// Step 5: Repair
// 	// ------------------------------------------------------------------

// 	var failedJobs []JobRef

// 	for _, ref := range selectedJobs {

// 		fmt.Printf("🔧 Repairing %s %s\n", ref.ProjectUID, ref.JobUID)

// 		modified, err := cli.JobImportReplaceSymlinks(
// 			ref.ProjectUID,
// 			ref.JobUID,
// 			prefixCut,
// 			prefixNew,
// 		)

// 		if err != nil {
// 			failedJobs = append(failedJobs, ref)
// 			fmt.Printf("❌ Failed to repair %s %s: %v\n\n",
// 				ref.ProjectUID, ref.JobUID, err)
// 			continue
// 		}

// 		if modified > 0 {
// 			fmt.Printf("%s✅ Finished. Modified %d links.%s\n\n",
// 				Bold, modified, Reset)
// 		} else {
// 			fmt.Printf("✅ Finished. Modified %d links.\n\n", modified)
// 		}
// 	}

// 	fmt.Printf("🎯 Completed. %d jobs failed to repair.\n", len(failedJobs))

// 	if len(failedJobs) > 0 {
// 		fmt.Println("🔻 Failed jobs:")
// 		for _, job := range failedJobs {
// 			fmt.Printf("  - %s %s\n", job.ProjectUID, job.JobUID)
// 		}
// 	}
// }

// func prompt(reader *bufio.Reader, text string) string {
// 	fmt.Print(text)
// 	input, _ := reader.ReadString('\n')
// 	return strings.TrimSpace(input)
// }

// func dirname(path string) string {
// 	i := strings.LastIndex(path, "/")
// 	if i == -1 {
// 		return path
// 	}
// 	return path[:i]
// }
