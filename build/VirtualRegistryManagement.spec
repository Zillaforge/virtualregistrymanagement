%define RPM_BUILD_ROOT _topdir
%define Name VirtualRegistryManagement
%define Version 0.0.6
%define Release 1
%define BuildRoot %{_topdir}/BUILDROOT/%{Name}-%{Version}-%{Release}.x86_64
%define Group pegasus
%define User pegasus
%define HomeDir /var/lib/ASUS/
%define Comment "pegasus cloud team"

Name:           %{Name}
Version:        %{Version}
License:        GPL
Release:        %{Release}
Summary:        ASUS VirtualRegistryManagement package
Url:            asus.com
Group:          developer
Source:         %{Name}_%{Version}.tar.gz
BuildRoot:      %BuildRoot

%pre
getent group %{Group} > /dev/null || groupadd -r %{Group}
getent passwd %{User} > /dev/null || useradd -r -g %{Group} -d %{HomeDir} -s /sbin/nologin -c %{Comment} %{User}

%description

%prep

%setup -q

%build

%install
install -d %{BuildRoot}/usr/bin/
install -d %{BuildRoot}/etc/ASUS/
install -d %{BuildRoot}/var/lib/ASUS
install -d %{BuildRoot}/var/log/ASUS
install -d %{BuildRoot}/etc/systemd/system/
cp %{Name}  %{BuildRoot}/usr/bin/
cp VirtualRegistryManagement.service %{BuildRoot}/etc/systemd/system/

%post
mkdir -p /etc/ASUS
mkdir -p /var/lib/ASUS
mkdir -p /var/log/ASUS

%preun

%files
%attr(755,root,root) /usr/bin/%{Name}
%attr(755,pegasus,pegasus) /etc/ASUS
%attr(755,pegasus,pegasus) /var/lib/ASUS
%attr(755,pegasus,pegasus) /var/log/ASUS
%config /etc/systemd/system/VirtualRegistryManagement.service

%clean
rm -rf %BuildRoot/
