package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	slogformatter "github.com/samber/slog-formatter"
	"github.com/spf13/cobra"

	pkg "github.com/andrew-womeldorf/mise-task-defs/go/internal"
	"github.com/andrew-womeldorf/mise-task-defs/go/internal/buildinfo"
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	verbose  bool
	jsonLogs bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "api",
	Short: "Example Connect RPC API CLI",
	Long: `Example Connect RPC API CLI provides commands for managing users.
It provides both server functionality and RPC client commands.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Configure logger based on verbose flag
		logLevel := slog.LevelInfo
		if verbose {
			logLevel = slog.LevelDebug
		}

		// Configure logger
		var handler slog.Handler
		if jsonLogs {
			handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: true,
				Level:     logLevel,
			})
		} else {
			handler = tint.NewHandler(os.Stderr, &tint.Options{
				AddSource: true,
				Level:     logLevel,
			})
		}

		logger := slog.New(slogformatter.NewFormatterHandler(
			slogformatter.ErrorFormatter("error"),
		)(handler))
		slog.SetDefault(logger)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		if err := cmd.Help(); err != nil {
			slog.Error("Failed to display help", "error", err)
			os.Exit(1)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version info",
	Long:  `Version info.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s %s\n", buildinfo.Version, buildinfo.Commit, buildinfo.Date)
	},
}

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "User management commands",
	Long:  `Commands for managing users in the database.`,
}

var createUserCmd = &cobra.Command{
	Use:   "create <email> <username>",
	Short: "Create a new user",
	Long:  `Create a new user with the specified email and username.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dbPath, _ := cmd.Flags().GetString("db")
		db, err := pkg.NewDatabase(dbPath)
		if err != nil {
			slog.Error("Failed to open database", "error", err)
			os.Exit(1)
		}
		defer func() {
			if err := db.Close(); err != nil {
				slog.Error("Failed to close database", "error", err)
			}
		}()

		email, username := args[0], args[1]
		ctx := context.Background()
		
		if err := db.CreateUser(ctx, email, username); err != nil {
			slog.Error("Failed to create user", "error", err)
			os.Exit(1)
		}
		
		fmt.Printf("User created successfully: %s (%s)\n", username, email)
	},
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Long:  `List all users in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbPath, _ := cmd.Flags().GetString("db")
		limit, _ := cmd.Flags().GetInt64("limit")
		
		db, err := pkg.NewDatabase(dbPath)
		if err != nil {
			slog.Error("Failed to open database", "error", err)
			os.Exit(1)
		}
		defer func() {
			if err := db.Close(); err != nil {
				slog.Error("Failed to close database", "error", err)
			}
		}()

		ctx := context.Background()
		users, err := db.GetUsers(ctx, limit)
		if err != nil {
			slog.Error("Failed to get users", "error", err)
			os.Exit(1)
		}
		
		fmt.Printf("Found %d users:\n", len(users))
		for _, user := range users {
			fmt.Printf("- %s (%s) [ID: %s]\n", user.Username, user.Email, user.ID)
		}
	},
}

func init() {
	// Add persistent flags that will be available to all commands
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	RootCmd.PersistentFlags().BoolVar(&jsonLogs, "json", false, "Output logs in JSON format (default: text)")
	
	// Add database flags
	usersCmd.PersistentFlags().String("db", "users.db", "Path to SQLite database file")
	listUsersCmd.Flags().Int64("limit", 10, "Maximum number of users to list")
	
	// Add commands
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(createUserCmd)
	usersCmd.AddCommand(listUsersCmd)
}
