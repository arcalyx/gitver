package cmd

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/gitops"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/spf13/viper"
	"log"
)

func loadConfig() {
	prints("load configuration")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := version.ReadVersion(); err != nil {
		log.Fatal(err)
	}
	prints("load configuration success")
}

func prepareGitOperation() error {
	prints("perpare git operations")
	if err := gitops.ReadRepository(); err != nil {
		return fmt.Errorf("git repository is not initialized: %w", err)
	}

	isClean, err := gitops.IsCleanRepo()
	if err != nil {
		return err
	}

	if !isClean {
		return fmt.Errorf("git repository is not clean")
	}

	prints("prepare git operations success")
	return nil

}

func prints(v ...any) {
	if verbose {
		log.Print(v...)
	}
}
