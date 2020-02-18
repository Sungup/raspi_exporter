package common

func CheckPrerequisite(opts *RaspiExpOpts) error {
	if !opts.Debug {
		// TODO 1. Check vcgencmd exists

		// TODO 2. Check /sys/class/thermal/thermal_zone0/temp exists

	}

	return nil
}
