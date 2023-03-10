package shared

import (
	"fmt"
)

func init() {
	SetOrigin(Origin{
		Region:     "Region",
		Zone:       "Zone",
		SubZone:    "SubZone",
		Service:    "Service",
		InstanceId: "InstanceId",
	})
}

func Example_Origin() {
	e := NewEmptyEntry()
	fmt.Printf("test: Value(Region) -> %v\n", e.Value(OriginRegionOperator))
	fmt.Printf("test: Value(Zone) -> %v\n", e.Value(OriginZoneOperator))
	fmt.Printf("test: Value(SubZone) -> %v\n", e.Value(OriginSubZoneOperator))
	fmt.Printf("test: Value(Service) -> %v\n", e.Value(OriginServiceOperator))
	fmt.Printf("test: Value(InstanceId) -> %v\n", e.Value(OriginInstanceIdOperator))

	//Output:
	//test: Value(Region) -> Region
	//test: Value(Zone) -> Zone
	//test: Value(SubZone) -> SubZone
	//test: Value(Service) -> Service
	//test: Value(InstanceId) -> InstanceId

}
