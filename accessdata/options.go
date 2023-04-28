package accessdata

var origin Origin

func SetOrigin(region, zone, subZone, service, instanceId string) {
	origin.Region = region
	origin.Zone = zone
	origin.SubZone = subZone
	origin.Service = service
	origin.InstanceId = instanceId
}
