package deployments

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/fuseml/fuseml/cli/helpers"
	"github.com/fuseml/fuseml/cli/kubernetes"
	"github.com/fuseml/fuseml/cli/paas/ui"
	"github.com/kyokomi/emoji"
	"github.com/pkg/errors"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Tekton struct {
	Debug      bool
	Secrets    []string
	ConfigMaps []string
	Timeout    int
}

const (
	TektonDeploymentID      = "tekton"
	tektonNamespace         = "tekton-pipelines"
	tektonPipelineYamlPath  = "tekton/pipeline-v0.22.0.yaml"
	tektonTriggersYamlPath  = "tekton/triggers-v0.12.1.yaml"
	tektonDashboardYamlPath = "tekton/dashboard-v0.15.0.yaml"
	tektonAdminRoleYamlPath = "tekton/admin-role.yaml"
	tektonFusemlYamlPath    = "tekton/fuseml.yaml"
	tektonKanikoYamlPath    = "tekton/kaniko.yaml"
)

func (k *Tekton) ID() string {
	return TektonDeploymentID
}

func (k *Tekton) Backup(c *kubernetes.Cluster, ui *ui.UI, d string) error {
	return nil
}

func (k *Tekton) Restore(c *kubernetes.Cluster, ui *ui.UI, d string) error {
	return nil
}

func (k Tekton) Describe() string {
	return emoji.Sprintf(":cloud:Tekton pipeline: %s\n:cloud:Tekton dashboard: %s\n:cloud:Tekton triggers: %s\n",
		tektonPipelineYamlPath, tektonDashboardYamlPath, tektonTriggersYamlPath)
}

// Delete removes Tekton from kubernetes cluster
func (k Tekton) Delete(c *kubernetes.Cluster, ui *ui.UI) error {
	ui.Note().KeeplineUnder(1).Msg("Removing Tekton...")

	existsAndOwned, err := c.NamespaceExistsAndOwned(tektonNamespace)
	if err != nil {
		return errors.Wrapf(err, "failed to check if namespace '%s' is owned or not", tektonNamespace)
	}
	if !existsAndOwned {
		ui.Exclamation().Msg("Skipping Tekton because namespace either doesn't exist or not owned by Fuseml")
		return nil
	}

	if out, err := helpers.KubectlDeleteEmbeddedYaml(tektonDashboardYamlPath, true); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Deleting %s failed:\n%s", tektonDashboardYamlPath, out))
	}
	if out, err := helpers.KubectlDeleteEmbeddedYaml(tektonAdminRoleYamlPath, true); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Deleting %s failed:\n%s", tektonAdminRoleYamlPath, out))
	}
	if out, err := helpers.KubectlDeleteEmbeddedYaml(tektonTriggersYamlPath, true); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Deleting %s failed:\n%s", tektonTriggersYamlPath, out))
	}
	if out, err := helpers.KubectlDeleteEmbeddedYaml(tektonPipelineYamlPath, true); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Deleting %s failed:\n%s", tektonPipelineYamlPath, out))
	}

	message := "Deleting Tekton namespace " + tektonNamespace
	_, err = helpers.WaitForCommandCompletion(ui, message,
		func() (string, error) {
			return "", c.DeleteNamespace(tektonNamespace)
		},
	)
	if err != nil {
		return errors.Wrapf(err, "Failed deleting namespace %s", tektonNamespace)
	}

	ui.Success().Msg("Tekton removed")

	return nil
}

