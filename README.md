[![Build Status](https://travis-ci.org/adtechpotok/go-health.svg?branch=master)](https://travis-ci.org/adtechpotok/go-health)
[![codecov](https://codecov.io/gh/adtechpotok/go-health/branch/master/graph/badge.svg)](https://codecov.io/gh/adtechpotok/go-health)
[![Go Report Card](https://goreportcard.com/badge/github.com/adtechpotok/go-health)](https://goreportcard.com/report/github.com/adtechpotok/go-health)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/adtechpotok/go-health/master/LICENSE)

# Go health
result
```$xslt
{
"Version":"1.0",
"Lifetime":1307, 
"Errors":0,
"Warnings":0,
"Heartbeat":0, // fixing problem with closed binlog connection
"BinLogPosition":62532416,
"BinLogFile":"mysql-bin.000003",
"Additional":null, //any usefull data you need
"CacheState":
    {
    "Animals.Cats":27,
    "Superheroes.Marvel":12
    }
}
```

#Configuration
```$xslt
var daemonHealth = health.New(version)
var log = logrus.New() 

func init(){
	log.Hooks.Add(&health.HealthHook{daemonHealth})
}

func industrialHealth(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	daemonHealth.CacheState = t
	err := daemonHealth.Health()
	if err != nil {
		log.Error(err)
		return
	}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
        js, _ := jsoniter.Marshal(result)
    	w.Write(js)
}

```

add to catch binlog updates
```$xslt
func (h *binlogHandler) OnRow(e *canal.RowsEvent) error {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r, " ", string(debug.Stack()))
		}
	}()
	go daemonHealth.CanalListener(e)
	return nil
}
```