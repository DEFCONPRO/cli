package update

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pio "github.com/confluentinc/cli/v3/pkg/io"
	"github.com/confluentinc/cli/v3/pkg/mock"
	updateMock "github.com/confluentinc/cli/v3/pkg/update/mock"
	"github.com/confluentinc/cli/v3/test"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		params *ClientParams
		want   *client
	}{
		{
			name:   "should set default values (interval=24h, clock=real clock, fs=real fs, os=real os)",
			params: &ClientParams{},
			want: &client{
				ClientParams: &ClientParams{CheckInterval: 24 * time.Hour, OS: runtime.GOOS},
				clock:        clockwork.NewRealClock(),
				fs:           &pio.RealFileSystem{},
			},
		},
		{
			name: "should set provided values",
			params: &ClientParams{
				CheckInterval: 48 * time.Hour,
				OS:            "duckduckgoos",
				DisableCheck:  true,
			},
			want: &client{
				ClientParams: &ClientParams{
					CheckInterval: 48 * time.Hour,
					OS:            "duckduckgoos",
					DisableCheck:  true,
				},
				clock: clockwork.NewRealClock(),
				fs:    &pio.RealFileSystem{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := NewClient(test.params); !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewClient() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestCheckForUpdates(t *testing.T) {
	tmpCheckFile1, err := os.CreateTemp("", "cli-test1-")
	require.NoError(t, err)
	defer os.Remove(tmpCheckFile1.Name())

	type args struct {
		name           string
		currentVersion string
		forceCheck     bool
	}
	tests := []struct {
		name      string
		client    *client
		args      args
		wantMajor string
		wantMinor string
		wantErr   bool
	}{
		{
			name: "should err if currentVersion isn't semver",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{},
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "gobbledegook",
			},
			wantErr: true,
		},
		{
			name: "should err if can't get versions",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						return nil, nil, fmt.Errorf("zap")
					},
				},
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
			},
			wantErr: true,
		},
		{
			name: "should return the new version",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v3, _ := version.NewSemver("v3")
						return v3, current, nil
					},
				},
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
			},
			wantMajor: "v3",
		},
		{
			name: "should not check for the new version if has checked recently",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v3, _ := version.NewSemver("v3")
						return v3, v3, nil
					},
				},
				CheckFile: tmpCheckFile1.Name(),
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
			},
		},
		{
			name: "should not check again if checked recently",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						require.Fail(t, "Shouldn't be called")
						return nil, nil, fmt.Errorf("whoops")
					},
				},
				// This check file was created by the TmpFile process, modtime is current, so should skip check
				CheckFile: tmpCheckFile1.Name(),
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
			},
		},
		{
			name: "should respect forceCheck even if you checked recently",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v3, _ := version.NewSemver("v3")
						return v3, current, nil
					},
				},
				// This check file was created by the TmpFile process, modtime is current, so should skip check
				CheckFile: tmpCheckFile1.Name(),
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
				forceCheck:     true,
			},
			wantMajor: "v3",
		},
		{
			name: "should err if you can't create the CheckFile",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v2, _ := version.NewSemver("v2")
						return v2, v2, nil
					},
				},
				// This file doesn't exist but you won't have permission to create it
				CheckFile: "/sbin/cant-write-here",
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
			},
			wantErr: true,
		},
		{
			name: "should err if you can't touch the CheckFile",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v2, _ := version.NewSemver("v2")
						return v2, v2, nil
					},
				},
				// This file doesn't exist but you won't have permission to touch it
				CheckFile: "/sbin/ping",
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
			},
			wantErr: true,
		},
		{
			name: "should not check if disabled",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						require.Fail(t, "Shouldn't be called")
						return nil, nil, fmt.Errorf("whoops")
					},
				},
				DisableCheck: true,
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
			},
		},
		{
			name: "checks - error",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						return nil, nil, fmt.Errorf("whoops")
					},
				},
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
			},
			wantErr: true,
		},
		{
			name: "checks - success - update",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v, _ := version.NewVersion("v1.2.4")
						return nil, v, nil
					},
				},
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.3",
			},
			wantMinor: "v1.2.4",
		},
		{
			name: "checks - success - same version",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v, _ := version.NewVersion("v1.2.4")
						return v, v, nil
					},
				},
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v1.2.4",
			},
		},
		{
			name: "checks - success - hyphen no update",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v, _ := version.NewVersion("v0.238.0")
						return v, v, nil
					},
				},
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v0.238.0-7-g5060ef4",
			},
		},
		{
			name: "checks - success - hyphen same version",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v, _ := version.NewVersion("v0.238.0-7-g5060ef4")
						return v, v, nil
					},
				},
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v0.238.0-7-g5060ef4",
			},
		},
		{
			name: "checks - success - hyphen update",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
						v, _ := version.NewVersion("v0.238.0-7-g5060ef4")
						return nil, v, nil
					},
				},
			}),
			args: args{
				name:           "my-cli",
				currentVersion: "v0.238.0",
			},
			wantMinor: "v0.238.0-7-g5060ef4",
		},
		{
			name: "different major and minor versions",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(_ string, current *version.Version) (*version.Version, *version.Version, error) {
						v0, _ := version.NewVersion("v0.1.0")
						v1, _ := version.NewVersion("v1.0.0")
						return v1, v0, nil
					},
				},
			}),
			args:      args{currentVersion: "v0.0.1"},
			wantMajor: "v1.0.0",
			wantMinor: "v0.1.0",
		},
		{
			name: "no latest major or minor versions",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestMajorAndMinorVersionFunc: func(_ string, current *version.Version) (*version.Version, *version.Version, error) {
						v0, _ := version.NewVersion("v0.0.0")
						return v0, v0, nil
					},
				},
			}),
			args:      args{currentVersion: "v0.0.0"},
			wantMajor: "",
			wantMinor: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			major, minor, err := test.client.CheckForUpdates(test.args.name, test.args.currentVersion, test.args.forceCheck)
			if (err != nil) != test.wantErr {
				t.Errorf("client.CheckForUpdates() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if major != test.wantMajor {
				t.Errorf("client.CheckForUpdates() major = %v, want %v", major, test.wantMajor)
			}
			if minor != test.wantMinor {
				t.Errorf("client.CheckForUpdates() minor = %v, want %v", minor, test.wantMinor)
			}
		})
	}
}

