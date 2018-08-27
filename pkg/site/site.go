package site

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	rcopy "github.com/otiai10/copy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/hugo/commands"
)

//go:generate rice embed-go

// Build builds the site HTML.
func Build(dir, outputDir string) error {
	templateBox, err := rice.FindBox("site_template")
	if err != nil {
		return errors.Wrap(err, "find embedded site template")
	}

	logrus.Info("staging embedded template")
	if err = templateBox.Walk("/", walkFn(templateBox, dir)); err != nil {
		return errors.Wrap(err, "walking embedded template")
	}

	logrus.Info("installing npm dependencies")
	if err = npmInstall(dir); err != nil {
		return errors.Wrap(err, "npm install")
	}

	logrus.Info("building hugo site")
	if err = hugoBuild(dir); err != nil {
		return errors.Wrap(err, "hugo build")
	}

	logrus.Infof("create output directory: %s", outputDir)
	if err = createOutput(outputDir, dir); err != nil {
		return errors.Wrap(err, "create output directory")
	}

	return nil
}

func walkFn(box *rice.Box, dir string) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dest := filepath.Join(dir, path)
		if fi.IsDir() {
			return os.MkdirAll(dest, 0700)
		}

		data, err := box.Bytes(path)
		if err != nil {
			return err
		}

		return copyData(data, dest)
	}
}

func copyData(data []byte, dest string) error {
	from := bytes.NewReader(data)

	to, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrapf(err, "create %s", dest)
	}

	_, err = io.Copy(to, from)
	if err != nil {
		return errors.Wrapf(err, "copying source to destination")
	}

	return nil
}

func npmInstall(dir string) error {
	cmd := exec.Command("npm", "i")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func hugoBuild(dir string) error {
	os.Setenv("HUGO_PARAMS", "{\"version\": \"(TBD)\"}")

	args := []string{"-s", dir}
	resp := commands.Execute(args)
	return resp.Err
}

func createOutput(outputDir, dir string) error {
	publicDir := filepath.Join(dir, "public")

	return rcopy.Copy(publicDir, outputDir)
}
