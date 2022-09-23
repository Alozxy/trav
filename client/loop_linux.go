package main

import "time"

func start() {
	if conf.get_conf("udp_mode").(bool) {

		go udp_loop()
	} else {

		go syn_loop()
	}

	var external_port uint16 = 0
	for {

		request(&external_port)
		time.Sleep(time.Duration(get_conf("interval").(int)) * time.Second)
	}
}
