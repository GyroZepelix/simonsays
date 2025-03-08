package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

type Command interface {
	Name() string
	Description() string
	Execute(ctx *cli.Context) error
}

type CommandRegistry struct {
	commands map[string]Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

func (r *CommandRegistry) Get(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

func (r *CommandRegistry) GetAll() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

type SortByDateCommand struct{}

func (c *SortByDateCommand) Name() string {
	return "sortbydate"
}

func (c *SortByDateCommand) Description() string {
	return "Sort files into directories based on their creation date"
}

func (c *SortByDateCommand) Execute(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return errors.New("please provide a directory path")
	}

	sourcePath := ctx.Args().Get(0)
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("error accessing path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", sourcePath)
	}

	return filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || path == sourcePath {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return fmt.Errorf("cannot get file info for %s: %w", path, err)
		}

		createTime := fileInfo.ModTime()
		dateFolder := createTime.Format("02012006")

		targetDir := filepath.Join(filepath.Dir(path), dateFolder)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		targetPath := filepath.Join(targetDir, d.Name())
		if err := os.Rename(path, targetPath); err != nil {
			return fmt.Errorf("failed to move file %s to %s: %w", path, targetPath, err)
		}

		fmt.Printf("Moved %s to %s\n", path, targetPath)
		return nil
	})
}

type BulkRenameCommand struct{}

func (c *BulkRenameCommand) Name() string {
	return "bulkrename"
}

func (c *BulkRenameCommand) Description() string {
	return "Bulk rename files based on a pattern"
}

func (c *BulkRenameCommand) Execute(ctx *cli.Context) error {
	if ctx.NArg() < 2 {
		return errors.New("please provide a directory path and rename pattern")
	}

	sourcePath := ctx.Args().Get(0)
	pattern := ctx.Args().Get(1)
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("error accessing path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", sourcePath)
	}

	dryRun := ctx.Bool("dry-run")
	startNum := ctx.Int("start")

	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	var files []fs.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry)
		}
	}

	for i, file := range files {
		oldPath := filepath.Join(sourcePath, file.Name())

		ext := filepath.Ext(file.Name())
		base := strings.TrimSuffix(file.Name(), ext)

		newFilename := strings.ReplaceAll(pattern, "{index}", fmt.Sprintf("%d", startNum+i))
		newFilename = strings.ReplaceAll(newFilename, "{name}", base)
		newFilename = strings.ReplaceAll(newFilename, "{ext}", ext)

		if !strings.Contains(newFilename, ".") {
			newFilename = newFilename + ext
		}

		newPath := filepath.Join(sourcePath, newFilename)

		if dryRun {
			fmt.Printf("Would rename: %s -> %s\n", oldPath, newPath)
		} else {
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newPath, err)
			}
			fmt.Printf("Renamed: %s -> %s\n", oldPath, newPath)
		}
	}

	return nil
}

type SortByTypeCommand struct{}

func (c *SortByTypeCommand) Name() string {
	return "sortbytype"
}

func (c *SortByTypeCommand) Description() string {
	return "Sort files into directories based on their file type/extension"
}

func (c *SortByTypeCommand) Execute(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return errors.New("please provide a directory path")
	}

	sourcePath := ctx.Args().Get(0)
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("error accessing path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", sourcePath)
	}

	return filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || path == sourcePath {
			return nil
		}

		ext := strings.TrimPrefix(filepath.Ext(d.Name()), ".")
		if ext == "" {
			ext = "other"
		}

		targetDir := filepath.Join(filepath.Dir(path), ext)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		targetPath := filepath.Join(targetDir, d.Name())
		if err := os.Rename(path, targetPath); err != nil {
			return fmt.Errorf("failed to move file %s to %s: %w", path, targetPath, err)
		}

		fmt.Printf("Moved %s to %s\n", path, targetPath)
		return nil
	})
}

type SortBySizeCommand struct{}

func (c *SortBySizeCommand) Name() string {
	return "sortbysize"
}

