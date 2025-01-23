package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type config struct {
	tgBotToken           string   `yaml:"tgBotToken"`
	trackedChannels      []string `yaml:"trackedChannels"`
	cacheCapacity        int      `yaml:"cacheCapacity"`
	alertMessageTemplate string   `yaml:"alertMessageTemplate"`
}

func loadFromEnv() (*config, error) {
	botToken, ok := os.LookupEnv("TG_BOT_TOKEN")
	if !ok {
		return nil, errors.New("missing environment variable TG_BOT_TOKEN")
	}
	alertMessageTemplate, ok := os.LookupEnv("ALERT_MESSAGE_TEMPLATE")
	if !ok {
		return nil, errors.New("missing environment variable ALERT_TEXT_TEMPLATE")
	}
	trackedChannels, ok := os.LookupEnv("TRACKED_CHANNELS")
	if !ok {
		return nil, errors.New("missing environment variable TRACKED_CHANNELS")
	}
	capacityS, ok := os.LookupEnv("CACHE_CAPACITY")
	if !ok {
		return nil, errors.New("missing environment variable CACHE_CAPACITY")
	}
	capacity, err := strconv.Atoi(capacityS)
	if err != nil {
		return nil, errors.New("invalid CACHE_CAPACITY number")
	}
	return &config{
		tgBotToken:           botToken,
		trackedChannels:      strings.Fields(trackedChannels),
		cacheCapacity:        capacity,
		alertMessageTemplate: alertMessageTemplate,
	}, nil
}
