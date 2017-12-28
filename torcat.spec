%global	import_path		github.com/iavael/torcat
# https://bugzilla.redhat.com/show_bug.cgi?id=995136#c12
%global	_dwz_low_mem_die_limit	0

Name:		torcat
Version:	0.1
Release:	1%{?dist}
Summary:	Netcat for onion router

License:	ASL 2.0
URL:		https://%{import_path}
Source0:	https://%{import_path}/archive/%{name}-%{version}.tar.gz

ExclusiveArch:	%{?go_arches:%{go_arches}}%{!?go_arches:%{ix86} x86_64 %{arm}}

BuildRequires:	%{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}

BuildRequires:	golang(github.com/yawning/bulb)
BuildRequires:	golang(github.com/yawning/bulb/utils)

%description
Netcat-like tool for tor network

%prep
%autosetup

%build
%{__mkdir_p} gopath/src/%{import_path}
%{__rm} -d gopath/src/%{import_path}
%{__ln_s} $(pwd) gopath/src/%{import_path}

%if ! 0%{?with_bundled}
export GOPATH=$(pwd)/gopath:%{gopath}
%else
echo "Unable to build from bundled deps. No Godeps nor vendor directory"
exit 1
%endif

%gobuild -o gopath/bin/%{name} %{import_path}

%install
rm -rf %{buildroot}
%{__install} -d %{buildroot}/%{_bindir}
%{__install} -p -m 755 gopath/bin/%{name} %{buildroot}%{_bindir}

%files
%license COPYING
%doc README.md
%{_bindir}/%{name}

%changelog
