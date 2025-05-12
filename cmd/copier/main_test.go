package main

import (
	"flag"
	"github.com/cloudogu/dogu-data-seeder/internal/copy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_handleCopyCommand(t *testing.T) {
	t.Run("should call copy subsequent src to destination volumes", func(t *testing.T) {
		// given
		oldCopyCmd := copyCmd
		copyCmd = flag.NewFlagSet("copy", flag.ExitOnError)
		defer func() { copyCmd = oldCopyCmd }()
		volumeCopierMock := newMockVolumeCopier(t)
		args := []string{"--source=/src1", "--target=/target1", "--source=/src2", "--target=/target2"}
		expectedCopyMap := copy.SrcToDestinationPaths{"/src1": "/target1", "/src2": "/target2"}
		volumeCopierMock.EXPECT().CopyVolumeMount(expectedCopyMap).Return(nil)

		// when
		err := handleCopyCommand(args, volumeCopierMock)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on copy error error", func(t *testing.T) {
		// given
		oldCopyCmd := copyCmd
		copyCmd = flag.NewFlagSet("copy", flag.ExitOnError)
		defer func() { copyCmd = oldCopyCmd }()
		volumeCopierMock := newMockVolumeCopier(t)
		args := []string{"--source=/src1", "--target=/target1", "--source=/src2", "--target=/target2"}
		expectedCopyMap := copy.SrcToDestinationPaths{"/src1": "/target1", "/src2": "/target2"}
		volumeCopierMock.EXPECT().CopyVolumeMount(expectedCopyMap).Return(assert.AnError)

		// when
		err := handleCopyCommand(args, volumeCopierMock)

		// then
		require.Error(t, err)
	})

	t.Run("should return nil on empty parameter", func(t *testing.T) {
		// given
		oldCopyCmd := copyCmd
		copyCmd = flag.NewFlagSet("copy", flag.ExitOnError)
		defer func() { copyCmd = oldCopyCmd }()
		volumeCopierMock := newMockVolumeCopier(t)
		var args []string

		// when
		err := handleCopyCommand(args, volumeCopierMock)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on odd parameters", func(t *testing.T) {
		// given
		oldCopyCmd := copyCmd
		copyCmd = flag.NewFlagSet("copy", flag.ExitOnError)
		defer func() { copyCmd = oldCopyCmd }()
		volumeCopierMock := newMockVolumeCopier(t)
		args := []string{"--source=/src1", "--target=/target1", "--source=/src2"}

		// when
		err := handleCopyCommand(args, volumeCopierMock)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "amount of source and target paths aren't equal")
	})
}
