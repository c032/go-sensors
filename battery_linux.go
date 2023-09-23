package sensors

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

var _ BatterySource = (*linuxBatterySource)(nil)

func isBatteryEntryName(name string) bool {
	const prefix = "BAT"

	if !strings.HasPrefix(name, prefix) {
		return false
	}

	nStr := strings.TrimPrefix(name, prefix)
	for _, c := range nStr {
		if c >= '0' && c <= '9' {
			continue
		}

		return false
	}

	return true
}

var (
	_ Battery               = (*linuxBattery)(nil)
	_ BatteryStatusProvider = (*linuxBattery)(nil)
)

type linuxBattery struct {
	Path string
}

func (lb *linuxBattery) Capacity() (float64, error) {
	p := filepath.Join(lb.Path, "capacity")

	var (
		err         error
		capacityStr string
	)
	capacityStr, err = readEnoughUTF8FromFile(p)
	if err != nil {
		return 0.0, fmt.Errorf("could not read battery capacity: %w", err)
	}
	capacityStr = strings.TrimSpace(capacityStr)

	var capacity float64
	capacity, err = strconv.ParseFloat(capacityStr, 64)
	if err != nil {
		return 0.0, fmt.Errorf("could not parse battery capacity: %w", err)
	}

	return capacity, nil
}

func (lb *linuxBattery) Status() (string, error) {
	p := filepath.Join(lb.Path, "status")

	var (
		err    error
		status string
	)
	status, err = readEnoughUTF8FromFile(p)
	if err != nil {
		return "", fmt.Errorf("could not read battery status: %w", err)
	}

	status = strings.TrimSpace(status)

	return status, nil
}

func (lb *linuxBattery) IsCharging() (bool, error) {
	status, err := lb.Status()
	if err != nil {
		return false, fmt.Errorf("could not check if battery is charging: %w", err)
	}

	isCharging := status == "Charging"

	return isCharging, nil
}

func (lb *linuxBattery) IsDischarging() (bool, error) {
	status, err := lb.Status()
	if err != nil {
		return false, fmt.Errorf("could not check if battery is discharging: %w", err)
	}

	isDischarging := status == "Discharging"

	return isDischarging, nil
}

var LinuxBatterySource = &linuxBatterySource{}

type linuxBatterySource struct{}

func (lbs *linuxBatterySource) Batteries() ([]Battery, error) {
	const powerSupplyDir = "/sys/class/power_supply"

	var (
		err     error
		entries []fs.DirEntry
	)

	entries, err = os.ReadDir(powerSupplyDir)
	if err != nil {
		return nil, fmt.Errorf("could not list batteries: %w", err)
	}

	var batteries []Battery
	for _, entry := range entries {
		name := entry.Name()
		if !isBatteryEntryName(name) {
			continue
		}

		path := filepath.Join(powerSupplyDir, name)
		battery := &linuxBattery{
			Path: path,
		}

		batteries = append(batteries, battery)
	}

	return batteries, nil
}

func readEnough(r io.Reader) ([]byte, error) {
	const readAtMostBytes = 512
	buf := make([]byte, readAtMostBytes)

	br := io.LimitReader(r, readAtMostBytes)

	n, err := br.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading from file: %w", err)
	}

	result := buf[:n]

	return result, nil
}

func readEnoughFromFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("could not open file to read: %w", err)
	}
	defer f.Close()

	return readEnough(f)
}

func readEnoughUTF8FromFile(file string) (string, error) {
	var (
		err     error
		rawData []byte
	)

	rawData, err = readEnoughFromFile(file)
	if err != nil {
		return "", fmt.Errorf("could not read bytes from file: %w", err)
	}

	if !utf8.Valid(rawData) {
		return "", fmt.Errorf("bytes read from file are not utf8: %#v", file)
	}

	utf8Data := string(rawData)

	return utf8Data, nil
}
