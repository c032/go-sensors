//go:build alltests || battery

package sensors_test

import (
	"testing"

	"github.com/c032/go-sensors"
)

// TestLinuxBatterySource tests for a 100% battery capacity while currently
// charging.
//
// To check that it works for different states just unplug the charger, re-run
// the tests, and see that it fails with expected values.
func TestLinuxBatterySource(t *testing.T) {
	lbs := sensors.LinuxBatterySource

	var (
		err       error
		batteries []sensors.Battery
	)

	batteries, err = lbs.Batteries()
	if err != nil {
		t.Fatalf("LinuxBatterySource.Batteries() error: %s", err)
	}
	if got, want := len(batteries), 1; got != want {
		t.Fatalf("len(LinuxBatterySource.Batteries()) = %#v; want %#v", got, want)
	}

	b := batteries[0]

	var capacity float64
	if capacity, err = b.Capacity(); err == nil {
		if got, want := capacity, 100.0; got != want {
			t.Errorf("(*linuxBattery).Capacity() = %#v; want %#v", got, want)
		}
	} else {
		t.Errorf("(*linuxBattery).Capacity() error: %s", err)
	}

	var status string
	if status, err = b.Status(); err == nil {
		if got, want := status, "Charging"; got != want {
			t.Errorf("(*linuxBattery).Status() = %#v; want %#v", got, want)
		}
	} else {
		t.Errorf("(*linuxBattery).Status() error: %s", err)
	}

	bsp := b.(sensors.BatteryStatusProvider)

	if isCharging, err := bsp.IsCharging(); err == nil {
		if got, want := isCharging, true; got != want {
			t.Errorf("(*linuxBattery).IsCharging() = %#v; want %#v", got, want)
		}
	}

	if isDischarging, err := bsp.IsDischarging(); err == nil {
		if got, want := isDischarging, false; got != want {
			t.Errorf("(*linuxBattery).IsDischarging() = %#v; want %#v", got, want)
		}
	}
}
