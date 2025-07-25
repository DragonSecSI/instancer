package instancer

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func (inst *Instancer) NewInstance(conf InstancerConfig) error {
	settings := cli.New()

	debugLog := func(format string, v ...any) {
		inst.Logger.Debug().Msgf(format, v...)
	}

	opts := genericclioptions.NewConfigFlags(false)
	opts.APIServer = &inst.KubernetesConfig.AuthURL
	//opts.BearerToken = &inst.KubernetesConfig.AuthToken
	//opts.CertFile = &inst.KubernetesConfig.AuthCert
	//opts.CAFile = &inst.KubernetesConfig.AuthCA
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

	install := action.NewInstall(actionConfig)
	install.ReleaseName = conf.Name
	install.Namespace = inst.KubernetesConfig.Namespace
	install.ChartPathOptions.Username = inst.HelmConfig.BasicAuthUsername
	install.ChartPathOptions.Password = inst.HelmConfig.BasicAuthPassword
	install.ChartPathOptions.Version = conf.Version

	chartPath, err := install.ChartPathOptions.LocateChart(conf.Repository, settings)
	if err != nil {
		inst.Logger.Error().Err(err).Msg("Failed to locate chart")
		return err
	}

	chartReq, err := loader.Load(chartPath)
	if err != nil {
		inst.Logger.Error().Err(err).Msg("Failed to load chart")
		return err
	}

	providers := getter.All(settings)
	valOpts := &values.Options{
		Values: conf.Values,
	}
	vals, err := valOpts.MergeValues(providers)
	if err != nil {
		inst.Logger.Error().Err(err).Msg("Failed to merge values")
		return err
	}

	rel, err := install.Run(chartReq, vals)
	if err != nil {
		inst.Logger.Error().Err(err).Msg("Failed to install chart")
		return err
	}

	inst.Logger.Info().Str("chart", conf.Repository).Str("release", rel.Name).Msg("Successfully installed Helm chart")
	return nil
}