func (c *SortBySizeCommand) Description() string {
	return "Sort files into directories based on their size (small, medium, large)"
}

func (c *SortBySizeCommand) Execute(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return errors.New("please provide a directory path")
	}

	sourcePath := ctx.Args().Get(0)
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("error accessing path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", sourcePath)
	}

	smallThreshold := int64(1024 * 1024)       // 1MB
	mediumThreshold := int64(10 * 1024 * 1024) // 10MB

	return filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || path == sourcePath {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return fmt.Errorf("cannot get file info for %s: %w", path, err)
		}

		var sizeCategory string
		if fileInfo.Size() < smallThreshold {
			sizeCategory = "small"
		} else if fileInfo.Size() < mediumThreshold {
			sizeCategory = "medium"
		} else {
			sizeCategory = "large"
		}

		targetDir := filepath.Join(filepath.Dir(path), sizeCategory)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		targetPath := filepath.Join(targetDir, d.Name())
		if err := os.Rename(path, targetPath); err != nil {
			return fmt.Errorf("failed to move file %s to %s: %w", path, targetPath, err)
		}

		fmt.Printf("Moved %s to %s\n", path, targetPath)
		return nil
	})
}

type ListCommand struct{}

func (c *ListCommand) Name() string {
	return "list"
}

func (c *ListCommand) Description() string {
	return "List files with various sorting options"
}

func (c *ListCommand) Execute(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return errors.New("please provide a directory path")
	}

	sourcePath := ctx.Args().Get(0)
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("error accessing path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", sourcePath)
	}

	sortBy := ctx.String("sort")
	recursive := ctx.Bool("recursive")

	var files []os.FileInfo

	if recursive {
		err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, info)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}
	} else {
		entries, err := os.ReadDir(sourcePath)
		if err != nil {
			return fmt.Errorf("error reading directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				info, err := entry.Info()
				if err != nil {
					return fmt.Errorf("error getting file info: %w", err)
				}
				files = append(files, info)
			}
		}
	}

	switch sortBy {
	case "name":
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})
	case "size":
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size() < files[j].Size()
		})
	case "time":
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Before(files[j].ModTime())
		})
	default:
		return fmt.Errorf("unknown sort option: %s", sortBy)
	}

	fmt.Println("Files:")
	for _, file := range files {
		fmt.Printf("- %s (Size: %d bytes, Modified: %s)\n",
			file.Name(),
			file.Size(),
			file.ModTime().Format(time.RFC3339))
	}

	return nil
}

func main() {
	registry := NewCommandRegistry()

	registry.Register(&SortByDateCommand{})
	registry.Register(&BulkRenameCommand{})
	registry.Register(&SortByTypeCommand{})
	registry.Register(&SortBySizeCommand{})
	registry.Register(&ListCommand{})

	app := &cli.App{
		Name:     "simonsays",
		Usage:    "A swiss-knife tool for file and directory manipulation",
		Commands: []*cli.Command{},
	}

	for _, cmd := range registry.GetAll() {
		command := cmd
		cliCommand := &cli.Command{
			Name:  command.Name(),
			Usage: command.Description(),
			Action: func(c *cli.Context) error {
				return command.Execute(c)
			},
		}

		switch command.Name() {
		case "bulkrename":
			cliCommand.Flags = []cli.Flag{
				&cli.BoolFlag{
					Name:  "dry-run",
					Usage: "Show what would be renamed without actually renaming",
				},
				&cli.IntFlag{
					Name:  "start",
					Usage: "Starting index for {index} pattern",
					Value: 1,
				},
			}
		case "list":
			cliCommand.Flags = []cli.Flag{
				&cli.StringFlag{
					Name:  "sort",
					Usage: "Sort by (name, size, time)",
					Value: "name",
				},
				&cli.BoolFlag{
					Name:  "recursive",
					Usage: "List files recursively",
				},
			}
		}

		app.Commands = append(app.Commands, cliCommand)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
