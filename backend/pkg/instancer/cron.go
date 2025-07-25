package instancer

import (
	"strings"
	"time"

	"github.com/DragonSecSI/instancer/backend/pkg/database/models"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/registry"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func (inst *Instancer) CronRun() error {
	debugLog := func(format string, v ...any) {
		inst.Logger.Debug().Msgf(format, v...)
	}

	opts := genericclioptions.NewConfigFlags(false)
	opts.APIServer = &inst.KubernetesConfig.AuthURL
	opts.BearerToken = &inst.KubernetesConfig.AuthToken
	opts.CAFile = &inst.KubernetesConfig.AuthCA
	opts.Namespace = &inst.KubernetesConfig.Namespace

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(opts, inst.KubernetesConfig.Namespace, "secret", debugLog); err != nil {
		inst.Logger.Error().Err(err).Msg("Failed to initialize Helm action configuration")
		return err
	}

	regClient, err := registry.NewClient(
		registry.ClientOptBasicAuth(inst.HelmConfig.BasicAuthUsername, inst.HelmConfig.BasicAuthPassword),
	)
	if err != nil {
		inst.Logger.Error().Err(err).Msg("Failed to create Helm registry client")
		return err
	}

	actionConfig.RegistryClient = regClient

	// List releases
	releases, err := actionConfig.Releases.ListReleases()
	if err != nil {
		inst.Logger.Error().Err(err).Msg("Failed to list Helm releases")
		return err
	}
	if releases == nil {
		return nil // No releases to process
	}

	// Delete releases that match the prefix
	for _, release := range releases {
		if strings.HasPrefix(release.Name, inst.Prefix) {
			if time.Since(release.Info.LastDeployed.Time) > inst.CleanupDuration {
				go func() {
					uninstall := action.NewUninstall(actionConfig)
					_, err := uninstall.Run(release.Name)
					if err != nil {
						inst.Logger.Error().Err(err).Msgf("Failed to uninstall release: %s", release.Name)
					} else {
						inst.Logger.Info().Msgf("Successfully uninstalled release: %s", release.Name)
					}

					instance, err := models.InstanceGetByName(inst.DB, release.Name)
					if err != nil {
						inst.Logger.Error().Err(err).Msgf("Failed to get instance by name: %s", release.Name)
						return
					}
					if inst != nil {
						instance.Active = false
						if err := models.InstanceUpdate(inst.DB, instance); err != nil {
							inst.Logger.Error().Err(err).Msgf("Failed to update instance status for: %s", release.Name)
						} else {
							inst.Logger.Info().Msgf("Updated instance status to inactive for: %s", release.Name)
						}
					} else {
						inst.Logger.Warn().Msgf("Instance not found for release: %s", release.Name)
					}
				}()
			}
		}
	}

	return nil
}
