package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/spf13/cobra"
)

// Register sub commands
func init() {
	cmd := getCmdS3Upload()
	addFlagsS3Upload(cmd)

	CmdS3.AddCommand(cmd)
}

func addFlagsS3Upload(cmd *cobra.Command) {
	cmd.Flags().String(CMD_S3_UPLOAD_BUCKET, "", "S3 bucket name")
	cmd.Flags().String(CMD_S3_UPLOAD_PREFIX, "", "The path prefix for S3 bucket that the objects will be uploaded to")
	cmd.Flags().BoolP(CMD_S3_UPLOAD_RECURSIVE, "r", false, "Recursively travel the given directory for all objects")
	cmd.Flags().String(CMD_S3_UPLOAD_EXCLUDE_FILES, "", "Exclude files with matching file names from upload. Multiple file names seperate by comma")
}

// cmd: upload
func getCmdS3Upload() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "upload objects to s3 bucket",
		Long:  `upload objects to s3 bucket`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing local objects path")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			recursive, _ := cmd.Flags().GetBool(CMD_S3_UPLOAD_RECURSIVE)

			for _, arg := range args {
				err := s3Upload(
					arg,
					cmd.Flags().Lookup(CMD_S3_UPLOAD_BUCKET).Value.String(),
					cmd.Flags().Lookup(CMD_S3_UPLOAD_PREFIX).Value.String(),
					recursive,
					cmd.Flags().Lookup(CMD_S3_UPLOAD_EXCLUDE_FILES).Value.String(),
				)

				silenceUsageOnError(cmd, err)

				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}

// Upload
func s3Upload(objPath, bucket, prefix string, recursive bool, exf string) error {
	// Default only current dir
	level := 1
	if recursive {
		level = 0
	}

	var exfiles []string
	if len(exf) > 0 {
		exfiles = strings.Split(exf, ",")
	}

	done := make(chan bool)
	defer close(done)

	isDir, err := utils.IsDir(objPath)
	if err != nil {
		return err
	}

	fmt.Println("")

	// If it's only one file
	if !isDir {
		content, err := ioutil.ReadFile(objPath)
		if err != nil {
			return err
		}

		// Upload all nested template to s3
		cfs3 := ctlaws.NewS3(s3.New(ctlaws.AWSSess))
		out, err := cfs3.Upload(bucket, path.Join(prefix, objPath), content)
		if err != nil {
			return err
		}

		s3UploadPrintToConsole(objPath, out.Location)
	} else {
		paths, errc := utils.ScanFiles(objPath, done, level)

		// Start 10 workers
		var wg sync.WaitGroup
		numProc := 10
		wg.Add(numProc)

		startPath, _ := filepath.Abs(objPath)
		startPath = path.Base(startPath)
		result := make(chan *uploadResult)
		for i := 0; i < numProc; i++ {
			go func() {
				uploadWorker(bucket, prefix, startPath, paths, result, done, exfiles)
				wg.Done()
			}()
		}

		// Close result when all workers
		go func() {
			wg.Wait()
			close(result)
		}()

		for r := range result {
			if r.err != nil {
				return r.err
			}

			s3UploadPrintToConsole(r.path, r.output.Location)
		}

		// Check whether the file scan failed.
		if err := <-errc; err != nil {
			return err
		}
	}

	return nil
}

// Output from upload worker
type uploadResult struct {
	path   string
	output *s3manager.UploadOutput
	err    error
}

// Worker to upload object to s3 bucket
func uploadWorker(bucket, prefix, startPath string, paths <-chan string, result chan<- *uploadResult, done <-chan bool, exfiles []string) {
	cfs3 := ctlaws.NewS3(s3.New(ctlaws.AWSSess))

	for p := range paths {
		// Skip file that's in exclusion list
		if utils.InSlice(exfiles, filepath.Base(p)) {
			continue
		}

		content, err := ioutil.ReadFile(p)
		if err != nil {
			result <- &uploadResult{err: err}
			continue
		}

		// Upload all nested template to s3
		out, err := cfs3.Upload(bucket, path.Join(prefix, utils.RewritePath(p, startPath)), content)
		if err != nil {
			result <- &uploadResult{err: err}
			continue
		}

		select {
		case result <- &uploadResult{path: p, output: out}:
		case <-done:
			return
		}
	}
}

// Print to console
func s3UploadPrintToConsole(path, s3url string) {
	utils.InfoPrint(
		fmt.Sprintf(
			"[ s3 | upload ] %s -> %s",
			path,
			s3url,
		),
	)
}
