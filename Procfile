
# Use goreman to run `go get github.com/mattn/goreman`
# Change the path of bin/etcd if etcd is located elsewhere

icsetcd1: etcd --name icsetcd1 --listen-client-urls http://127.0.0.1:12379 --advertise-client-urls http://127.0.0.1:12379 --listen-peer-urls http://127.0.0.1:12380 --initial-advertise-peer-urls http://127.0.0.1:12380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'icsetcd1=http://127.0.0.1:12380,icsetcd2=http://127.0.0.1:22380,icsetcd3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof
icsetcd2: etcd --name icsetcd2 --listen-client-urls http://127.0.0.1:22379 --advertise-client-urls http://127.0.0.1:22379 --listen-peer-urls http://127.0.0.1:22380 --initial-advertise-peer-urls http://127.0.0.1:22380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'icsetcd1=http://127.0.0.1:12380,icsetcd2=http://127.0.0.1:22380,icsetcd3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof
icsetcd3: etcd --name icsetcd3 --listen-client-urls http://127.0.0.1:32379 --advertise-client-urls http://127.0.0.1:32379 --listen-peer-urls http://127.0.0.1:32380 --initial-advertise-peer-urls http://127.0.0.1:32380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'icsetcd1=http://127.0.0.1:12380,icsetcd2=http://127.0.0.1:22380,icsetcd3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof
icsetcdproxy: etcd grpc-proxy start --endpoints=127.0.0.1:12379,127.0.0.1:22379,127.0.0.1:32379 --listen-addr=127.0.0.1:2379 --advertise-client-url=127.0.0.1:2379 --enable-pprof

# A learner node can be started using Procfile.learner