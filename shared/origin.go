package shared

func OriginValue(value string) string {
	switch value {
	case OriginRegionOperator:
		return opt.origin.Region
	case OriginZoneOperator:
		return opt.origin.Zone
	case OriginSubZoneOperator:
		return opt.origin.SubZone
	case OriginServiceOperator:
		return opt.origin.Service
	case OriginInstanceIdOperator:
		return opt.origin.InstanceId
	}
	return ""
}