func TestCheckForUpdates_BehaviorOverTime(t *testing.T) {
	req := require.New(t)

	tmpDir, err := os.MkdirTemp("", "cli-test3-")
	req.NoError(err)
	defer os.RemoveAll(tmpDir)
	checkFile := filepath.FromSlash(fmt.Sprintf("%s/new-check-file", tmpDir))

	repo := &updateMock.Repository{
		GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
			v3, _ := version.NewSemver("v3")
			return v3, v3, nil
		},
	}
	clock := clockwork.NewFakeClockAt(time.Now())
	client := NewClient(&ClientParams{
		Repository: repo,
		CheckFile:  checkFile,
	})
	client.clock = clock

	// Should check and find update
	latestMajorVersion, latestMinorVersion, err := client.CheckForUpdates("my-cli", "v1.2.3", false)
	req.NoError(err)
	req.Equal("v3", latestMajorVersion)
	req.Equal("v3", latestMinorVersion)
	req.True(repo.GetLatestMajorAndMinorVersionCalled())

	// Shouldn't check anymore for 24 hours
	for i := 0; i < 3; i++ {
		clock.Advance(8*time.Hour + -1*time.Second)
		repo.Reset()

		_, _, _ = client.CheckForUpdates("my-cli", "v1.2.3", false)
		req.False(repo.GetLatestMajorAndMinorVersionCalled())
	}

	// 5 days pass...
	clock.Advance(5 * 24 * time.Hour)

	// Should check and find update
	latestMajorVersion, latestMinorVersion, err = client.CheckForUpdates("my-cli", "v1.2.3", false)
	req.NoError(err)
	req.Equal("v3", latestMajorVersion)
	req.Equal("v3", latestMinorVersion)
	req.True(repo.GetLatestMajorAndMinorVersionCalled())

	// Shouldn't check anymore for 24 hours
	for i := 0; i < 3; i++ {
		clock.Advance(8*time.Hour + -1*time.Second)
		repo.Reset()

		_, _, _ = client.CheckForUpdates("my-cli", "v1.2.3", false)
		req.False(repo.GetLatestMajorAndMinorVersionCalled())
	}

	// Finally we should check once more
	clock.Advance(3 * time.Second)
	repo.Reset()
	_, _, _ = client.CheckForUpdates("my-cli", "v1.2.3", false)
	req.True(repo.GetLatestMajorAndMinorVersionCalled())
}

