Name:           autobrowser
Version:        1.0
Release:        1%{?dist}
Summary:        Automatically chose browser depends on context

License:        GPLv3
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang

Provides:       %{name} = %{version}

%description
Automatically chose which browser to opend depends on context

%prep
%autosetup


%build
make build-linux


%install
install -Dpm 0755 out/autobrowser %{buildroot}%{_bindir}/%{name}


%files
%{_bindir}/autobrowser
