Name: tsa
Version: %{_version}
Release: %{_release}%{?dist}
Summary: Starts a CA (Certification Authority) server
Group: Tools/Docker

License: GPLv3+

URL: https://github.com/kassisol/tsa
Vendor: Kassisol
Packager: Kassisol <support@kassisol.com>

BuildArch: x86_64
BuildRoot: %{_tmppath}/%{name}-buildroot

Source: tsa.tar.gz

%description
HBM TSA is an application acting as a CA (Certification Authority) server to issue certificates for Docker Engine with TLS enabled

%package daemon
Summary: The TSA daemon
Group: Tools/Docker
%description daemon
Starts a CA (Certification Authority) server

%package adm
Summary: The TSA Client
Group: Tools/Docker
%description adm
Client to manage TSA daemon

%prep
%setup -n %{name}

%install
# install binary
install -d $RPM_BUILD_ROOT%{_sbindir}
install -d $RPM_BUILD_ROOT%{_bindir}
install -p -m 755 bin/tsad $RPM_BUILD_ROOT%{_sbindir}/
install -p -m 755 bin/tsa $RPM_BUILD_ROOT%{_bindir}/

# add init scripts
install -d $RPM_BUILD_ROOT/%{_unitdir}
install -p -m 644 tsad.service $RPM_BUILD_ROOT/%{_unitdir}/tsad.service

# add bash completions
install -d $RPM_BUILD_ROOT/usr/share/bash-completion/completions
install -p -m 644 tsad/shellcompletion/bash $RPM_BUILD_ROOT/usr/share/bash-completion/completions/tsad
install -p -m 644 tsa/shellcompletion/bash $RPM_BUILD_ROOT/usr/share/bash-completion/completions/tsa

# install manpages
install -d $RPM_BUILD_ROOT%{_mandir}/man8
install -p -m 644 tsad/man/man8/*.8 $RPM_BUILD_ROOT%{_mandir}/man8/
install -p -m 644 tsa/man/man8/*.8 $RPM_BUILD_ROOT%{_mandir}/man8/

ls -lR $RPM_BUILD_ROOT%{_mandir}/man8/

%files daemon
%{_sbindir}/tsad
/%{_unitdir}/tsad.service
/usr/share/bash-completion/completions/tsad
%{_mandir}/man8/tsad.8

%files adm
%{_bindir}/tsa
/usr/share/bash-completion/completions/tsa
%{_mandir}/man8/tsa-*
%{_mandir}/man8/tsa.8

%post daemon
%systemd_post tsad

%preun daemon
%systemd_preun tsad

%postun daemon
rm -f %{_sbindir}/tsad
rm -f /usr/share/bash-completion/completions/tsad
rm -f %{_mandir}/man8

%postun adm
rm -f %{_bindir}/tsa

%clean
rm -rf $RPM_BUILD_ROOT
