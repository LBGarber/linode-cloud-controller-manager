package framework

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	scriptDirectory = "scripts"
	RetryInterval   = 5 * time.Second
	RetryTimout     = 15 * time.Minute
)

func RunScript(script string, args ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return runCommand(path.Join(wd, scriptDirectory, script), args...)
}

func runCommand(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = append(c.Env, append(os.Environ())...)
	glog.Infof("Running command %q\n", cmd)
	return c.Run()
}

func deleteInForeground() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationForeground
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}

func ApplyManifest(commandName, manifest string) error {
	args := []string{commandName, "--kubeconfig", KubeConfigFile, "-f", "-"}
	cmd := exec.Command("kubectl", args...)
	cmd.Stdin = strings.NewReader(manifest)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}