func (k Tekton) apply(c *kubernetes.Cluster, ui *ui.UI, options kubernetes.InstallationOptions, upgrade bool) error {
	if out, err := helpers.KubectlApplyEmbeddedYaml(tektonAdminRoleYamlPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Installing %s failed:\n%s", tektonAdminRoleYamlPath, out))
	}
	if out, err := helpers.KubectlApplyEmbeddedYaml(tektonPipelineYamlPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Installing %s failed:\n%s", tektonPipelineYamlPath, out))
	}
	if out, err := helpers.KubectlApplyEmbeddedYaml(tektonTriggersYamlPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Installing %s failed:\n%s", tektonTriggersYamlPath, out))
	}
	if out, err := helpers.KubectlApplyEmbeddedYaml(tektonDashboardYamlPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Installing %s failed:\n%s", tektonDashboardYamlPath, out))
	}

	err := c.LabelNamespace(tektonNamespace, kubernetes.FusemlDeploymentLabelKey, kubernetes.FusemlDeploymentLabelValue)
	if err != nil {
		return err
	}

	for _, crd := range []string{
		"clustertasks.tekton.dev",
		"clustertriggerbindings.triggers.tekton.dev",
		"conditions.tekton.dev",
		"eventlisteners.triggers.tekton.dev",
		"pipelineresources.tekton.dev",
		"pipelineruns.tekton.dev",
		"pipelines.tekton.dev",
		"runs.tekton.dev",
		"taskruns.tekton.dev",
		"tasks.tekton.dev",
		"triggerbindings.triggers.tekton.dev",
		"triggers.triggers.tekton.dev",
		"triggertemplates.triggers.tekton.dev",
	} {
		message := fmt.Sprintf("Establish CRD %s", crd)
		out, err := helpers.WaitForCommandCompletion(ui, message,
			func() (string, error) {
				return helpers.Kubectl("wait --for=condition=established --timeout=" + strconv.Itoa(k.Timeout) + "s crd/" + crd)
			},
		)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("%s failed:\n%s", message, out))
		}
	}

	for _, c := range []string{"pipelines", "triggers", "dashboard"} {
		message := fmt.Sprintf("Starting tekton %s pods", c)
		out, err := helpers.WaitForCommandCompletion(ui, message,
			func() (string, error) {
				return helpers.Kubectl(fmt.Sprintf("wait --for=condition=Ready --timeout=%ds -n %s --selector=app.kubernetes.io/part-of=tekton-%s pod",
					k.Timeout, tektonNamespace, c))
			},
		)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("%s failed:\n%s", message, out))
		}
	}

	message := "Creating registry certificates in fuseml-workloads"
	out, err := helpers.WaitForCommandCompletion(ui, message,
		func() (string, error) {
			out1, err := helpers.ExecToSuccessWithTimeout(
				func() (string, error) {
					return helpers.Kubectl("get secret -n fuseml-workloads registry-tls-self-ca")
				}, time.Duration(k.Timeout)*time.Second, 3*time.Second)
			if err != nil {
				return out1, err
			}

			out2, err := helpers.ExecToSuccessWithTimeout(
				func() (string, error) {
					return helpers.Kubectl("get secret -n fuseml-workloads registry-tls-self")
				}, time.Duration(k.Timeout)*time.Second, 3*time.Second)

			return fmt.Sprintf("%s\n%s", out1, out2), err
		},
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s failed:\n%s", message, out))
	}

	message = "Installing FuseML pipelines and triggers"
	out, err = helpers.WaitForCommandCompletion(ui, message,
		func() (string, error) {
			return helpers.KubectlApplyEmbeddedYaml(tektonFusemlYamlPath)
		},
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s failed:\n%s", message, out))
	}

	message = "Applying tekton Kaniko resources"
	out, err = helpers.WaitForCommandCompletion(ui, message,
		func() (string, error) {
			return applyTektonKaniko(c, ui)
		},
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s failed:\n%s", message, out))
	}

	domain, err := options.GetString("system_domain", TektonDeploymentID)
	if err != nil {
		return errors.Wrap(err, "Couldn't get system_domain option")
	}

	if c.HasIstio() {
		message := "Creating Tekton dashboard istio ingress gateway"
		_, err = helpers.WaitForCommandCompletion(ui, message,
			func() (string, error) {
				return helpers.CreateIstioIngressGateway("tekton", tektonNamespace, TektonDeploymentID+"."+domain, "tekton-dashboard", 9097)
			},
		)
	} else {
		message = "Creating Tekton dashboard ingress"
		_, err = helpers.WaitForCommandCompletion(ui, message,
			func() (string, error) {
				return "", createTektonIngress(c, TektonDeploymentID+"."+domain)
			},
		)
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s failed", message))
	}

	if err := c.WaitUntilPodBySelectorExist(ui, WorkloadsDeploymentID, "eventlistener=mlflow-listener,app.kubernetes.io/part-of=Triggers", k.Timeout); err != nil {
		return errors.Wrap(err, "failed waiting Tekton event listener deployment to exist")
	}
	if err := c.WaitForPodBySelectorRunning(ui, WorkloadsDeploymentID, "eventlistener=mlflow-listener,app.kubernetes.io/part-of=Triggers", k.Timeout); err != nil {
		return errors.Wrap(err, "failed waiting Tekton event listener deployment to come up")
	}

	ui.Success().Msg("Tekton deployed")

	return nil
}

