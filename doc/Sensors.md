# Sensors

## INA226 I2C Current Sensor

Discover the INA266 address and edit the constant in ```ina266monitor.go```

```sudo i2cdetect -y 1```

set ```kFinalPaI2cAddress``` to the address 

## DS18B20 Temperature Sensors

```
DEVICE               CONNECTION
PaSensorSlaveId      pin 7 GPIO_4
PreAmpSensorSlaveId  pin 7 GPIO_4
```

Discover the DS18B20 addresses and edit the constants in ```ds18b20monitor.go```

```ls /sys/bus/w1/devices/```

Set ```kPaSensorSlaveId``` and ```kPreAmpSensorSlaveId``` accordingly

## Fan Speed Sensors

```
DEVICE          CONNECTION
EncIntakePin    pin 29 GPIO_5
EncExtractPi    pin 31 GPIO_6
PaIntakePin     pin 33 GPIO_13
PaExtractPin    pin 35 GPIO_19
```

Fans do not have addresses, so nothing to discover
