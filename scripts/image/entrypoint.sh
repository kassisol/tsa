#!/bin/bash

TSA_LIB="/var/lib/tsa"
TSA="/usr/local/sbin/tsa"

function report_error() {
	local var=$1

	echo "var \"${var}\" is not set"
	exit 1
}

if [ -z $TSA_COUNTRY ]; then
	report_error TSA_COUNTRY
fi

if [ -z $TSA_STATE ]; then
	report_error TSA_STATE
fi

if [ -z $TSA_CITY ]; then
	report_error TSA_CITY
fi

if [ -z "$TSA_ORG" ]; then
	report_error TSA_ORG
fi

if [ -z $TSA_ORG_UNIT ]; then
	report_error TSA_ORG_UNIT
fi

if [ -z $TSA_EMAIL ]; then
	report_error TSA_EMAIL
fi

if [ -z $TSA_API_FQDN ]; then
	report_error TSA_API_FQDN
fi

if [ ! -z $TSA_AUTH_TYPE ]; then
	if [ ${TSA_AUTH_TYPE} == "ldap" ]; then
		if [ -z $TSA_AUTH_HOST ]; then
			report_error TSA_AUTH_HOST
		fi

		if [ -z $TSA_AUTH_PORT ]; then
			report_error TSA_AUTH_PORT
		fi

		if [ -z $TSA_AUTH_TLS ]; then
			report_error TSA_AUTH_TLS
		fi

		if [ -z $TSA_AUTH_BIND_USERNAME ]; then
			report_error TSA_AUTH_BIND_USERNAME
		fi

		if [ -z $TSA_AUTH_BIND_PASSWORD ]; then
			report_error TSA_AUTH_BIND_PASSWORD
		fi

		if [ -z $TSA_AUTH_SEARCH_BASE_USER ]; then
			report_error TSA_AUTH_SEARCH_BASE_USER
		fi

		if [ -z $TSA_AUTH_SEARCH_FILTER ]; then
			report_error TSA_AUTH_SEARCH_FILTER
		fi

		if [ -z $TSA_AUTH_ATTR_MEMBERS ]; then
			report_error TSA_AUTH_ATTR_MEMBERS
		fi

		if [ -z $TSA_AUTH_GROUP_ADMIN ]; then
			report_error TSA_AUTH_GROUP_ADMIN
		fi

		if [ -z $TSA_AUTH_GROUP_USER ]; then
			report_error TSA_AUTH_GROUP_USER
		fi
	fi
fi

if [ ! -d ${TSA_LIB} ]; then
        echo "${TSA_LIB} is not mounted"
        exit 1
fi

if [ `ls -1 ${TSA_LIB} | wc -l` -eq 0 ]; then
        ${TSA} init \
                --country ${TSA_COUNTRY} \
                --state ${TSA_STATE} \
                --city ${TSA_CITY} \
                --org ${TSA_ORG} \
                --org-unit ${TSA_ORG_UNIT} \
                --email ${TSA_EMAIL} \
                --api-fqdn ${TSA_API_FQDN}
fi

if [ `${TSA} info | awk -F ': ' '/Auth Type/ { print $2 }'` != ${TSA_AUTH_TYPE} ]; then
	${TSA} auth enable ${TSA_AUTH_TYPE}

	if [ ${TSA_AUTH_TYPE} == "ldap" ]; then
		${TSA} auth add auth_host ${TSA_AUTH_HOST}
		${TSA} auth add auth_port ${TSA_AUTH_PORT}
		${TSA} auth add auth_tls ${TSA_AUTH_TLS}
		${TSA} auth add auth_bind_username ${TSA_AUTH_BIND_USERNAME}
		${TSA} auth add auth_bind_password ${TSA_AUTH_BIND_PASSWORD}
		${TSA} auth add auth_search_base_user ${TSA_AUTH_SEARCH_BASE_USER}
		${TSA} auth add auth_search_filter ${TSA_AUTH_SEARCH_FILTER}
		${TSA} auth add auth_attr_members ${TSA_AUTH_ATTR_MEMBERS}
		${TSA} auth add auth_group_admin ${TSA_AUTH_GROUP_ADMIN}
		${TSA} auth add auth_group_user ${TSA_AUTH_GROUP_USER}
	fi

fi

exec ${TSA} server start