func (k Tekton) GetVersion() string {
	return fmt.Sprintf("pipelines: %s, triggers %s, dashboard: %s",
		tektonPipelineYamlPath, tektonTriggersYamlPath, tektonDashboardYamlPath)
}

func (k Tekton) Deploy(c *kubernetes.Cluster, ui *ui.UI, options kubernetes.InstallationOptions) error {

	_, err := c.Kubectl.CoreV1().Namespaces().Get(
		context.Background(),
		TektonDeploymentID,
		metav1.GetOptions{},
	)
	if err == nil {
		ui.Exclamation().Msg("Namespace " + TektonDeploymentID + " already present, skipping installation")
		return nil
	}

	ui.Note().KeeplineUnder(1).Msg("Deploying Tekton...")

	err = k.apply(c, ui, options, false)
	if err != nil {
		return err
	}

	return nil
}

func (k Tekton) Upgrade(c *kubernetes.Cluster, ui *ui.UI, options kubernetes.InstallationOptions) error {
	_, err := c.Kubectl.CoreV1().Namespaces().Get(
		context.Background(),
		TektonDeploymentID,
		metav1.GetOptions{},
	)
	if err != nil {
		return errors.New("Namespace " + TektonDeploymentID + " not present")
	}

	ui.Note().Msg("Upgrading Tekton...")

	return k.apply(c, ui, options, true)
}

// The equivalent of:
// kubectl get secret -n fuseml-workloads registry-tls-self -o json | jq -r '.["data"]["ca"]' | base64 -d | openssl x509 -hash -noout
// written in golang.
func getRegistryCAHash(c *kubernetes.Cluster, ui *ui.UI) (string, error) {
	secret, err := c.Kubectl.CoreV1().Secrets("fuseml-workloads").
		Get(context.Background(), "registry-tls-self", metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return helpers.OpenSSLSubjectHash(string(secret.Data["ca"]))
}

func applyTektonKaniko(c *kubernetes.Cluster, ui *ui.UI) (string, error) {
	caHash, err := getRegistryCAHash(c, ui)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get registry CA from fuseml-workloads namespace")
	}

	yamlPathOnDisk, err := helpers.ExtractFile(tektonKanikoYamlPath)
	if err != nil {
		return "", errors.New("Failed to extract embedded file: " + tektonKanikoYamlPath + " - " + err.Error())
	}
	defer os.Remove(yamlPathOnDisk)

	fileContents, err := ioutil.ReadFile(yamlPathOnDisk)
	if err != nil {
		return "", err
	}

	// Constructing the name of the cert file as required by openssl.
	// Lookup "subject_hash" in the docs: https://www.openssl.org/docs/man1.0.2/man1/x509.html
	re := regexp.MustCompile(`{{CA_SELF_HASHED_NAME}}`)
	renderedFileContents := re.ReplaceAll(fileContents, []byte(caHash+".0"))

	tmpFilePath, err := helpers.CreateTmpFile(string(renderedFileContents))
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFilePath)

	return helpers.Kubectl(fmt.Sprintf("apply -n fuseml-workloads --filename %s", tmpFilePath))
}

func createTektonIngress(c *kubernetes.Cluster, subdomain string) error {
	_, err := c.Kubectl.ExtensionsV1beta1().Ingresses("tekton-pipelines").Create(
		context.Background(),
		// TODO: Switch to networking v1 when we don't care about <1.18 clusters
		// Like this (which has been reverted):
		// https://github.com/SUSE/carrier/commit/7721d610fdf27a79be980af522783671d3ffc198
		&v1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "tekton-dashboard",
				Namespace: "tekton-pipelines",
				Annotations: map[string]string{
					"kubernetes.io/ingress.class": "traefik",
				},
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: subdomain,
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "/",
										Backend: v1beta1.IngressBackend{
											ServiceName: "tekton-dashboard",
											ServicePort: intstr.IntOrString{
												Type:   intstr.Int,
												IntVal: 9097,
											},
										}}}}}}}}},
		metav1.CreateOptions{},
	)

	return err
}
