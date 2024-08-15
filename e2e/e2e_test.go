package e2e_test

import (
	"bytes"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var kubeConfigPath string

var _ = BeforeSuite(func() {
	kubeConfigPath = setupKindCluster()
	buildImage()
	importImage()
})

var _ = AfterSuite(func() {
	teardownKindCluster(kubeConfigPath)
})

var _ = Describe("E2E Test", func() {
	Context("When deploying the application", func() {
		It("should deploy successfully and produce logs", func() {
			deployApplication(kubeConfigPath)
			// takes 30 seconds for the pods to start being created
			time.Sleep(30 * time.Second)

			// Check logs
			podLogs := getPodLogs(kubeConfigPath)
			Expect(podLogs).To(ContainSubstring("Auditor successfully created pod:"))
			Expect(podLogs).To(ContainSubstring("Watch Controller successfully deleted pod:"))
		})
	})
})

func buildImage() {
	cmd := exec.Command("docker", "build", "-t", "watch-auditor:dev", ".", "-f", "../Dockerfile.arm")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		GinkgoWriter.Write(out.Bytes())
	}
}

func importImage() {
	cmd := exec.Command("kind", "load", "docker-image", "watch-auditor:dev")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		GinkgoWriter.Write(out.Bytes())
		Expect(err).NotTo(HaveOccurred(), "Failed to import the Docker image into kind")
	}
}

func setupKindCluster() string {
	cmd := exec.Command("kind", "create", "cluster")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		GinkgoWriter.Write(out.Bytes())
		Expect(err).NotTo(HaveOccurred(), "Failed to create the kind cluster")
	}

	kubeConfigPath := out.String()
	return kubeConfigPath
}

func teardownKindCluster(kubeConfigPath string) {
	// Delete the kind cluster
	cmd := exec.Command("kind", "delete", "cluster", "--kubeconfig", kubeConfigPath)
	err := cmd.Run()
	Expect(err).NotTo(HaveOccurred())
}

func deployApplication(kubeConfigPath string) {
	// Deploy your application using kubectl or any other method
	cmd := exec.Command("kubectl", "apply", "-k", "../kustomize/overlays/dev", "--kubeconfig", kubeConfigPath)
	err := cmd.Run()
	Expect(err).NotTo(HaveOccurred())
}

func getPodLogs(kubeConfigPath string) string {
	// Get logs from the deployed pod
	cmd := exec.Command("kubectl", "logs", "deploy/watch-auditor", "-n", "watch-auditor", "--kubeconfig", kubeConfigPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	return out.String()
}
