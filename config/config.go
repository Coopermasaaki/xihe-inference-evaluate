package config

import (
	"os"

	"github.com/opensourceways/community-robot-lib/utils"
	"gopkg.in/yaml.v2"

	"github.com/opensourceways/xihe-inference-evaluate/infrastructure/cloudimpl"
	"github.com/opensourceways/xihe-inference-evaluate/infrastructure/evaluateimpl"
	"github.com/opensourceways/xihe-inference-evaluate/infrastructure/inferenceimpl"
	"github.com/opensourceways/xihe-inference-evaluate/k8sclient"
)

func LoadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func LoadConfig(path string, cfg interface{}) error {
	if err := LoadFromYaml(path, cfg); err != nil {
		return err
	}

	if f, ok := cfg.(ConfigSetDefault); ok {
		f.SetDefault()
	}

	if f, ok := cfg.(ConfigValidate); ok {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type ConfigValidate interface {
	Validate() error
}

type ConfigSetDefault interface {
	SetDefault()
}

type Config struct {
	Evaluate  evaluateimpl.Config  `json:"evaluate"   required:"true"`
	Inference inferenceimpl.Config `json:"inference"  required:"true"`
	Cloud     cloudimpl.Config     `json:"cloud"      required:"true"`
	K8sClient k8sclient.Config     `json:"k8s"        required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Inference,
		&cfg.Evaluate,
		&cfg.Cloud,
		&cfg.K8sClient,
	}
}

func (cfg *Config) SetDefault() {
	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(ConfigSetDefault); ok {
			f.SetDefault()
		}
	}

	cfg.Evaluate.OBS = cfg.Inference.OBS
}

func (cfg *Config) Validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(ConfigValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
