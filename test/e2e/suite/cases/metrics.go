package cases

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cert-manager/csi-driver/test/e2e/framework"
)

const (
	certManagerNamespace = "cert-manager"
)

var _ = framework.CasesDescribe("Metrics", func() {
	f := framework.NewDefaultFramework("metrics")

	FIt("Should serve Go and process metrics on port 9402", func() {
		By("Discovering the name of the csi-driver Pod")
		pods, err := f.KubeClientSet.CoreV1().Pods(certManagerNamespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/instance=csi-driver",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(pods.Items).To(HaveLen(1))
		p := pods.Items[0]
		By("Connecting to Pod on default metrics port 9402 and sending a GET request to the /metrics endpoint")
		respBytes, err := f.KubeClientSet.
			CoreV1().
			Pods(p.Namespace).
			ProxyGet("http", p.Name, "9402", "/metrics", map[string]string{}).
			DoRaw(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		resp := string(respBytes)
		Expect(resp).To(ContainSubstring("# HELP go_threads Number of OS threads created."),
			"go_collector metrics should be available")
		Expect(resp).To(ContainSubstring("# HELP process_open_fds Number of open file descriptors."),
			"process_collector metrics should be available")
	})
})
