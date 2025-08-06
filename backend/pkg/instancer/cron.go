package instancer

import (
	"time"

	"github.com/DragonSecSI/instancer/backend/pkg/database/models"
	"github.com/DragonSecSI/instancer/backend/pkg/metrics"
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

	// Get active instances
	activeInstances, err := models.InstanceGetActive(inst.DB)
	if err != nil {
		inst.Logger.Error().Err(err).Msg("Failed to get active instances")
		return err
	}

	// Delete instances that are older than the cleanup duration
	count := 0
	for _, instance := range activeInstances {
		if time.Since(instance.CreatedAt).Seconds() > float64(instance.Duration) {
			inst.Logger.Info().Str("instance_name", instance.Name).Msg("Uninstalling instance")
			count++
			go func(inst *Instancer, instance *models.Instance) {
				uninstall := action.NewUninstall(actionConfig)
				_, err := uninstall.Run(instance.Name)
				if err != nil {
					inst.Logger.Error().Err(err).Msgf("Failed to uninstall instance: %s", instance.Name)
				}

				instance.Active = false
				if err := models.InstanceUpdate(inst.DB, instance); err != nil {
					inst.Logger.Error().Err(err).Msgf("Failed to update instance status for: %s", instance.Name)
				}

				metrics.InstancesDeletedCounter.Inc()
			}(inst, &instance)
		}
	}

	if count > 0 {
		inst.Logger.Info().Int("instances_deleted", count).Msg("Instance cleanup completed")
	}

	return nil
}
