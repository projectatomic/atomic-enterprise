Hello, Atomic!
-----------------

This example will serve an http response of "Hello Atomic!" to [http://localhost:6061](http://localhost:6061).  To create the pod run:

        $ oc create -f examples/hello-atomic/hello-pod.json

If you need to rebuild the image:
$ go build -tags netgo   # ensures static binary
$ mv hello-atomic bin
$ docker build -t docker.io/atomic-enterprise/hello-atomic .
