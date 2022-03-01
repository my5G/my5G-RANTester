package config

type Config struct {
	GNodeB struct {
		ControlIF struct {
			Ip   string `yaml:"ip"`
			Port int    `yaml:"port"`
		} `yaml:"controlif"`
		DataIF struct {
			Ip   string `yaml:"ip"`
			Port int    `yaml:"port"`
		} `yaml:"dataif"`
		PlmnList struct {
			Mcc   string `yaml:"mmc"`
			Mnc   string `yaml:"mnc"`
			Tac   string `yaml:"tac"`
			GnbId string `yaml:"gnbid"`
		} `yaml:"plmnlist"`
		SliceSupportList struct {
			Sst string `yaml:"sst"`
			Sd  string `yaml:"sd"`
		} `yaml:"slicesupportlist"`
	} `yaml:"gnodeb"`
	Ue struct {
		Msin  string `yaml:"msin"`
		Key   string `yaml:"key"`
		Opc   string `yaml:"opc"`
		Amf   string `yaml:"amf"`
		Sqn   string `yaml:"sqn"`
		Hplmn struct {
			Mcc string `yaml:"mcc"`
			Mnc string `yaml:"mnc"`
		} `yaml:"hplmn"`
		Snssai struct {
			Sst int    `yaml:"sst"`
			Sd  string `yaml:"sd"`
		} `yaml:"snssai"`
	} `yaml:"ue"`
	AMF struct {
		Ip   string `yaml:"ip"`
		Port int    `yaml:"port"`
		Name string `yaml:"name"`
	} `yaml:"amfif"`
}
