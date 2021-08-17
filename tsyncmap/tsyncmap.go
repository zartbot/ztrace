//Package tsyncmap : a timeout based syncmap
package tsyncmap

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"sync"
	"time"
)

// Map is the base structure for tsyncmap.
type Map struct {
	Name       string
	Data       sync.Map
	Timeout    int64
	CheckFreq  int64
	ExpireTime sync.Map
	Verbose    bool
}

//NewMap is a construct function to create tsyncmap.
func NewMap(name string, timeout time.Duration, checkfreq time.Duration, verbose bool) *Map {
	t := int64(timeout)
	f := int64(checkfreq)

	return &Map{
		Name:      name,
		Timeout:   t,
		CheckFreq: f,
		Verbose:   verbose,
	}
}

//Load returns the value from tsyncmap
func (tmap *Map) Load(key interface{}) (value interface{}, ok bool) {
	return tmap.Data.Load(key)
}

type keyjson struct {
	Key string
}

func (tmap *Map) LoadRestApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	v := new(keyjson)
	_ = json.NewDecoder(r.Body).Decode(v)

	if v.Key == "internal_fetch_keylist" {
		var result []interface{}
		tmap.Data.Range(func(k, v interface{}) bool {
			result = append(result, k)
			return true
		})
		json.NewEncoder(w).Encode(result)
	} else {
		result, find := tmap.Data.Load(v.Key)
		if find {
			json.NewEncoder(w).Encode(result)
		} else {
			json.NewEncoder(w).Encode("")

		}
	}
}

func (tmap *Map) GetRemainTime(key interface{}) (time.Duration, error) {
	exp, ok := tmap.ExpireTime.Load(key)
	if ok {
		remainTime := exp.(time.Time).Sub(time.Now())
		return remainTime, nil
	}
	return 0 * time.Second, fmt.Errorf("key does not exist")
}

//Store is used save the key,value pairs in tsyncmap
func (tmap *Map) Store(key interface{}, value interface{}, currentTime time.Time) {
	//Check ExpireTime Map.
	exp, ok := tmap.ExpireTime.Load(key)
	if !ok {
		expireTime := currentTime.Add(time.Duration(tmap.Timeout))
		tmap.ExpireTime.Store(key, expireTime)
	} else {
		elapsed := exp.(time.Time).Sub(currentTime)
		//elapsed time less than half of timeout, update ExpireTime Store.
		if elapsed < time.Duration(tmap.Timeout/2) {
			expireTime := currentTime.Add(time.Duration(tmap.Timeout))
			tmap.ExpireTime.Store(key, expireTime)
		}
	}
	tmap.Data.Store(key, value)
}

//UpdateTime is used update specific key's expiretime.
func (tmap *Map) UpdateTime(key interface{}, currentTime time.Time) {
	expireTime := currentTime.Add(time.Duration(tmap.Timeout))
	tmap.ExpireTime.Store(key, expireTime)
}

func (tmap *Map) Delete(key interface{}) {
	tmap.Data.Delete(key)
	tmap.ExpireTime.Delete(key)
}

//Run is a coroutine to help tsyncmap manage the expire data.
func (tmap *Map) Run() {
	rand.Seed(time.Now().UnixNano())
	r := tmap.CheckFreq / 5
	for {
		currentTime := time.Now()
		tmap.ExpireTime.Range(func(k, v interface{}) bool {
			value := v.(time.Time)
			if value.Sub(currentTime) < 0 {
				//fmt.Println("DEBUG:::DELETE-KEY", reflect.ValueOf(k))
				tmap.Data.Delete(k)
				tmap.ExpireTime.Delete(k)
			}
			return true
		})
		if tmap.Verbose {
			tmap.ShowExpireTime()
			tmap.ShowData()
		}
		time.Sleep(time.Duration(tmap.CheckFreq + rand.Int63n(r)))
	}
}

func (tmap *Map) ShowExpireTime() {
	fmt.Printf("%10s:--------------------Expire Time Table-------------------------------\n", tmap.Name)
	i := 1
	tmap.ExpireTime.Range(func(k, v interface{}) bool {
		value := v.(time.Time)
		key := reflect.ValueOf(k)
		fmt.Printf("%10s:[%6d]Key:%v | ExipreTime: %v \n", tmap.Name, i, key, value)
		i++
		return true
	})
	fmt.Printf("%10s:--------------------------------------------------------------------\n\n", tmap.Name)
}

func (tmap *Map) ShowData() {
	fmt.Printf("%10s:---------------------Data Table-------------------------------------\n", tmap.Name)
	i := 1
	tmap.Data.Range(func(k, v interface{}) bool {
		value := reflect.ValueOf(v)
		key := reflect.ValueOf(k)
		fmt.Printf("%10s:[%6d]Key:%v | Value: %v \n", tmap.Name, i, key, value)
		i++
		return true
	})
	fmt.Printf("%10s:--------------------------------------------------------------------\n\n", tmap.Name)
}
