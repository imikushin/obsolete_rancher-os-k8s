package config

func NewConfig() *Config {
	return &Config{
		Debug: DEBUG,
		Dns: []string{
			"8.8.8.8",
			"8.8.4.4",
		},
		State: ConfigState{
			Required: false,
			Dev:      "LABEL=RANCHER_STATE",
			FsType:   "auto",
		},
		SystemDockerArgs: []string{"docker", "-d", "-s", "overlay", "-b", "none", "--restart=false", "-H", DOCKER_SYSTEM_HOST},
		Modules:          []string{},
		Userdocker: UserDockerInfo{
			UseTLS: true,
		},
		Network: NetworkConfig{
			Interfaces: []InterfaceConfig{
				{
					Match: "eth*",
					DHCP:  true,
				},
				{
					Match:   "lo",
					Address: "127.0.0.1/8",
				},
			},
		},
		CloudInit: CloudInit{
			Datasources: []string{"configdrive:/media/config-2"},
		},
		SystemContainers: []ContainerConfig{
			{
				Id: "system-volumes",
				Cmd: "--name=system-volumes " +
					"--net=none " +
					"--read-only " +
					"-v=/var/lib/rancher/conf:/var/lib/rancher/conf " +
					"-v=/lib/modules:/lib/modules:ro " +
					"-v=/var/run:/var/run " +
					"-v=/var/log:/var/log " +
					"state",
			},
			{
				Id: "command-volumes",
				Cmd: "--name=command-volumes " +
					"--net=none " +
					"--read-only " +
					"-v=/init:/sbin/halt:ro " +
					"-v=/init:/sbin/poweroff:ro " +
					"-v=/init:/sbin/reboot:ro " +
					"-v=/init:/sbin/shutdown:ro " +
					"-v=/init:/sbin/netconf:ro " +
					"-v=/init:/usr/bin/cloud-init:ro " +
					"-v=/init:/usr/bin/rancherctl:ro " +
					"-v=/init:/usr/bin/respawn:ro " +
					"-v=/init:/usr/bin/system-docker:ro " +
					"-v=/lib/modules:/lib/modules:ro " +
					"-v=/usr/bin/docker:/usr/bin/docker:ro " +
					"state",
			},
			{
				Id: "user-volumes",
				Cmd: "--name=user-volumes " +
					"--net=none " +
					"--read-only " +
					"-v=/var/lib/rancher/state/home:/home " +
					"-v=/var/lib/rancher/state/opt:/opt " +
					"state",
			},
			{
				Id: "udev",
				Cmd: "--name=udev " +
					"--net=none " +
					"--privileged " +
					"--rm " +
					"-v=/dev:/host/dev " +
					"-v=/lib/modules:/lib/modules:ro " +
					"udev",
			},
			{
				Id: "network",
				Cmd: "--name=network " +
					"--rm " +
					"--cap-add=NET_ADMIN " +
					"--net=host " +
					"--volumes-from=command-volumes " +
					"--volumes-from=system-volumes " +
					"network",
			},
			{
				Id: "cloud-init",
				Cmd: "--name=cloud-init " +
					"--rm " +
					"--privileged " +
					"--net=host " +
					"--volumes-from=command-volumes " +
					"--volumes-from=system-volumes " +
					"cloudinit",
				ReloadConfig: true,
			},
			{
				Id: "ntp",
				Cmd: "--name=ntp " +
					"--rm " +
					"-d " +
					"--privileged " +
					"--net=host " +
					"ntp",
			},
			{
				Id: "syslog",
				Cmd: "--name=syslog " +
					"-d " +
					"--rm " +
					"--privileged " +
					"--net=host " +
					"--ipc=host " +
					"--pid=host " +
					"--volumes-from=system-volumes " +
					"syslog",
			},
			{
				Id: "userdocker",
				Cmd: "--name=userdocker " +
					"-d " +
					"--rm " +
					"--restart=always " +
					"--ipc=host " +
					"--pid=host " +
					"--net=host " +
					"--privileged " +
					"--volumes-from=command-volumes " +
					"--volumes-from=user-volumes " +
					"--volumes-from=system-volumes " +
					"-v=/var/lib/rancher/state/docker:/var/lib/docker " +
					"userdocker",
			},
			{
				Id: "console",
				Cmd: "--name=console " +
					"-d " +
					"--rm " +
					"--privileged " +
					"--volumes-from=command-volumes " +
					"--volumes-from=user-volumes " +
					"--volumes-from=system-volumes " +
					"--restart=always " +
					"--ipc=host " +
					"--net=host " +
					"--pid=host " +
					"console",
			},
		},
		EnabledAddons: []string{/*"ubuntu-console", "etcd", "flannel-conf", "flannel", "k8s-docker", "k8s-master", "k8s-minion"*/},
		Addons: map[string]Config{
			"ubuntu-console": {
				SystemContainers: []ContainerConfig{
					{
						Id: "console",
						Cmd: "--name=ubuntu-console " +
							"-d " +
							"--rm " +
							"--privileged " +
							"--volumes-from=command-volumes " +
							"--volumes-from=user-volumes " +
							"--volumes-from=system-volumes " +
							"--restart=always " +
							"--ipc=host " +
							"--net=host " +
							"--pid=host " +
							"rancher/ubuntuconsole:v0.0.2",
					},
				},
			},
			"etcd": {
				SystemContainers: []ContainerConfig{
					{
						Id: "etcd",
						Cmd: "--name=etcd " +
							"-d " +
							"--restart=always " +
							"--net=host " +
							"quay.io/coreos/etcd:v2.0.4 " +
							"--listen-client-urls 'http://0.0.0.0:2379,http://0.0.0.0:4001'",
					},
				},
			},
			"flannel-conf": {
				SystemContainers: []ContainerConfig{
					{
						Id: "flannel-conf",
						Cmd: "--name=flannel-conf " +
							"--rm " +
							"--net=host " +
							"--volumes-from=command-volumes " +
							"--volumes-from=system-volumes " +
							"etcdctl " +
							"mk /coreos.com/network/config '{\"Network\":\"10.244.0.0/16\", \"Backend\": {\"Type\": \"vxlan\"}}'",
					},
				},
			},
			"flannel": {
				SystemContainers: []ContainerConfig{
					{
						Id: "flannel",
						Cmd: "--name=flannel " +
							"-d " +
							"--rm " +
							"--restart=always " +
							"--ipc=host " +
							"--pid=host " +
							"--net=host " +
							"--privileged " +
							"--volumes-from=command-volumes " +
							"--volumes-from=system-volumes " +
							"flannel",
					},
				},
			},
			"k8s-docker": {
				SystemContainers: []ContainerConfig{
					{
						Id: "stopuserdocker",
						Cmd: "--rm " +
							"--ipc=host " +
							"--pid=host " +
							"--net=host " +
							"--privileged " +
							"--volumes-from=command-volumes " +
							"--volumes-from=system-volumes " +
							"stopuserdocker",
					},
					{
						Id: "k8s-docker",
						Cmd: "--name=k8s-docker " +
							"-d " +
							"--restart=always " +
							"--ipc=host " +
							"--pid=host " +
							"--net=host " +
							"--privileged " +
							"--volumes-from=command-volumes " +
							"--volumes-from=user-volumes " +
							"--volumes-from=system-volumes " +
							"-v=/var/lib/rancher/state/docker:/var/lib/docker " +
							"k8s-docker",
					},
				},
			},
			"k8s-master": {
				SystemContainers: []ContainerConfig{
					{
						Id: "kube-apiserver",
						Cmd: "--name=kube-apiserver " +
							"-d " +
							"--restart=always " +
							"--ipc=host " +
							"--pid=host " +
							"--net=host " +
							"--privileged " +
							"--volumes-from=command-volumes " +
							"--volumes-from=system-volumes " +
							"k8s " +
							"/kube-apiserver " +
							"--address=0.0.0.0 --port=8080 " +
							"--portal_net=10.100.0.0/16 " +
							"--etcd_servers=http://127.0.0.1:4001 " +
							"--public_address_override=172.17.7.101 " +
							"--logtostderr=true",
					},
					{
						Id: "kube-controller-manager",
						Cmd: "--name=kube-controller-manager " +
						"-d " +
						"--restart=always " +
						"--ipc=host " +
						"--pid=host " +
						"--net=host " +
						"--privileged " +
						"--volumes-from=command-volumes " +
						"--volumes-from=system-volumes " +
						"k8s " +
						"/kube-controller-manager " +
						"--master=127.0.0.1:8080 " +
						"--logtostderr=true",
					},
					{
						Id: "kube-scheduler",
						Cmd: "--name=kube-scheduler " +
						"-d " +
						"--restart=always " +
						"--ipc=host " +
						"--pid=host " +
						"--net=host " +
						"--privileged " +
						"--volumes-from=command-volumes " +
						"--volumes-from=system-volumes " +
						"k8s " +
						"/kube-scheduler " +
						"--master=127.0.0.1:8080 " +
						"--logtostderr=true",
					},
				},
			},
			"k8s-minion": {
				SystemContainers: []ContainerConfig{
					{
						Id: "kube-proxy",
						Cmd: "--name=kube-proxy " +
						"-d " +
						"--restart=always " +
						"--ipc=host " +
						"--pid=host " +
						"--net=host " +
						"--privileged " +
						"--volumes-from=command-volumes " +
						"--volumes-from=system-volumes " +
						"k8s " +
						"/kube-proxy " +
						"--etcd_servers=http://172.17.7.101:4001 " +
						"--logtostderr=true",
					},
					{
						Id: "kubelet",
						Cmd: "--name=kubelet " +
						"-d " +
						"--restart=always " +
						"--ipc=host " +
						"--pid=host " +
						"--net=host " +
						"--privileged " +
						"--volumes-from=command-volumes " +
						"--volumes-from=system-volumes " +
						"k8s " +
						"/kubelet " +
						"--address=0.0.0.0 --port=10250 " +
						"--hostname_override=$public_ipv4 " + //FIXME replace $public_ipv4 with the node IP
						"--api_servers=172.17.7.101:8080 " +
						"--logtostderr=true",
					},
				},
			},
		},
		RescueContainer: &ContainerConfig{
			Id: "console",
			Cmd: "--name=rescue " +
				"-d " +
				"--rm " +
				"--privileged " +
				"--volumes-from=console-volumes " +
				"--volumes-from=user-volumes " +
				"--volumes-from=system-volumes " +
				"--restart=always " +
				"--ipc=host " +
				"--net=host " +
				"--pid=host " +
				"rescue",
		},
	}
}
