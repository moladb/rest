package rest

type Config struct {
	BindAddr              string
	EnableDebug           bool
	EnableMetrics         bool
	EnableDiscovery       bool
	GraceShutdownTimeoutS int
}

type Version struct {
	Version   string
	BuildDate string
	GoVersion string
}
