package deployments

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fuseml/fuseml/cli/helpers"
	"github.com/fuseml/fuseml/cli/kubernetes"
	"github.com/fuseml/fuseml/cli/paas/ui"
	"github.com/kyokomi/emoji"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Gitea struct {
	Debug   bool
	Timeout int
}

const (
	GiteaDeploymentID = "gitea"
	giteaVersion      = "1.15.3"
	giteaChartURL     = "https://dl.gitea.io/charts/gitea-4.1.1.tgz"
)

func (k *Gitea) ID() string {
	return GiteaDeploymentID
}

func (k *Gitea) Backup(c *kubernetes.Cluster, ui *ui.UI, d string) error {
	return nil
}

func (k *Gitea) Restore(c *kubernetes.Cluster, ui *ui.UI, d string) error {
	return nil
}

func (k Gitea) Describe() string {
	return emoji.Sprintf(":cloud:Gitea version: %s\n:clipboard:Gitea chart: %s", giteaVersion, giteaChartURL)
}

// Delete removes Gitea from kubernetes cluster
func (k Gitea) Delete(c *kubernetes.Cluster, ui *ui.UI) error {
	ui.Note().KeeplineUnder(1).Msg("Removing Gitea...")

	existsAndOwned, err := c.NamespaceExistsAndOwned(GiteaDeploymentID)
	if err != nil {
		return errors.Wrapf(err, "failed to check if namespace '%s' is owned or not", GiteaDeploymentID)
	}
	if !existsAndOwned {
		ui.Exclamation().Msg("Skipping Gitea because namespace either doesn't exist or not owned by Fuseml")
		return nil
	}

	currentdir, err := os.Getwd()
	if err != nil {
		return errors.New("Failed uninstalling Gitea: " + err.Error())
	}

	message := "Removing helm release " + GiteaDeploymentID
	out, err := helpers.WaitForCommandCompletion(ui, message,
		func() (string, error) {
			helmCmd := fmt.Sprintf("helm uninstall gitea --namespace %s", GiteaDeploymentID)
			return helpers.RunProc(helmCmd, currentdir, k.Debug)
		},
	)
	if err != nil {
		if strings.Contains(out, "release: not found") {
			ui.Exclamation().Msgf("%s helm release not found, skipping.\n", GiteaDeploymentID)
		} else {
			return errors.Wrapf(err, "Failed uninstalling helm release %s: %s", GiteaDeploymentID, out)
		}
	}

	message = "Deleting Gitea namespace " + GiteaDeploymentID
	_, err = helpers.WaitForCommandCompletion(ui, message,
		func() (string, error) {
			return "", c.DeleteNamespace(GiteaDeploymentID)
		},
	)
	if err != nil {
		return errors.Wrapf(err, "Failed deleting namespace %s", GiteaDeploymentID)
	}

	ui.Success().Msg("Gitea removed")

	return nil
}

