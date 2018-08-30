# haproxy-kube-agent

haproxy-kube-agent will work together with [cloud-provider-baremetal](https://github.com/chongzii6/cloud-provider-baremetal)

haproxy-kube-agent publish status of haproxy in localhost, including IP, listen section, ports.
cloud-provider-baremetal will consume eht status and sned add/del request of lb, they communicate via etcd

# The working procedure is:
1. haproxy-kube-agent publish state to etcd, key=/chongzii6/lb/<lb_name>, value=[<lb_ip>:<listen_port>]. currently we only use haproxy listen configuration
2. cloud-provider-baremetal watch key /chongzii6/lb/ for lb status
3. cloud-provider-baremetal publish request to etcd, key=/chongzii6/reqlb/(add|del)/<lb_name>, value=[<endpoints>]
4. haproxy-kube-agent watch key /chongzii6/reqlb for request
5. multiple haproxy-kube-agent use distributed lock which supported by etcd clientv3.concurrency to acquire right of handling /chongzii6/reqlb
6. if haproxy-kube-agent is able to fulfill the request, it will delete /chongzii6/reqlb/(add|del)/<lb_name>
7. whatever the result haproxy-kube-agent can fulfill, it should release the lock, let next one to handle request
8. cloud-provider-baremetal will get result after haproxy-kube-agent update etcd key=/chongzii6/lb/<lb_name>, 