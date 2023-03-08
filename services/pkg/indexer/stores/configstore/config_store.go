package configstore

import (
	"github.com/playmint/ds-node/pkg/api/model"
)

type ConfigStore struct{}

func New() *ConfigStore {
	return &ConfigStore{}
}

func (s *ConfigStore) GetContracts() []*model.ContractConfig {
	return []*model.ContractConfig{}
}
