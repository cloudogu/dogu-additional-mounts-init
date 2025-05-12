package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/cloudogu/dogu-data-seeder/internal/copy"
	"log"
	"os"
	"strings"
)

var (
	copyCmd = flag.NewFlagSet("copy", flag.ExitOnError)
)

type volumeCopier interface {
	CopyVolumeMount(srcToDest copy.SrcToDestinationPaths) error
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("expected at least on of the following commands: \n"+
			"%s - copy files from specified volumes to destination paths", copyCmd.Name())
	}

	var err error
	switch os.Args[1] {
	case copyCmd.Name():
		err = handleCopyCommand(os.Args[2:], copy.NewVolumeMountCopier())
	default:
		err = errors.New("unknown command")
		// print help
	}

	if err != nil {
		log.Fatal(err.Error())
	}
}

func handleCopyCommand(args []string, copier volumeCopier) error {
	var sourcePaths stringSliceFlag
	var targetPaths stringSliceFlag
	copyCmd.Var(&sourcePaths, "source", "")
	copyCmd.Var(&targetPaths, "target", "")
	err := copyCmd.Parse(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	if len(sourcePaths) != len(targetPaths) {
		return fmt.Errorf("amount of source and target paths aren't equal")
	}

	if len(sourcePaths) == 0 && len(targetPaths) == 0 {
		log.Println("no source and target paths given")
		return nil
	}

	copyMap := make(copy.SrcToDestinationPaths, len(sourcePaths))
	for i := range sourcePaths {
		copyMap[sourcePaths[i]] = targetPaths[i]
	}

	err = copier.CopyVolumeMount(copyMap)
	if err != nil {
		return err
	}

	return nil
}

type stringSliceFlag []string

func (i *stringSliceFlag) String() string {
	return strings.Join(*i, ",")
}

func (i *stringSliceFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
