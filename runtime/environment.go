package runtime

import (
	"os"
	"strings"
)

const (
	DevEnv = iota
	ReviewEnv
	TestEnv
	StageEnv
	ProdEnv

	runtimeEnvKey  = "RUNTIME_ENV"
	DevEnvValue    = "dev"
	ReviewEnvValue = "review"
	TestEnvValue   = "test"
	StageEnvValue  = "stage"
	ProdEnvValue   = "prod"
)

func matchEnvironment(env int) bool {
	s := GetRuntimeEnv()
	switch env {
	case DevEnv:
		return strings.EqualFold(s, DevEnvValue)
	case ReviewEnv:
		return strings.EqualFold(s, ReviewEnvValue)
	case TestEnv:
		return strings.EqualFold(s, TestEnvValue)
	case StageEnv:
		return strings.EqualFold(s, StageEnvValue)
	case ProdEnv:
		return strings.EqualFold(s, ProdEnvValue)
	}
	return false
}

// IsDevEnv - determine environment
func IsDevEnv() bool {
	return matchEnvironment(DevEnv)
}

func IsReviewEnv() bool {
	return matchEnvironment(ReviewEnv)
}

func IsTestEnv() bool {
	return matchEnvironment(TestEnv)
}

func IsStageEnv() bool {
	return matchEnvironment(StageEnv)
}

func IsProdEnv() bool {
	return matchEnvironment(ProdEnv)
}

// GetRuntimeEnv - function to get the runtime environment
func GetRuntimeEnv() string {
	s := os.Getenv(runtimeEnvKey)
	if s == "" {
		return DevEnvValue
	}
	return s
}

// SetRuntimeEnv - function to set the runtime environment
func SetRuntimeEnv(s string) {
	os.Setenv(runtimeEnvKey, s)
}

const (
	EnvVariableReference = "{env}"
)

// EnvExpansion - expand a template with an environment variable reference
func EnvExpansion(s string) string {
	index := strings.Index(s, EnvVariableReference)
	if index == -1 {
		return "invalid or missing environment variable reference: {env}"
	}
	t := s[:index] + GetRuntimeEnv()
	u := s[index+len(EnvVariableReference):]
	return t + u
}
