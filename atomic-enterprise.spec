#debuginfo not supported with Go
%global debug_package %{nil}
%global gopath      %{_datadir}/gocode
%global import_path github.com/projectatomic/appinfra-next
%global kube_plugin_path /usr/libexec/kubernetes/kubelet-plugins/net/exec/redhat~openshift-ovs-subnet
%global sdn_import_path github.com/openshift/openshift-sdn

# %%commit and %%ldflags are intended to be set by tito custom builders provided
# in the rel-eng directory. The values in this spec file will not be kept up to date.
%{!?commit:
%global commit 54e7bfc9b4765a22ddb4c9c8a0c37c42eeab0dbd
}
%global shortcommit %(c=%{commit}; echo ${c:0:7})
# OpenShift AE specific ldflags from hack/common.sh os::build:ldflags
%{!?ldflags:
%global ldflags -X github.com/projectatomic/appinfra-next/pkg/version.majorFromGit 0 -X github.com/projectatomic/appinfra-next/pkg/version.minorFromGit 0+ -X github.com/projectatomic/appinfra-next/pkg/version.versionFromGit v0.0.1 -X github.com/projectatomic/appinfra-next/pkg/version.commitFromGit e3c46fd -X github.com/GoogleCloudPlatform/kubernetes/pkg/version.gitCommit 496be63 -X github.com/GoogleCloudPlatform/kubernetes/pkg/version.gitVersion v0.17.1-804-g496be63
}

Name:           atomic-enterprise
# Version is not kept up to date and is intended to be set by tito custom
# builders provided in the rel-eng directory of this project
Version:        0.0.1
Release:        1.git.0.097952c%{?dist}
Summary:        Open Source Platform as a Service by Red Hat
License:        ASL 2.0
URL:            https://%{import_path}
ExclusiveArch:  x86_64
Source0:        %{name}-git-0.%{shortcommit}.tar.gz

BuildRequires:  systemd
BuildRequires:  golang >= 1.4


%description
%{summary}

%package master
Summary:        Atomic Enterprise Master
Requires:       %{name} = %{version}-%{release}
Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd

%description master
%{summary}

%package node
Summary:        Atomic Enterprise Node
Requires:       %{name} = %{version}-%{release}
Requires:       docker >= 1.6.2
Requires:       tuned-profiles-%{name}-node
Requires:       util-linux
Requires:       socat
Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd

%description node
%{summary}

%package -n tuned-profiles-%{name}-node
Summary:        Tuned profiles for OpenShift AE Node hosts
Requires:       tuned >= 2.3
Requires:       %{name} = %{version}-%{release}

%description -n tuned-profiles-%{name}-node
%{summary}

%package clients
Summary:      Atomic Enterprise Client binaries for Linux, Mac OSX, and Windows
BuildRequires: golang-pkg-darwin-amd64
BuildRequires: golang-pkg-windows-386

%description clients
%{summary}

%package dockerregistry
Summary:        Docker Registry v2 for Atomic Enterprise
Requires:       %{name} = %{version}-%{release}

%description dockerregistry
%{summary}

%package pod
Summary:        Atomic Enterprise Pod
Requires:       %{name} = %{version}-%{release}

%description pod
%{summary}

%package sdn-ovs
Summary:          Atomic Enterprise SDN Plugin for Open vSwitch
Requires:         openvswitch >= 2.3.1
Requires:         %{name}-node = %{version}-%{release}
Requires:         bridge-utils
Requires:         ethtool

%description sdn-ovs
%{summary}

%prep
%setup -q -n %{name}-git-0.%{shortcommit}

%build

# Don't judge me for this ... it's so bad.
mkdir _build

# Horrid hack because golang loves to just bundle everything
pushd _build
    mkdir -p src/github.com/projectatomic
    ln -s $(dirs +1 -l) src/%{import_path}
popd


# Gaming the GOPATH to include the third party bundled libs at build
# time. This is bad and I feel bad.
mkdir _thirdpartyhacks
pushd _thirdpartyhacks
    ln -s \
        $(dirs +1 -l)/Godeps/_workspace/src/ \
            src
popd
export GOPATH=$(pwd)/_build:$(pwd)/_thirdpartyhacks:%{buildroot}%{gopath}:%{gopath}
# Build all linux components we care about
for cmd in openshift dockerregistry
do
        go install -ldflags "%{ldflags}" %{import_path}/cmd/${cmd}
done


# TODO: Do we need to do this?
# Build only 'openshift' for other platforms
#GOOS=windows GOARCH=386 go install -ldflags "%{ldflags}" %{import_path}/cmd/openshift
#GOOS=darwin GOARCH=amd64 go install -ldflags "%{ldflags}" %{import_path}/cmd/openshift

#Build our pod
pushd images/pod/
    go build -ldflags "%{ldflags}" pod.go
popd

%install

install -d %{buildroot}%{_bindir}
install -d %{buildroot}%{_datadir}/%{name}/linux
#{linux,macosx,windows}

mv %{buildroot}%{_datadir}/%{name} %{buildroot}%{_datadir}/openshift

# Install linux components
echo "+++ INSTALLING atomic-enterprise"
install -p -m 755 _build/bin/openshift %{buildroot}%{_bindir}/%{name}
echo "+++ INSTALLING dockerregistry"
install -p -m 755 _build/bin/dockerregistry %{buildroot}%{_bindir}/dockerregistry

#for bin in openshift dockerregistry
#do
#  install -p -m 755 _build/bin/${bin} %{buildroot}%{_bindir}/${bin}
#done

