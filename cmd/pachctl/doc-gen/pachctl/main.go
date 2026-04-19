// pachctl-doc

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/laityjet/mammoth/v0/internal/cmdutil"
	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/log"
	"github.com/laityjet/mammoth/v0/internal/pachconfig"
	"github.com/laityjet/mammoth/v0/internal/server/cmd/pachctl/cmd"

	"github.com/spf13/cobra/doc"
)

func main() {
	log.InitPachctlLogger()
	cmdutil.Main(context.Background(), do, &pachconfig.EmptyConfig{})
}

func do(ctx context.Context, _ *pachconfig.EmptyConfig) error {
	path := "./docs/"

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return errors.Wrap(err, "make output directory")
	}

	rootCmd, err := cmd.PachctlCmd()
	if err != nil {
		return errors.Wrap(err, "generate pachctl command")
	}
	rootCmd.DisableAutoGenTag = true

	const fmTemplate = `---
date: %s
title: "%s"
description: "Learn about the %s command"
---

`

	filePrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, filepath.Ext(name))
		title := strings.Replace(base, "_", " ", -1)
		description := strings.Replace(base, "_", " ", -1)
		return fmt.Sprintf(fmTemplate, now, title, description)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, filepath.Ext(name))
		url := "../" + strings.ToLower(base)
		return url
	}

	err = doc.GenMarkdownTreeCustom(rootCmd, path, filePrepender, linkHandler)

	if err != nil {
		return errors.Wrap(err, "generate Markdown documentation")
	}

	if err := reformatMarkdownOutput(path); err != nil {
		return errors.Wrap(err, "replace ./ in Markdown files")
	}

	return nil
}

func reformatMarkdownOutput(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return errors.Wrap(err, "read directory")
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			filePath := filepath.Join(path, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			// Replace any instance of "./" in the generated markdown files
			updatedContent := strings.Replace(string(content), "./pachctl", "pachctl", -1)
			// Replace any instance of "	-" in the generated markdown examples
			updatedContent2 := strings.Replace(updatedContent, "	-", "", -1)

			// Write the updated content back to the file
			if err := os.WriteFile(filePath, []byte(updatedContent2), os.ModePerm); err != nil {
				return errors.Wrap(err, "write file")
			}
		}
	}

	return nil
}
