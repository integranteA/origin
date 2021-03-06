package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"

	"github.com/openshift/origin/test/extended/util"
)

var _ = g.Describe("[cli] oc adm must-gather", func() {
	defer g.GinkgoRecover()
	oc := util.NewCLI("oc-adm-must-gather", util.KubeConfigPath()).AsAdmin()
	g.It("runs successfully", func() {
		tempDir, err := ioutil.TempDir("", "test.oc-adm-must-gather.")
		o.Expect(err).ToNot(o.HaveOccurred())
		defer os.RemoveAll(tempDir)
		o.Expect(oc.Run("adm", "must-gather").Args("--dest-dir", tempDir).Execute()).To(o.Succeed())

		expectedDirectories := [][]string{
			{tempDir, "cluster-scoped-resources", "config.openshift.io"},
			{tempDir, "cluster-scoped-resources", "operator.openshift.io"},
			{tempDir, "cluster-scoped-resources", "core"},
			{tempDir, "cluster-scoped-resources", "apiregistration.k8s.io"},
			{tempDir, "namespaces", "openshift"},
			{tempDir, "namespaces", "openshift-kube-apiserver-operator"},
		}

		expectedFiles := [][]string{
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "apiservers.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "authentications.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "builds.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "clusteroperators.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "clusterversions.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "consoles.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "dnses.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "featuregates.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "images.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "infrastructures.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "ingresses.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "networks.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "oauths.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "projects.yaml"},
			{tempDir, "cluster-scoped-resources", "config.openshift.io", "schedulers.yaml"},
			{tempDir, "namespaces", "openshift-kube-apiserver", "core", "configmaps.yaml"},
			{tempDir, "namespaces", "openshift-kube-apiserver", "core", "secrets.yaml"},
			{tempDir, "audit_logs", "kube-apiserver.audit_logs_listing"},
			{tempDir, "audit_logs", "openshift-apiserver.audit_logs_listing"},
			{tempDir, "host_service_logs", "masters", "crio_service.log"},
			{tempDir, "host_service_logs", "masters", "kubelet_service.log"},
		}

		for _, expectedDirectory := range expectedDirectories {
			o.Expect(path.Join(expectedDirectory...)).To(o.BeADirectory())
		}

		emptyFiles := []string{}
		for _, expectedFile := range expectedFiles {
			expectedFilePath := path.Join(expectedFile...)
			o.Expect(expectedFilePath).To(o.BeAnExistingFile())
			stat, err := os.Stat(expectedFilePath)
			o.Expect(err).ToNot(o.HaveOccurred())
			if size := stat.Size(); size < 100 {
				emptyFiles = append(emptyFiles, expectedFilePath)
			}
		}
		if len(emptyFiles) > 0 {
			o.Expect(fmt.Errorf("expected files should not be empty: %s", strings.Join(emptyFiles, ","))).NotTo(o.HaveOccurred())
		}

	})

	g.It("runs successfully with options", func() {
		tempDir, err := ioutil.TempDir("", "test.oc-adm-must-gather.")
		o.Expect(err).ToNot(o.HaveOccurred())
		defer os.RemoveAll(tempDir)
		args := []string{
			"--dest-dir", tempDir,
			"--source-dir", "/artifacts",
			"--",
			"/bin/bash", "-c",
			"ls -l > /artifacts/ls.log",
		}
		o.Expect(oc.Run("adm", "must-gather").Args(args...).Execute()).To(o.Succeed())
		expectedFilePath := path.Join(tempDir, "ls.log")
		o.Expect(expectedFilePath).To(o.BeAnExistingFile())
		stat, err := os.Stat(expectedFilePath)
		o.Expect(err).ToNot(o.HaveOccurred())
		o.Expect(stat.Size()).To(o.BeNumerically(">", 0))
	})
})
