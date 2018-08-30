# haproxy-kube-agent

haproxy-kube-agent will work together with [cloud-provider-baremetal](https://github.com/chongzii6/cloud-provider-baremetal)

haproxy-kube-agent publish status of haproxy in localhost, including IP, listen section, ports.
cloud-provider-baremetal will consume eht status and sned add/del request of lb, they communicate via etcd

# The working procedure is:
1. haproxy-kube-agent publish state to etcd, key=`<env:ETCD_LBAGENT_KEY>`/<lb_name>, value=[<lb_ip>:<listen_port>]. currently we only use haproxy listen configuration
2. cloud-provider-baremetal watch key `<env:ETCD_LBAGENT_KEY>` for lb status
3. cloud-provider-baremetal publish request to etcd, key=`<env:ETCD_LBREQ_KEY>`/(add|del)/<lb_name>, value=[<endpoints>]
4. haproxy-kube-agent watch key `<env:ETCD_LBREQ_KEY>` for request
5. multiple haproxy-kube-agent use distributed lock which supported by etcd clientv3.concurrency to acquire right of handling `<env:ETCD_LBREQ_KEY>`
6. if haproxy-kube-agent is able to fulfill the request, it will delete `<env:ETCD_LBREQ_KEY>`/(add|del)/<lb_name>
7. whatever the result haproxy-kube-agent can fulfill, it should release the lock, let next one to handle request
8. cloud-provider-baremetal will get result after haproxy-kube-agent update etcd key=`<env:ETCD_LBAGENT_KEY>`/<lb_name>, 

# Required Envs
* ETCD_LBAGENT_KEY, default=`/chongzii6/lb`
* ETCD_LBREQ_KEY, default=`/chongzii6/reqlb`
* ETCD_CA_FILE, path to ca.pem, default=`ca.pem`
* ETCD_CERT_FILE, path to cert.pem, default=`cert.pem`
* ETCD_KEY_FILE, path to key.pem, default=`key.pem`
