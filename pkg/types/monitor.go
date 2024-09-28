package types

type Load struct {
	Load1  []float64 `json:"load1"`
	Load5  []float64 `json:"load5"`
	Load15 []float64 `json:"load15"`
}

type CPU struct {
	Percent []string `json:"percent"`
}

type Mem struct {
	Total     string   `json:"total"`
	Available []string `json:"available"`
	Used      []string `json:"used"`
}

type SWAP struct {
	Total string   `json:"total"`
	Used  []string `json:"used"`
	Free  []string `json:"free"`
}

type Network struct {
	Sent []string `json:"sent"`
	Recv []string `json:"recv"`
	Tx   []string `json:"tx"`
	Rx   []string `json:"rx"`
}

type MonitorData struct {
	Times []string `json:"times"`
	Load  Load     `json:"load"`
	CPU   CPU      `json:"cpu"`
	Mem   Mem      `json:"mem"`
	SWAP  SWAP     `json:"swap"`
	Net   Network  `json:"net"`
}
