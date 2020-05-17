//+build wireinject

package main

import "github.com/google/wire"

func InitializeScript() (*RInfo, error) {
	wire.Build(NewRInfo, NewConfig, NewSource)
	return &RInfo{}, nil
}

func InitializeConfig() (*Config, error) {
	wire.Build(NewConfig)
	return &Config{}, nil
}

func InitializeClient() (*Source, error) {
	wire.Build(NewSource, NewConfig)
	return &Source{}, nil
}
