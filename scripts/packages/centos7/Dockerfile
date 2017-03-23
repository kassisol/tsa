FROM centos:7

COPY entrypoint.sh /entrypoint.sh
COPY build /usr/local/src/tsa

ENV DAPPER_SOURCE /tmp
ENV DAPPER_OUTPUT dist
ENV SHELL /bin/bash

WORKDIR ${DAPPER_SOURCE}

ENV RPMBUILD_PATH="/srv/rpmbuild"

RUN build="rpm-build" \
	&& set -x \
	&& yum -y install $build \
	&& yum -y autoremove \
	&& yum clean all

RUN mkdir -p ${RPMBUILD_PATH} \
	&& mkdir ${RPMBUILD_PATH}/BUILD \
	&& mkdir ${RPMBUILD_PATH}/RPMS \
	&& mkdir ${RPMBUILD_PATH}/SOURCES \
	&& mkdir ${RPMBUILD_PATH}/SPECS \
	&& mkdir ${RPMBUILD_PATH}/SRPMS \
	&& mkdir ${RPMBUILD_PATH}/tmp \
	&& echo "%_topdir  ${RPMBUILD_PATH}" > /root/.rpmmacros \
	&& echo "%_tmppath  ${RPMBUILD_PATH}/tmp" >> /root/.rpmmacros

COPY tsa.spec ${RPMBUILD_PATH}/SPECS/tsa.spec

RUN set -x \
	&& tar cvzf ${RPMBUILD_PATH}/SOURCES/tsa.tar.gz -C /usr/local/src tsa	 \
	&& rm -rf /usr/local/src/tsa

ENTRYPOINT ["/entrypoint.sh"]
