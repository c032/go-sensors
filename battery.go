package sensors

type Battery interface {
	Capacity() (float64, error)
	Status() (string, error)
}

type BatteryStatusProvider interface {
	IsCharging() (bool, error)
	IsDischarging() (bool, error)
}

type BatterySource interface {
	Batteries() ([]Battery, error)
}
