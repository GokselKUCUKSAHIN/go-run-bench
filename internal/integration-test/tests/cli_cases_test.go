package tests_test

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"

	"integration-tests/container"

	"github.com/stretchr/testify/require"
)

// TestCLICases expands to exactly 50 integration subtests.
func (s *cliSuite) TestCLICases() {
	type tc struct {
		name string
		fn   func()
	}

	cases := []tc{
		// 1-3 help
		{"help_bare", func() {
			out, code := s.runCLI("help")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Usage:")
		}},
		{"help_dash", func() {
			out, code := s.runCLI("-help")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Usage:")
		}},
		{"help_double_dash", func() {
			out, code := s.runCLI("--help")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Usage:")
		}},
		// 4-6 help content
		{"help_mentions_cooldown", func() {
			out, _ := s.runCLI("-help")
			require.Contains(s.T(), out, "-cooldown")
		}},
		{"help_mentions_benchtime", func() {
			out, _ := s.runCLI("-help")
			require.Contains(s.T(), out, "-benchtime")
		}},
		{"help_mentions_benchmem", func() {
			out, _ := s.runCLI("-help")
			require.Contains(s.T(), out, "-benchmem")
		}},
		{"help_mentions_save", func() {
			out, _ := s.runCLI("-help")
			require.Contains(s.T(), out, "-save")
		}},
		{"help_mentions_wd", func() {
			out, _ := s.runCLI("-help")
			require.Contains(s.T(), out, "-wd")
		}},
		{"help_mentions_windows_path", func() {
			out, _ := s.runCLI("-help")
			require.Contains(s.T(), out, `C:\`)
		}},
		// 10 binary / env
		{"binary_on_path", func() {
			stdout, _, code, err := s.cli.Exec("which", "go-run-bench")
			require.NoError(s.T(), err)
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), stdout, "go-run-bench")
		}},
		{"go_toolchain_present", func() {
			stdout, _, code, err := s.cli.Exec("go", "version")
			require.NoError(s.T(), err)
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), stdout, "go version")
		}},
		{"fixture_dir_exists", func() {
			_, _, code, err := s.cli.Exec("test", "-d", container.FixtureDir)
			require.NoError(s.T(), err)
			require.Equal(s.T(), 0, code)
		}},
		{"fixture_benchmark_file_exists", func() {
			_, _, code, err := s.cli.Exec("test", "-f", path.Join(container.FixtureDir, "sample_benchmark_test.go"))
			require.NoError(s.T(), err)
			require.Equal(s.T(), 0, code)
		}},
		{"fixture_go_mod_exists", func() {
			_, _, code, err := s.cli.Exec("test", "-f", path.Join(container.FixtureDir, "go.mod"))
			require.NoError(s.T(), err)
			require.Equal(s.T(), 0, code)
		}},
		{"work_dir_exists", func() {
			_, _, code, err := s.cli.Exec("test", "-d", container.WorkDir)
			require.NoError(s.T(), err)
			require.Equal(s.T(), 0, code)
		}},
		// 16-30 invalid flags (parse errors print help and exit 0)
		{"invalid_arg_no_hyphen", func() {
			out, code := s.runCLI("cooldown=0")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
			require.Contains(s.T(), out, "Usage:")
		}},
		{"invalid_arg_no_equals", func() {
			out, code := s.runCLI("-cooldown")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_unknown_flag", func() {
			out, code := s.runCLI("-foo=bar")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "unknown parameter")
		}},
		{"invalid_cooldown_text", func() {
			out, code := s.runCLI("-cooldown=abc")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_cooldown_negative", func() {
			out, code := s.runCLI("-cooldown=-1")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_cooldown_too_large", func() {
			out, code := s.runCLI("-cooldown=301")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_benchtime_zero", func() {
			out, code := s.runCLI("-benchtime=0")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_benchtime_negative", func() {
			out, code := s.runCLI("-benchtime=-5")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_benchtime_too_large", func() {
			out, code := s.runCLI("-benchtime=31")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_benchtime_text", func() {
			out, code := s.runCLI("-benchtime=xx")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_benchmem_value", func() {
			out, code := s.runCLI("-benchmem=maybe")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_save_format", func() {
			out, code := s.runCLI("-save=xml")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_wd_relative", func() {
			out, code := s.runCLI("-wd=relative/path")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_wd_missing", func() {
			out, code := s.runCLI("-wd=/this/path/does/not/exist/go-run-bench")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		{"invalid_wd_is_file", func() {
			out, code := s.runCLI("-wd=/usr/local/bin/go-run-bench")
			require.Equal(s.T(), 0, code)
			require.Contains(s.T(), out, "Error:")
		}},
		// 31-40 JSON cached asserts
		{"json_is_valid_array", func() {
			var results []benchResult
			require.NoError(s.T(), json.Unmarshal(s.jsonRaw, &results))
			require.NotEmpty(s.T(), results)
		}},
		{"json_os_present", func() {
			require.NotEmpty(s.T(), s.jsonResults[0].Os)
		}},
		{"json_arch_present", func() {
			require.NotEmpty(s.T(), s.jsonResults[0].Arch)
		}},
		{"json_os_is_linux", func() {
			require.Equal(s.T(), "linux", s.jsonResults[0].Os)
		}},
		{"json_runtime_non_negative", func() {
			require.GreaterOrEqual(s.T(), s.jsonResults[0].RunTime, 0.0)
		}},
		{"json_has_runs", func() {
			require.NotEmpty(s.T(), s.jsonResults[0].Runs)
		}},
		{"json_sample_name", func() {
			require.True(s.T(), strings.HasPrefix(s.sample.Name, "BenchmarkSample"))
		}},
		{"json_sample_nsop_positive", func() {
			require.Greater(s.T(), s.sample.NsPerOp, 0.0)
		}},
		{"json_sample_score_positive", func() {
			require.Greater(s.T(), s.sample.Score, 0.0)
		}},
		{"json_sample_bop_present", func() {
			require.NotEmpty(s.T(), s.sample.BytesPerOp)
		}},
		{"json_sample_allocsop_present", func() {
			require.NotEmpty(s.T(), s.sample.AllocsPerOp)
		}},
		{"json_raw_starts_with_bracket", func() {
			require.True(s.T(), strings.HasPrefix(strings.TrimSpace(string(s.jsonRaw)), "["))
		}},
		// 43-48 CSV cached asserts
		{"csv_has_header", func() {
			lines := strings.Split(strings.TrimSpace(string(s.csvRaw)), "\n")
			require.GreaterOrEqual(s.T(), len(lines), 2)
			require.Equal(s.T(), "Name;Score;ns/op;B/op;allocs/op", strings.TrimSpace(lines[0]))
		}},
		{"csv_has_sample_row", func() {
			require.Contains(s.T(), string(s.csvRaw), "BenchmarkSample")
		}},
		{"csv_uses_semicolon", func() {
			require.Contains(s.T(), string(s.csvRaw), ";")
		}},
		{"csv_sample_row_field_count", func() {
			lines := strings.Split(strings.TrimSpace(string(s.csvRaw)), "\n")
			var row string
			for _, line := range lines[1:] {
				if strings.Contains(line, "BenchmarkSample") {
					row = line
					break
				}
			}
			require.NotEmpty(s.T(), row)
			require.Equal(s.T(), 5, len(strings.Split(row, ";")))
		}},
		{"csv_sample_score_parseable", func() {
			lines := strings.Split(strings.TrimSpace(string(s.csvRaw)), "\n")
			var row string
			for _, line := range lines[1:] {
				if strings.Contains(line, "BenchmarkSample") {
					row = line
					break
				}
			}
			parts := strings.Split(row, ";")
			require.Len(s.T(), parts, 5)
			score, err := strconv.ParseFloat(parts[1], 64)
			require.NoError(s.T(), err)
			require.Greater(s.T(), score, 0.0)
		}},
		{"csv_sample_nsop_parseable", func() {
			lines := strings.Split(strings.TrimSpace(string(s.csvRaw)), "\n")
			var row string
			for _, line := range lines[1:] {
				if strings.Contains(line, "BenchmarkSample") {
					row = line
					break
				}
			}
			parts := strings.Split(row, ";")
			nsop, err := strconv.ParseFloat(parts[2], 64)
			require.NoError(s.T(), err)
			require.Greater(s.T(), nsop, 0.0)
		}},
		// 49-50 live runs + empty dir
		{"live_json_workingdirectory_alias", func() {
			require.NoError(s.T(), s.cli.CleanWorkDir())
			out, code := s.runCLI(
				"-cooldown=0", "-benchtime=1", "-save=json",
				"-workingdirectory="+container.FixtureDir,
			)
			require.Equalf(s.T(), 0, code, out)
			name, err := s.cli.FindResultFile("json")
			require.NoError(s.T(), err)
			raw, err := s.cli.ReadWorkFile(name)
			require.NoError(s.T(), err)
			var results []benchResult
			require.NoError(s.T(), json.Unmarshal(raw, &results))
			require.NotNil(s.T(), findSampleRun(results))
		}},
		{"live_csv_benchmem_false_empty_dir", func() {
			require.NoError(s.T(), s.cli.CleanWorkDir())
			_, _, code, err := s.cli.Exec("mkdir", "-p", "/empty-fixture")
			require.NoError(s.T(), err)
			require.Equal(s.T(), 0, code)
			out, code := s.runCLI(
				"-cooldown=0", "-benchtime=1", "-benchmem=false", "-save=csv",
				"-wd=/empty-fixture",
			)
			require.Equalf(s.T(), 0, code, out)
			name, err := s.cli.FindResultFile("csv")
			require.NoError(s.T(), err)
			raw, err := s.cli.ReadWorkFile(name)
			require.NoError(s.T(), err)
			lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
			require.Equal(s.T(), "Name;Score;ns/op;B/op;allocs/op", strings.TrimSpace(lines[0]))
			require.Len(s.T(), lines, 1) // header only
		}},
	}

	require.Len(s.T(), cases, 50, "integration suite must expose exactly 50 cases, got %d", len(cases))

	for i, c := range cases {
		name := fmt.Sprintf("%02d_%s", i+1, c.name)
		s.Run(name, c.fn)
	}
}
