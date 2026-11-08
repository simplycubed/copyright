package main

import (
	"context"
	"fmt"
	"os"

	"github.com/simplycubed/copyright/internal/config"
	"github.com/simplycubed/copyright/internal/processor"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
	quiet   bool
	dryRun  bool
	backup  bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "copyright",
	Short: "A tool for managing copyright headers in source files",
	Long: `Copyright is a CLI tool that checks and adds copyright headers
to source code files. It supports multiple programming languages
and can be configured via a .copyright.yaml file.`,
}

var checkCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Check files for copyright headers",
	Long:  `Check all files in the specified path (default: current directory) for copyright headers.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCheck,
}

var fixCmd = &cobra.Command{
	Use:   "fix [path]",
	Short: "Add missing copyright headers to files",
	Long:  `Add copyright headers to files that are missing them in the specified path (default: current directory).`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runFix,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a .copyright.yaml configuration file",
	Long:  `Create a new .copyright.yaml configuration file with default values in the current directory.`,
	RunE:  runInit,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("copyright v0.1.0")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is .copyright.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-error output")

	fixCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be changed without modifying files")
	fixCmd.Flags().BoolVar(&backup, "backup", false, "create .bak files before modifying")

	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(fixCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		// Use defaults if config not found
		cfg = config.NewDefault()
	}

	// Apply CLI flags
	cfg.Output.Verbose = verbose
	cfg.Output.Quiet = quiet

	return cfg, nil
}

func runCheck(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	proc, err := processor.NewProcessor(cfg)
	if err != nil {
		return fmt.Errorf("failed to create processor: %w", err)
	}

	ctx := context.Background()
	result, err := proc.Process(ctx, path, processor.ModeCheck)
	if err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}

	printResult(result, cfg)
	os.Exit(result.ExitCode())
	return nil
}

func runFix(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	proc, err := processor.NewProcessor(cfg)
	if err != nil {
		return fmt.Errorf("failed to create processor: %w", err)
	}

	ctx := context.Background()
	opts := &processor.ProcessorOptions{
		DryRun: dryRun,
		Backup: backup,
	}
	result, err := proc.ProcessWithOptions(ctx, path, processor.ModeFix, opts)
	if err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}

	printResult(result, cfg)
	return nil
}

func runInit(cmd *cobra.Command, args []string) error {
	configPath := ".copyright.yaml"

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file %s already exists", configPath)
	}

	cfg := config.NewDefault()
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Printf("Created %s\n", configPath)
	return nil
}

func printResult(result *processor.ProcessResult, cfg *config.Config) {
	if cfg.Output.Quiet {
		return
	}

	if cfg.Output.Verbose {
		for _, fr := range result.FileResults {
			symbol := getStatusSymbol(fr.Status)
			fmt.Printf("%s %s: %s\n", symbol, fr.RelativePath, fr.Message)
		}
		fmt.Println()
	}

	fmt.Println(result.Summary())

	if result.LicenseStatus != "" {
		fmt.Printf("License: %s\n", result.LicenseStatus)
	}
}

func getStatusSymbol(status processor.FileStatus) string {
	switch status {
	case processor.StatusHasCopyright:
		return "[OK]"
	case processor.StatusMissingCopyright:
		return "[MISSING]"
	case processor.StatusFixed:
		return "[FIXED]"
	case processor.StatusSkipped:
		return "[SKIP]"
	case processor.StatusFailed:
		return "[FAIL]"
	default:
		return "[?]"
	}
}
