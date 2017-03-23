Name: tsa
Version: %{_version}
Release: %{_release}%{?dist}
Summary: HBM TSA is an application acting as a CA (Certificate Authority) for HBM TWIC
Group: Tools/Docker

License: GPL

URL: https://github.com/kassisol/tsa
Vendor: Kassisol
Packager: Kassisol <support@kassisol.com>

BuildArch: x86_64
BuildRoot: %{_tmppath}/%{name}-buildroot

Source: tsa.tar.gz

%description
HBM TSA is an application acting as a CA (Certificate Authority) for HBM TWIC.

%prep
%setup -n %{name}

%install
# install binary
install -d $RPM_BUILD_ROOT/%{_sbindir}
install -p -m 755 tsa $RPM_BUILD_ROOT/%{_sbindir}/

# add init scripts
install -d $RPM_BUILD_ROOT/%{_unitdir}
install -p -m 644 tsa.service $RPM_BUILD_ROOT/%{_unitdir}/tsa.service

# add bash completions
install -d $RPM_BUILD_ROOT/usr/share/bash-completion/completions
install -p -m 644 shellcompletion/bash $RPM_BUILD_ROOT/usr/share/bash-completion/completions/tsa

# install manpages
install -d $RPM_BUILD_ROOT/%{_mandir}/man8
install -p -m 644 man/man8/*.8 $RPM_BUILD_ROOT/%{_mandir}/man8

# list files owned by the package here
%files
#%doc README.md
%{_sbindir}/tsa
/%{_unitdir}/tsa.service
/usr/share/bash-completion/completions/tsa
%doc
/%{_mandir}/man8/*

%post
%systemd_post tsa

%preun
%systemd_preun tsa

%postun
rm -f %{_sbindir}/tsa

%clean
rm -rf $RPM_BUILD_ROOT
