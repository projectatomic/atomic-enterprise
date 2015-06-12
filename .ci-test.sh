#!/bin/bash

# Pretend like we're Travis so that scripts behave accordingly.
export TRAVIS=true

# NB: set to false if we don't want to test assets
export TEST_ASSETS=true

echo "GOPATH = $GOPATH"
echo "PATH = $PATH"
go version
go env

cd $GOPATH/src/github.com/projectatomic/appinfra-next

# overall job status
failed=false

# run the passed command and capture output in a log
log_eval() {

	echo "RUNNING: $@" | tee -a $HOME/log
	eval $@ &>> $HOME/log

	if [ $? == 0 ]; then
		res=PASSED
	else
		res=FAILED
		failed=true
	fi

	echo "$res: $@" | tee -a $HOME/log
}

# installs

log_eval ./hack/verify-jsonformat.sh
log_eval ./hack/install-etcd.sh
log_eval ./hack/install-std-race.sh
log_eval ./hack/install-tools.sh
log_eval ./hack/build-go.sh
log_eval ./hack/install-assets.sh

# tests

# NB: because of the eval in log_eval(), we want to escape the quotes so that
# they don't reduce to nothing, which will cause that argument to be skipped
# when passed to test-go.sh as $(WHAT) (see the Makefile).
log_eval \
	PATH=./_output/etcd/bin:$PATH \
		make check-test WHAT="\'\'" TESTFLAGS="-p=4"

log_eval ./hack/test-assets.sh

if [ "$failed" = true ]; then
	exit 1
fi
