package copy

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
)

type SrcToDestinationPaths map[string]string

type Copier func(src, dest string, filesystem Filesystem) error

type VolumeMountCopier struct {
	fileSystem Filesystem
	copier     Copier
}

func NewVolumeMountCopier() *VolumeMountCopier {
	return &VolumeMountCopier{&fileSystem{}, copyFile}
}

// CopyVolumeMount copies all files from the given src path in srcToDest parameter to the given paths.
// It only handles regular files and stops if an error occurs.
// Existing files will be overwritten.
// If the volume was mounted without the subPath Attribute it resolves the data symlink and copies the real files
// from the mount. If the subPath attribute was used it just copies all regular files to the destination.
func (v *VolumeMountCopier) CopyVolumeMount(srcToDest SrcToDestinationPaths) error {
	var multiErr []error
	isSubPathMount := true

	for src, dest := range srcToDest {
		log.Printf("Start copy files from dir %s to %s", src, dest)
		data := filepath.Join(src, "..data")
		log.Printf("Checking data symlink %s", data)
		dataFileInfo, err := v.fileSystem.Lstat(data)

		if err == nil && dataFileInfo.Mode()&os.ModeSymlink != 0 {
			log.Println("Detected data symlink")
			// this volume was mounted without a subPath and all regular files are actually behind symlinks
			// set src to the real folder e.g. src/..2025_05_07_4643786234
			isSubPathMount = false
			var symErr error
			src, symErr = v.resolveDataSymlink(data)
			if symErr != nil {
				return fmt.Errorf("failed to resolve data dir symlink %s: %w", data, err)
			}
		}

		err = v.fileSystem.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				multiErr = append(multiErr, fmt.Errorf("error during filepath walk: %w", err))
				return nil
			}

			multiErr = append(multiErr, v.walk(src, dest, path, isSubPathMount, d))
			return nil
		})

		if err != nil {
			multiErr = append(multiErr, fmt.Errorf("error copy files from directory %s to %s: %w", src, dest, err))
		}
	}
	return errors.Join(multiErr...)
}

// walk will be executed on every path in src by [CopyVolumeMount].
// It only copies regular files and if symLinkChain is not empty it will remove this path from the base path.
// This is needed in volumeMounts from configmaps and secrets without the subPath attributes. In this case
// the files are behind symlinks and the resolved folder is used as source. This path from src to the resolved folder
// should not be copied to the destination.
func (v *VolumeMountCopier) walk(srcVolume, destVolume, filePath string, isSubPathMount bool, d fs.DirEntry) error {
	log.Printf("Processing file %s", filePath)
	if d.IsDir() {
		log.Printf("Skip dir %s", filePath)
		return nil
	}

	sourceFileInfo, err := d.Info()
	if err != nil {
		return err
	}

	if !sourceFileInfo.Mode().IsRegular() {
		return fmt.Errorf("source file %s is not a regular file", filePath)
	}

	var rel string
	if isSubPathMount {
		rel, err = filepath.Rel(srcVolume, filePath)
		if err != nil {
			return fmt.Errorf("can't get the relative path of the source file %s and the source volume %s: %w", filePath, srcVolume, err)
		}
	} else {
		// There can't be nested folders in the mount. Just use the file name from example /mount/..20250504/filename
		// to determine destination path.
		_, rel = path.Split(filePath)
	}

	destinationFilePath := path.Join(destVolume, rel)
	destFileInfo, err := v.fileSystem.Stat(destinationFilePath)
	if err == nil {
		if !destFileInfo.Mode().IsRegular() {
			return fmt.Errorf("destination file %s exists and is not a regular file", destinationFilePath)
		}

		if v.fileSystem.SameFile(sourceFileInfo, destFileInfo) {
			log.Printf("source file %s and destination file %s are equal", filePath, destinationFilePath)
			return nil
		}
	}

	err = v.copier(filePath, destinationFilePath, v.fileSystem)
	if err != nil {
		return err
	}

	return nil
}

// resolveDataSymlink follows the symlink and returns the path from the real file and the relative to the dir of the symlink
func (v *VolumeMountCopier) resolveDataSymlink(symlink string) (string, error) {
	resolvedDataLink, err := v.fileSystem.EvalSymlinks(symlink)
	if err != nil {
		return "", err
	}

	dirInfo, err := v.fileSystem.Stat(resolvedDataLink)
	if err != nil {
		return "", err
	}

	if !dirInfo.IsDir() {
		return "", fmt.Errorf("data symlink %s should point to a dir", symlink)
	}

	return resolvedDataLink, nil
}
