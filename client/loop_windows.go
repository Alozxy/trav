package main

import "time"

func start() {

	var external_port uint16 = 0
	if conf.get_conf("udp_mode").(bool) {

		go udp_loop()
		for {

			if conf.get_conf("ifconfig_me").(bool) {
				ifconfig_me_request(&external_port)
			} else {
				stun_request(&external_port)
			}

			time.Sleep(time.Duration(get_conf("interval").(int)) * time.Second)
		}

	} else {

		for {

			send_syn()
			if conf.get_conf("ifconfig_me").(bool) {
				ifconfig_me_request(&external_port)
			} else {
				stun_request(&external_port)
			}

			for i := 0; i < get_conf("interval").(int); i++ {

				send_syn()
				time.Sleep(1 * time.Second)
			}
		}
	}
}
