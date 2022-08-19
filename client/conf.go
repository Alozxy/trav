package main

type Config struct {
	kvstore map[string]interface{}
}

var conf *Config

func init() {

	c := new(Config)
	c.kvstore = make(map[string]interface{})

	conf = c
}

func set_conf(name string, value interface{}) {

	conf.set_conf(name, value)
}
func (conf *Config) set_conf(name string, value interface{}) {

	conf.kvstore[name] = value
}

func get_conf(name string) interface{} {

	v := conf.get_conf(name)
	return v
}
func (conf *Config) get_conf(name string) interface{} {

	v, exist := conf.kvstore[name]
	if !exist {
		return nil
	}
	return v
}
