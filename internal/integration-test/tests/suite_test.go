package tests_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"integration-tests/container"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestCLISuite(t *testing.T) {
	suite.Run(t, &cliSuite{})
}

type benchRun struct {
	Name        string  `json:"name"`
	BytesPerOp  string  `json:"bop"`
	AllocsPerOp string  `json:"allocsop"`
	Score       float64 `json:"score"`
	NsPerOp     float64 `json:"nsop"`
}

type benchResult struct {
	Os      string     `json:"os"`
	Arch    string     `json:"arch"`
	Runs    []benchRun `json:"runs"`
	RunTime float64    `json:"runTime"`
}

type cliSuite struct {
	suite.Suite
	cli         *container.CLIContainer
	jsonRaw     []byte
	csvRaw      []byte
	jsonResults []benchResult
	sample      *benchRun
}

func (s *cliSuite) SetupSuite() {
	s.cli = container.NewCLIContainer(context.Background())
	require.NoError(s.T(), s.cli.Run(), "start go-run-bench container")

	require.NoError(s.T(), s.cli.CleanWorkDir())
	stdout, _, code, err := s.cli.RunCLI(
		"-cooldown=0", "-benchtime=1", "-benchmem=true", "-save=json",
		"-wd="+container.FixtureDir,
	)
	require.NoError(s.T(), err)
	require.Equalf(s.T(), 0, code, "canonical json run failed: %s", stdout)
	jsonName, err := s.cli.FindResultFile("json")
	require.NoError(s.T(), err)
	s.jsonRaw, err = s.cli.ReadWorkFile(jsonName)
	require.NoError(s.T(), err)
	require.NoError(s.T(), json.Unmarshal(s.jsonRaw, &s.jsonResults))
	require.NotEmpty(s.T(), s.jsonResults)
	s.sample = findSampleRun(s.jsonResults)
	require.NotNil(s.T(), s.sample, "BenchmarkSample missing in canonical JSON")

	require.NoError(s.T(), s.cli.CleanWorkDir())
	stdout, _, code, err = s.cli.RunCLI(
		"-cooldown=0", "-benchtime=1", "-benchmem=true", "-save=csv",
		"-wd="+container.FixtureDir,
	)
	require.NoError(s.T(), err)
	require.Equalf(s.T(), 0, code, "canonical csv run failed: %s", stdout)
	csvName, err := s.cli.FindResultFile("csv")
	require.NoError(s.T(), err)
	s.csvRaw, err = s.cli.ReadWorkFile(csvName)
	require.NoError(s.T(), err)
}

func (s *cliSuite) TearDownSuite() {
	if s.cli != nil {
		_ = s.cli.Terminate()
	}
}

func findSampleRun(results []benchResult) *benchRun {
	for i := range results {
		for j := range results[i].Runs {
			if strings.HasPrefix(results[i].Runs[j].Name, "BenchmarkSample") {
				return &results[i].Runs[j]
			}
		}
	}
	return nil
}

func (s *cliSuite) runCLI(args ...string) (string, int) {
	stdout, _, code, err := s.cli.RunCLI(args...)
	require.NoError(s.T(), err)
	return stdout, code
}
