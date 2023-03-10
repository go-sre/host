package accessdata

func originValue(value string) string {
	switch value {
	case OriginRegionOperator:
	case OriginZoneOperator:
	case OriginSubZoneOperator:
	case OriginServiceOperator:
	case OriginInstanceIdOperator:
	}
	return ""
}
