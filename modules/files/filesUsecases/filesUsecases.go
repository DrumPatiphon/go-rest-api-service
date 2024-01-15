package filesUsecases

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/files"
)

type IFilesUsecases interface {
	UploadToGCP(req []*files.FileReq) ([]*files.FileRes, error)
	DeleteFileOnGCP(req []*files.DeleteFileReq) error
}

type filesUsecase struct {
	cfg config.Iconfig
}

func FilesUsecase(cfg config.Iconfig) IFilesUsecases {
	return &filesUsecase{
		cfg: cfg,
	}
}

type filesPublic struct {
	bucket      string
	destination string
	file        *files.FileRes
}

// makePublic gives all users read access to an object.
func (f *filesPublic) makePublic(ctx context.Context, client *storage.Client) error {
	acl := client.Bucket(f.bucket).Object(f.destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %v", err)
	}
	fmt.Printf("Blob %v is now publicly accessible.\n", f.destination)
	return nil
}

func (u *filesUsecase) uploadWorkers(ctx context.Context, client *storage.Client, jobsCh <-chan *files.FileReq, resultsCh chan<- *files.FileRes, errsCh chan<- error) {
	for job := range jobsCh {
		cotainer, err := job.File.Open()
		if err != nil {
			errsCh <- err
			return
		}
		b, err := ioutil.ReadAll(cotainer)
		if err != nil {
			errsCh <- err
			return
		}

		buf := bytes.NewBuffer(b)

		// Upload an object with storage.Writer.
		wc := client.Bucket(u.cfg.App().Gcpbucket()).Object(job.Destination).NewWriter(ctx)

		if _, err = io.Copy(wc, buf); err != nil {
			errsCh <- fmt.Errorf("io.Copy: %v", err)
			return
		}
		// Data can continue to be added to the file until the writer is closed.
		if err := wc.Close(); err != nil {
			errsCh <- fmt.Errorf("Writer.Close: %v", err)
			return
		}
		fmt.Printf("%v uploaded to %v.\n", job.FileName, job.Extension)

		newFile := &filesPublic{
			file: &files.FileRes{
				FileName: job.FileName,
				Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.cfg.App().Gcpbucket(), job.Destination),
			},
			bucket:      u.cfg.App().Gcpbucket(),
			destination: job.Destination,
		}

		if err := newFile.makePublic(ctx, client); err != nil {
			errsCh <- err
			return
		}

		errsCh <- nil
		resultsCh <- newFile.file
	}

}

// pull worker
func (u *filesUsecase) UploadToGCP(req []*files.FileReq) ([]*files.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	jobsCh := make(chan *files.FileReq, len(req))
	resultsCh := make(chan *files.FileRes, len(req))
	errsCh := make(chan error, len(req))

	res := make([]*files.FileRes, 0)

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.uploadWorkers(ctx, client, jobsCh, resultsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		if err != nil {
			return nil, err
		}

		result := <-resultsCh
		res = append(res, result)
	}

	return res, nil
}

// deleteFile removes specified object.
func (u *filesUsecase) deleteFileWorkers(ctx context.Context, jobsCh <-chan *files.DeleteFileReq, client *storage.Client, errsCh chan<- error) {

	for job := range jobsCh {

		o := client.Bucket(u.cfg.App().Gcpbucket()).Object(job.Destination)

		// Optional: set a generation-match precondition to avoid potential race
		// conditions and data corruptions. The request to delete the file is aborted
		// if the object's generation number does not match your precondition.
		attrs, err := o.Attrs(ctx)
		if err != nil {
			errsCh <- fmt.Errorf("object.Attrs: %v", err)
			return
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			errsCh <- fmt.Errorf("Object(%q).Delete: %v", job.Destination, err)
			return
		}
		fmt.Printf("Blob %v deleted.\n", job.Destination)

		errsCh <- nil
	}
}

func (u *filesUsecase) DeleteFileOnGCP(req []*files.DeleteFileReq) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	jobsCh := make(chan *files.DeleteFileReq, len(req))
	errsCh := make(chan error, len(req))

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.deleteFileWorkers(ctx, jobsCh, client, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		return err
	}

	return nil
}