func TestCheckForUpdates_NoCheckFileGiven(t *testing.T) {
	req := require.New(t)

	repo := &updateMock.Repository{
		GetLatestMajorAndMinorVersionFunc: func(name string, current *version.Version) (*version.Version, *version.Version, error) {
			v3, _ := version.NewSemver("v3")
			return v3, v3, nil
		},
	}
	client := NewClient(&ClientParams{
		Repository: repo,
	})
	client.clock = clockwork.NewFakeClockAt(time.Now())

	// Should check for updates every time if no CheckFile given to serve as the "last check" cache
	for i := 0; i < 3; i++ {
		latestMajorVersion, latestMinorVersion, err := client.CheckForUpdates("my-cli", "v1.2.3", false)
		req.NoError(err)
		req.Equal("v3", latestMajorVersion)
		req.Equal("v3", latestMinorVersion)
		req.True(repo.GetLatestMajorAndMinorVersionCalled())
		repo.Reset()
	}
}

func TestDownloadChecksum(t *testing.T) {
	checksums := test.LoadFixture(t, "../input/update/checksums.golden")

	mockRepository := &updateMock.Repository{
		DownloadChecksumsFunc: func(name, version string) (string, error) {
			if version == "2.5.1" {
				return checksums, nil
			} else {
				return "", fmt.Errorf("no checksums for given version")
			}
		},
	}

	tests := []struct {
		name            string
		version         string
		wantDownloadErr bool
	}{
		{
			name:            "valid checksum for valid version verifies successfully",
			version:         "2.5.1",
			wantDownloadErr: false,
		},
		{
			name:            "invalid checksum for valid version fails",
			version:         "2.5.1",
			wantDownloadErr: false,
		},
		{
			name:            "checksum for invalid version fails",
			version:         "0.1234.0",
			wantDownloadErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := mockRepository.DownloadChecksums("confluent", test.version)
			if test.wantDownloadErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetLatestReleaseNotes(t *testing.T) {
	currentVersion := "0.1.0"
	releaseNotesVersion := "1.0.0"
	releaseNotes := "nice release notes"

	tests := []struct {
		name             string
		client           *client
		wantVersion      string
		wantReleaseNotes []string
		wantErr          bool
	}{
		{
			name: "success",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestReleaseNotesVersionsFunc: func(_, _ string) (version.Collection, error) {
						v, _ := version.NewSemver(releaseNotesVersion)
						return version.Collection{v}, nil
					},
					DownloadReleaseNotesFunc: func(_, _ string) (string, error) {
						return releaseNotes, nil
					},
				},
			}),
			wantVersion:      releaseNotesVersion,
			wantReleaseNotes: []string{releaseNotes},
			wantErr:          false,
		},
		{
			name: "error getting release notes version",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestReleaseNotesVersionsFunc: func(_, _ string) (version.Collection, error) {
						return nil, fmt.Errorf("whoops")
					},
					DownloadReleaseNotesFunc: func(_, _ string) (string, error) {
						return "", nil
					},
				},
			}),
			wantErr: true,
		},
		{
			name: "error downloading release notes",
			client: NewClient(&ClientParams{
				Repository: &updateMock.Repository{
					GetLatestReleaseNotesVersionsFunc: func(_, _ string) (version.Collection, error) {
						v1, _ := version.NewSemver("v1")
						return version.Collection{v1}, nil
					},
					DownloadReleaseNotesFunc: func(_, _ string) (string, error) {
						return "", fmt.Errorf("whoops")
					},
				},
			}),
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotReleaseNotesVersion, gotReleaseNotes, err := test.client.GetLatestReleaseNotes("confluent", currentVersion)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.wantVersion, gotReleaseNotesVersion)
			require.Equal(t, test.wantReleaseNotes, gotReleaseNotes)
		})
	}
}

