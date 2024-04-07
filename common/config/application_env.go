package config

type ApplicationEnv string

const (
	production  ApplicationEnv = "production"
	staging     ApplicationEnv = "staging"
	development ApplicationEnv = "development"
)

func (e ApplicationEnv) IsProduction() bool {
	return e == production
}

func (e ApplicationEnv) IsStaging() bool {
	return e == staging
}

func (e ApplicationEnv) IsDevelopment() bool {
	return e == development
}

func (e ApplicationEnv) String() string {
	switch e {
	case production:
		return "production"
	case staging:
		return "staging"
	case development:
		return "development"
	default:
		return ""
	}
}
