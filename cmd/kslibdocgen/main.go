package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ksonnet/kslibdocs/pkg/docgen"
	"github.com/ksonnet/kslibdocs/pkg/site"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	var groups arrayFlags
	flag.Var(&groups, "groups", "Groups to render. If blank, it will render all groups")

	var k8sLib string
	flag.StringVar(&k8sLib, "k8sLib", "", "Path to k8s.libsonnet")

	var outDir string
	flag.StringVar(&outDir, "outDir", "", "Output directory")

	var templateDir string
	flag.StringVar(&templateDir, "templateDir", "", "Template directory or leave blank to use embedded")

	flag.Parse()

	if err := run(groups, k8sLib, outDir, templateDir); err != nil {
		logrus.WithError(err).Fatal("create ksonnet lib docs")
	}
}

func run(groups []string, k8sLib, outDir, templateDir string) error {
	if k8sLib == "" {
		return errors.New("k8sLib was blank")
	}

	if outDir == "" {
		return errors.New("outDir was blank")
	}

	if templateDir == "" {
		logrus.Info("staging template directory")
		var err error
		templateDir, err = ioutil.TempDir("", "")
		if err != nil {
			return errors.Wrap(err, "creating temp directory")
		}
		defer os.RemoveAll(templateDir)
	}

	if err := generateContent(groups, k8sLib, templateDir); err != nil {
		return errors.Wrap(err, "generate content")
	}

	if err := buildSite(templateDir, outDir); err != nil {
		return errors.Wrap(err, "build site")
	}

	return nil
}

func generateContent(groups []string, k8sLib, outPath string) error {
	logrus.Info("generate content")
	return docgen.Generate(k8sLib, outPath, []string(groups)...)
}

func buildSite(dir, outPath string) error {
	logrus.Info("build site")
	return site.Build(dir, outPath)
}

type arrayFlags []string

func (f *arrayFlags) String() string {
	return strings.Join(*f, ",")
}

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}
