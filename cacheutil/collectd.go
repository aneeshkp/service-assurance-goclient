package cacheutil

import (
	"encoding/json"
	"fmt"
  "math/rand"
	"strconv"
  "time"
)

type Collectd struct {
	Values          []float64
	Dstypes         []string
	Dsnames         []string
	Time            float64
	Interval        float64
	Host            string
	Plugin          string
	Plugin_instance string
	Type            string `json:"type"`
	Type_instance   string
	new						bool
}
//GenrateSampleData
func GenrateSampleData(hostname string, plugincount int, collectdjson string, cacheserver *CacheServer) {
	//100 plugins
	for j := 0; j < plugincount; j++ {
		var pluginname = fmt.Sprintf("%s_%d", "plugin_name", j)
		go func() {
			c := ParseCollectdJSON(collectdjson)
			c.Host = hostname
			c.Plugin = pluginname
			c.Type = pluginname
			c.Plugin_instance = pluginname
			c.Dstypes[0] = "gauge"
			c.Dstypes[1] = "gauge"
			c.Dsnames[0] = "value1"
			c.Dsnames[1] = "value2"
			c.Values[0] = rand.Float64()
			c.Values[1] = rand.Float64()
			c.Time = float64((time.Now().UnixNano())) / 1000000
			cacheserver.Put(*c)
		}()
	}
}
//ParseCollectdJSON   ...
func ParseCollectdJSON(collectdJson string) *Collectd {
	c:=make([]Collectd,1)
	var jsonBlob = []byte(collectdJson)
	err := json.Unmarshal(jsonBlob, &c)
	if err != nil {
		fmt.Println("error:", err)
	}
	c1:=c[0]
	c1.SetNew(true)
	return &c1

}
//ParseCollectdByte ....
func ParseCollectdByte(amqpCollectd []byte) *Collectd {
	c:=make([]Collectd,1)
	//var jsonBlob = []byte(collectdJson)
	err := json.Unmarshal(amqpCollectd, &c)
	if err != nil {
		fmt.Println("error:", err)
	}
	c1:=c[0]
	c1.SetNew(true)
	return &c1

}

func (c *Collectd) SetNew(new bool) {
	c.new=new
}
func (c *Collectd) ISNew() bool {
	return c.new
}

// newName converts one data source of a value list to a string representation.
func (vl *Collectd) DSName(index int) string {
	if vl.Dsnames != nil {
		return vl.Dsnames[index]
	} else if len(vl.Values) != 1 {
		return strconv.FormatInt(int64(index), 10)
	}

	return "value"
}

//generateCollectdJson   for samples
func GenerateCollectdJson(hostname string, pluginname string) string {
	return `[{
      "values":  [0.0,0.0],
      "dstypes":  ["gauge","guage"],
      "dsnames":    ["value1","value2"],
      "time":      0,
      "interval":          10,
      "host":            "hostname",
      "plugin":          "pluginname",
      "plugin_instance": "0",
      "type":            "pluginname",
      "type_instance":   "idle"
    }]`
}