# Install 'openshift' as client executable for windows and mac
install -p -m 755 _build/bin/openshift %{buildroot}%{_datadir}/openshift/linux/oc
#install -p -m 755 _build/bin/darwin_amd64/openshift %{buildroot}%{_datadir}/openshift/macosx/oc
#install -p -m 755 _build/bin/windows_386/openshift.exe %{buildroot}%{_datadir}/openshift/windows/oc.exe
#Install openshift pod
install -p -m 755 images/pod/pod %{buildroot}%{_bindir}/

install -d -m 0755 %{buildroot}/etc/%{name}/{master,node}
mv %{buildroot}/etc/%{name} %{buildroot}/etc/openshift
install -d -m 0755 %{buildroot}%{_unitdir}
install -m 0644 -t %{buildroot}%{_unitdir} rel-eng/%{name}-master.service
install -m 0644 -t %{buildroot}%{_unitdir} rel-eng/%{name}-node.service

mkdir -p %{buildroot}%{_sysconfdir}/sysconfig
install -m 0644 rel-eng/openshift-master.sysconfig %{buildroot}%{_sysconfdir}/sysconfig/%{name}-master
install -m 0644 rel-eng/openshift-node.sysconfig %{buildroot}%{_sysconfdir}/sysconfig/%{name}-node

mkdir -p %{buildroot}%{_sharedstatedir}/%{name}

ln -s %{_bindir}/%{name} %{buildroot}%{_bindir}/oc
ln -s %{_bindir}/%{name} %{buildroot}%{_bindir}/oadm

install -d -m 0755 %{buildroot}%{_prefix}/lib/tuned/%{name}-node-{guest,host}
install -m 0644 tuned/openshift-node-guest/tuned.conf %{buildroot}%{_prefix}/lib/tuned/%{name}-node-guest/
install -m 0644 tuned/openshift-node-host/tuned.conf %{buildroot}%{_prefix}/lib/tuned/%{name}-node-host/
install -d -m 0755 %{buildroot}%{_mandir}/man7
install -m 0644 tuned/man/tuned-profiles-openshift-node.7 %{buildroot}%{_mandir}/man7/tuned-profiles-%{name}-node.7

# Install sdn scripts
install -d -m 0755 %{buildroot}%{kube_plugin_path}
pushd _thirdpartyhacks/src/%{sdn_import_path}/ovssubnet/bin
   install -p -m 755 openshift-ovs-subnet %{buildroot}%{kube_plugin_path}/openshift-ovs-subnet
   install -p -m 755 openshift-sdn-kube-subnet-setup.sh %{buildroot}%{_bindir}/
popd
install -d -m 0755 %{buildroot}%{_prefix}/lib/systemd/system/%{name}-node.service.d
install -p -m 0644 rel-eng/openshift-sdn-ovs.conf %{buildroot}%{_prefix}/lib/systemd/system/%{name}-node.service.d/%{name}-sdn-ovs.conf
install -d -m 0755 %{buildroot}%{_prefix}/lib/systemd/system/docker.service.d
install -p -m 0644 rel-eng/docker-sdn-ovs.conf %{buildroot}%{_prefix}/lib/systemd/system/docker.service.d/

# Install bash completions
install -d -m 755 %{buildroot}/etc/bash_completion.d/
install -p -m 644 rel-eng/completions/bash/* %{buildroot}/etc/bash_completion.d/

%files
%defattr(-,root,root,-)
%doc README.md LICENSE
%{_bindir}/%{name}
%{_bindir}/oc
%{_bindir}/oadm
%{_sharedstatedir}/%{name}
/etc/bash_completion.d/*

%files master
%defattr(-,root,root,-)
%{_unitdir}/%{name}-master.service
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}-master
%config(noreplace) /etc/openshift/master

%post master
%systemd_post %{basename:%{name}-master.service}

%preun master
%systemd_preun %{basename:%{name}-master.service}

%postun master
%systemd_postun


%files node
%defattr(-,root,root,-)
%{_unitdir}/%{name}-node.service
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}-node
%config(noreplace) /etc/openshift/node

%post node
%systemd_post %{basename:%{name}-node.service}

%preun node
%systemd_preun %{basename:%{name}-node.service}

%postun node
%systemd_postun

%files sdn-ovs
%defattr(-,root,root,-)
%{_bindir}/openshift-sdn-kube-subnet-setup.sh
%{kube_plugin_path}/openshift-ovs-subnet
%{_prefix}/lib/systemd/system/%{name}-node.service.d/%{name}-sdn-ovs.conf
%{_prefix}/lib/systemd/system/docker.service.d/docker-sdn-ovs.conf

%files -n tuned-profiles-%{name}-node
%defattr(-,root,root,-)
%{_prefix}/lib/tuned/%{name}-node-host
%{_prefix}/lib/tuned/%{name}-node-guest
%{_mandir}/man7/tuned-profiles-%{name}-node.7*

%post -n tuned-profiles-%{name}-node
recommended=`/usr/sbin/tuned-adm recommend`
if [[ "${recommended}" =~ guest ]] ; then
  /usr/sbin/tuned-adm profile %{name}-node-guest > /dev/null 2>&1
else
  /usr/sbin/tuned-adm profile %{name}-node-host > /dev/null 2>&1
fi

%preun -n tuned-profiles-%{name}-node
# reset the tuned profile to the recommended profile
# $1 = 0 when we're being removed > 0 during upgrades
if [ "$1" = 0 ]; then
  recommended=`/usr/sbin/tuned-adm recommend`
  /usr/sbin/tuned-adm profile $recommended > /dev/null 2>&1
fi

%files clients
%{_datadir}/openshift/linux/oc
#%{_datadir}/openshift/macosx/oc
#%{_datadir}/openshift/windows/oc.exe

%files dockerregistry
%defattr(-,root,root,-)
%{_bindir}/dockerregistry

%files pod
%defattr(-,root,root,-)
%{_bindir}/pod

%changelog