func TestUpdateBinary(t *testing.T) {
	req := require.New(t)

	binName := "fake_cli"

	installDir, err := os.MkdirTemp("", "cli-test4-")
	require.NoError(t, err)
	defer os.Remove(installDir)

	err = os.WriteFile(filepath.Join(installDir, binName), []byte("old version"), os.ModePerm)
	require.NoError(t, err)

	clock := clockwork.NewFakeClockAt(time.Now())

	type args struct {
		name    string
		version string
	}
	tests := []struct {
		name    string
		client  *client
		args    args
		wantErr bool
	}{
		{
			name: "can update application binary",
			client: &client{
				ClientParams: &ClientParams{
					Repository: &updateMock.Repository{
						DownloadVersionFunc: func(name, version string) ([]byte, error) {
							req.Equal(binName, name)
							req.Equal("v123.456.789", version)
							clock.Advance(23 * time.Second)
							return []byte("new version"), nil
						},
					},
				},
				clock: clock,
				fs:    &pio.RealFileSystem{},
			},
			args: args{
				name:    binName,
				version: "v123.456.789",
			},
		},
		{
			name: "err if unable to download package",
			client: &client{
				ClientParams: &ClientParams{
					Repository: &updateMock.Repository{
						DownloadVersionFunc: func(name, version string) ([]byte, error) {
							return nil, fmt.Errorf("out of disk")
						},
					},
				},
				clock: clock,
				fs:    &pio.RealFileSystem{},
			},
			args: args{
				name:    binName,
				version: "v1",
			},
			wantErr: true,
		},
		{
			name: "no attempt to mv binary (darwin)",
			client: &client{
				ClientParams: &ClientParams{
					Repository: &updateMock.Repository{
						DownloadVersionFunc: func(name, version string) ([]byte, error) {
							req.Equal(binName, name)
							req.Equal("v1", version)
							clock.Advance(23 * time.Second)
							return []byte("new version"), nil
						},
					},
					OS: "darwin",
				},
				clock: clock,
				fs: &mock.PassThroughFileSystem{
					Mock: &mock.FileSystem{
						MoveFunc: func(src, dst string) error {
							return fmt.Errorf("move func intentionally failed")
						},
					},
					FS: &pio.RealFileSystem{},
				},
			},
			args: args{
				name:    binName,
				version: "v1",
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		if test.client.OS != "" && test.client.OS != runtime.GOOS {
			continue
		}
		t.Run(test.name, func(t *testing.T) {
			test.client.Out = os.Stdout
			if err := test.client.UpdateBinary(test.args.name, test.args.version, true); (err != nil) != test.wantErr {
				t.Errorf("client.UpdateBinary() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestPromptToDownload(t *testing.T) {
	req := require.New(t)

	clock := clockwork.NewFakeClockAt(time.Now())
	countRepeated := 0
	countNoConfirm := 0
	countNoPrompt := 0

	makeFS := func(terminal bool, input string) pio.FileSystem {
		return &mock.PassThroughFileSystem{
			Mock: &mock.FileSystem{
				IsTerminalFunc: func(fd uintptr) bool {
					return terminal
				},
				NewBufferedReaderFunc: func(rd io.Reader) pio.Reader {
					req.Equal(os.Stdin, rd)
					_, _ = fmt.Println() // to go to newline after test prompt
					return bytes.NewBuffer([]byte(input + "\n"))
				},
			},
			FS: &pio.RealFileSystem{},
		}
	}

	makeClient := func(fs pio.FileSystem) *client {
		client := NewClient(&ClientParams{
			Repository: &updateMock.Repository{},
		})
		client.clock = clock
		client.fs = fs
		return client
	}

	type args struct {
		name          string
		currVersion   string
		latestVersion string
		confirm       bool
	}

	basicArgs := args{
		name:          "my-cli",
		currVersion:   "v1.2.0",
		latestVersion: "v2.0.0",
		confirm:       true,
	}

	tests := []struct {
		name   string
		client *client
		args   args
		want   bool
	}{
		{
			name:   "should prompt interactively and return true for yes",
			client: makeClient(makeFS(true, "yes")),
			args:   basicArgs,
			want:   true,
		},
		{
			name:   "should prompt interactively and return true for y",
			client: makeClient(makeFS(true, "y")),
			args:   basicArgs,
			want:   true,
		},
		{
			name:   "should prompt interactively and return true for Y",
			client: makeClient(makeFS(true, "Y")),
			args:   basicArgs,
			want:   true,
		},
		{
			name:   "should prompt interactively and return false for no",
			client: makeClient(makeFS(true, "no")),
			args:   basicArgs,
			want:   false,
		},
		{
			name:   "should prompt interactively and return false for n",
			client: makeClient(makeFS(true, "n")),
			args:   basicArgs,
			want:   false,
		},
		{
			name:   "should prompt interactively and return false for N",
			client: makeClient(makeFS(true, "N")),
			args:   basicArgs,
			want:   false,
		},
		{
			name:   "should prompt interactively and ignore trailing whitespace",
			client: makeClient(makeFS(true, "y ")),
			args:   basicArgs,
			want:   true,
		},
		{
			name: "should prompt repeatedly until user enters yes/no",
			client: makeClient(&mock.PassThroughFileSystem{
				Mock: &mock.FileSystem{
					IsTerminalFunc: func(fd uintptr) bool {
						return true
					},
					NewBufferedReaderFunc: func(rd io.Reader) pio.Reader {
						req.Equal(os.Stdin, rd)
						_, _ = fmt.Println() // to go to newline after test prompt
						countRepeated++
						switch countRepeated {
						case 1:
							return bytes.NewBuffer([]byte("maybe"))
						case 2:
							return bytes.NewBuffer([]byte("youwish"))
						case 3:
							return bytes.NewBuffer([]byte("YES"))
						case 4:
							return bytes.NewBuffer([]byte("never"))
						case 5:
							return bytes.NewBuffer([]byte("no"))
						}
						return bytes.NewBuffer([]byte("n"))
					},
				},
				FS: &pio.RealFileSystem{},
			}),
			args: basicArgs,
			want: false,
		},
		{
			name: "should skip confirmation if not requested",
			client: makeClient(&mock.PassThroughFileSystem{
				Mock: &mock.FileSystem{
					IsTerminalFunc: func(fd uintptr) bool {
						return true
					},
					NewBufferedReaderFunc: func(rd io.Reader) pio.Reader {
						countNoConfirm++
						return bytes.NewBuffer([]byte("n"))
					},
				},
				FS: &pio.RealFileSystem{},
			}),
			args: args{
				name:          "my-cli",
				currVersion:   "v1.2.0",
				latestVersion: "v2.0.0",
				confirm:       false,
			},
			want: true,
		},
		{
			name: "should skip confirmation if not a TTY",
			client: makeClient(&mock.PassThroughFileSystem{
				Mock: &mock.FileSystem{
					IsTerminalFunc: func(fd uintptr) bool {
						return false
					},
					NewBufferedReaderFunc: func(rd io.Reader) pio.Reader {
						countNoPrompt++
						return bytes.NewBuffer([]byte("n"))
					},
				},
				FS: &pio.RealFileSystem{},
			}),
			args: args{
				name:          "my-cli",
				currVersion:   "v1.2.0",
				latestVersion: "v2.0.0",
				confirm:       false,
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.client.Out == nil {
				test.client.Out = os.Stdout
			}
			if got := test.client.PromptToDownload(test.args.name, test.args.currVersion, test.args.latestVersion, "", test.args.confirm); got != test.want {
				t.Errorf("client.PromptToDownload() = %v, want %v", got, test.want)
			}
		})
	}
	req.Equal(5, countRepeated)
	req.Equal(0, countNoConfirm)
	req.Equal(0, countNoPrompt)
}

func TestGetBinaryName(t *testing.T) {
	assert.Equal(t, "confluent_3.13.0_darwin_amd64", getBinaryName("3.13.0", "darwin", "amd64"))
	assert.Equal(t, "confluent_3.13.0_windows_amd64.exe", getBinaryName("3.13.0", "windows", "amd64"))
}

func TestFindChecksum(t *testing.T) {
	content := strings.Join([]string{
		"0e3b559127d31a3f4bd9833e31ddd60d74efbd52d088e7a8b81ea402c4b80c37  confluent_3.13.0_linux_amd64",
		"495bfcb16f1b33a37a6c0d3941ea4b82756ee5d3329f9cc223269daeadd08e7c  confluent_3.13.0_darwin_amd64",
		"cf1f7f14c5bc31e502f8b75f98fa6caff02617261318810ed93fed358e28f994  confluent_3.13.0_linux_amd64.tar.gz",
		"e0e3377b2297060bfe6cf918cd926ff0e240d4115bd314bd9ac53c0f5c47ebcd  confluent_3.13.0_darwin_amd64.tar.gz",
	}, "\n")

	checksum, err := findChecksum(content, "confluent_3.13.0_darwin_amd64")
	require.NoError(t, err)
	require.Equal(t, "495bfcb16f1b33a37a6c0d3941ea4b82756ee5d3329f9cc223269daeadd08e7c", checksum)
}
