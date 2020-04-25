package ant

import "fmt"

const (
	deviceClassBikePowerSensor                   = "Bike Power Sensor"
	deviceClassControl                           = "Control"
	deviceClassFitnessEquipmentDevice            = "Fitness Equipment Device"
	deviceClassBloodPressureMonitor              = "Blood Pressure Monitor"
	deviceClassGeocacheTransmitter               = "Geocache Transmitter"
	deviceClassEnvironmentSensor                 = "Environment Sensor"
	deviceClassWeightSensor                      = "Weight Sensor"
	deviceClassHeartRateSensor                   = "Heart Rate Sensor"
	deviceClassBikeSpeedAndCadenceSensor         = "Bike Speed and Cadence Sensor"
	deviceClassBikeCadenceSensor                 = "Bike Cadence Sensor"
	deviceClassBikeSpeedSensor                   = "Bike Speed Sensor"
	deviceClassStrideBasedSpeedAndDistanceSensor = "Stride-Based Speed and Distance Sensor"
	deviceClassUnknown                           = "Unknown Device"
)

func deviceClasses() map[byte]string {
	return map[byte]string{
		11:  deviceClassBikePowerSensor,
		16:  deviceClassControl,
		17:  deviceClassFitnessEquipmentDevice,
		18:  deviceClassBloodPressureMonitor,
		19:  deviceClassGeocacheTransmitter,
		25:  deviceClassEnvironmentSensor,
		119: deviceClassWeightSensor,
		120: deviceClassHeartRateSensor,
		121: deviceClassBikeSpeedAndCadenceSensor,
		122: deviceClassBikeCadenceSensor,
		123: deviceClassBikeSpeedSensor,
		124: deviceClassStrideBasedSpeedAndDistanceSensor,
		255: deviceClassUnknown,
	}
}

func deviceClassDecoder(classes map[byte]string) func(b byte) (string, error) {
	return func(b byte) (string, error) {
		class, ok := classes[b]
		if !ok {
			return deviceClassUnknown, unknownDeviceError(b)
		}
		return class, nil
	}
}

type unknownDeviceError byte

func (dev unknownDeviceError) Error() string {
	return fmt.Sprintf("UnknownHandler device: %X", byte(dev))
}