func (k Gitea) apply(c *kubernetes.Cluster, ui *ui.UI, options kubernetes.InstallationOptions, upgrade bool) error {
	action := "install"
	if upgrade {
		action = "upgrade"
	}

	currentdir, err := os.Getwd()
	if err != nil {
		return err
	}
	if action == "install" {
		helmCmd := fmt.Sprintf("helm list --namespace %s --deployed -q | grep gitea", GiteaDeploymentID)
		out, _ := helpers.RunProc(helmCmd, currentdir, k.Debug)
		if strings.TrimSpace(out) == "gitea" {
			ui.Exclamation().Msg("gitea already present under " + GiteaDeploymentID + " namespace, skipping installation")
			return nil
		}
	}

	// Setup Gitea helm values
	var helmArgs []string

	domain, err := options.GetString("system_domain", GiteaDeploymentID)
	if err != nil {
		return err
	}
	subdomain := GiteaDeploymentID + "." + domain

	hasIstio := c.HasIstio()

	config := fmt.Sprintf(`
ingress:
  enabled: %t
  hosts:
    - host: %s
      paths:
       - path: /
         pathType: Prefix
  annotations:
    kubernetes.io/ingress.class: traefik

gitea:
  admin:
    username: "dev"
    password: "changeme"
    email: "admin@fuseml.sh"
  config:
    RUN_MODE: prod
    repository:
      ROOT:  "/data/git/gitea-repositories"
    server:
      DOMAIN:  %s
      ROOT_URL: %s
    security:
      INSTALL_LOCK: true
      SECRET_KEY: generated-by-quarks-secret
      INTERNAL_TOKEN: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE2MDIzNzc3NzZ9.uvJPCMSDTPlVMAUwNzW9Jbl5487kbj5T_pWu3dGirnA
    service:
      ENABLE_REGISTRATION_CAPTCHA: false
      DISABLE_REGISTRATION: true
    oauth2:
      ENABLE: true
      JWT_SECRET: HLNn92qqtznZSMkD_TzR_XFVdiZ5E87oaus6pyH7tiI
`, !hasIstio, subdomain, subdomain, "http://"+subdomain)

	configPath, err := helpers.CreateTmpFile(config)
	if err != nil {
		return err
	}
	defer os.Remove(configPath)

	helmCmd := fmt.Sprintf("helm %s gitea --create-namespace --values %s --namespace %s %s %s", action, configPath, GiteaDeploymentID, giteaChartURL, strings.Join(helmArgs, " "))

	if out, err := helpers.RunProc(helmCmd, currentdir, k.Debug); err != nil {
		return errors.New("Failed installing Gitea: " + out)
	}
	err = c.LabelNamespace(GiteaDeploymentID, kubernetes.FusemlDeploymentLabelKey, kubernetes.FusemlDeploymentLabelValue)
	if err != nil {
		return err
	}

	if hasIstio {
		message := "Creating istio ingress gateway"
		out, err := helpers.WaitForCommandCompletion(ui, message,
			func() (string, error) {
				return helpers.CreateIstioIngressGateway("gitea", GiteaDeploymentID, subdomain, "gitea-http", 3000)
			},
		)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("%s failed:\n%s", message, out))
		}
	}

	for _, podname := range []string{
		"memcached",
		"postgresql",
		"gitea",
	} {
		if err := c.WaitUntilPodBySelectorExist(ui, GiteaDeploymentID, "app.kubernetes.io/name="+podname, k.Timeout); err != nil {
			return errors.Wrap(err, "failed waiting Gitea "+podname+" deployment to exist")
		}
		if err := c.WaitForPodBySelectorRunning(ui, GiteaDeploymentID, "app.kubernetes.io/name="+podname, k.Timeout); err != nil {
			return errors.Wrap(err, "failed waiting Gitea "+podname+" deployment to come up")
		}
	}

	ui.Success().Msg(fmt.Sprintf("Gitea deployed (http://%s).", subdomain))

	return nil
}

func (k Gitea) GetVersion() string {
	return giteaVersion
}

func (k Gitea) Deploy(c *kubernetes.Cluster, ui *ui.UI, options kubernetes.InstallationOptions) error {

	_, err := c.Kubectl.CoreV1().Namespaces().Get(
		context.Background(),
		GiteaDeploymentID,
		metav1.GetOptions{},
	)
	if err == nil {
		ui.Note().Msg("Namespace " + GiteaDeploymentID + " already present")
	}

	ui.Note().KeeplineUnder(1).Msg("Deploying Gitea...")

	err = k.apply(c, ui, options, false)
	if err != nil {
		return err
	}

	return nil
}

func (k Gitea) Upgrade(c *kubernetes.Cluster, ui *ui.UI, options kubernetes.InstallationOptions) error {
	_, err := c.Kubectl.CoreV1().Namespaces().Get(
		context.Background(),
		GiteaDeploymentID,
		metav1.GetOptions{},
	)
	if err != nil {
		return errors.New("Namespace " + GiteaDeploymentID + " not present")
	}

	ui.Note().Msg("Upgrading Gitea...")

	return k.apply(c, ui, options, true)
}
