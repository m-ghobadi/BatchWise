package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Weights                Weights    `yaml:"weights"`
	Thresholds             Thresholds `yaml:"thresholds"`
	BatchSizeLimits        Limits     `yaml:"batch_size_limits"`
	IntervalLimits         Limits     `yaml:"interval_limits"`
	SamplingInterval       int        `yaml:"sampling_interval"`
	ProcessingIntervalBase float64    `yaml:"processing_interval_base"`
	Constants              Constants  `yaml:"constants"`
	StaticBatchSize        int        `yaml:"static_batch_size"`
	WorkerCount            int        `yaml:"worker_count"`
}

type Weights struct {
	W1 float64 `yaml:"w1"`
	W2 float64 `yaml:"w2"`
	W3 float64 `yaml:"w3"`
	W4 float64 `yaml:"w4"`
}

type Thresholds struct {
	Priority float64 `yaml:"priority"`
}

type Limits struct {
	Min float64 `yaml:"min"`
	Max float64 `yaml:"max"`
}

type Constants struct {
	Alpha float64 `yaml:"alpha"`
	Beta  float64 `yaml:"beta"`
	Gamma float64 `yaml:"gamma"`
	C     float64 `yaml:"c"`
}

func LoadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
