package utils

func GetIp() []string {

	interfaces, err := GetUpnRunninginterfaces()
	var addrs []string
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	for _, interf := range interfaces {
		addr, _ := ExtractIPV4Address(interf)
		addrs = append(addrs, addr)
	}
	return addrs
}