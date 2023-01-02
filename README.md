Zabbix agent plugin for monitoring APT package updates.

This plugin uses -s simulation option when invoking apt-get, so no root access is needed for Zabbix user during polling.

# Notes

To work properly, the list of packages needs to be updated periodically. For example, this action can be configured with APT::Periodic::Enable [https://wiki.debian.org/UnattendedUpgrades](UnattendedUpgrades).

# Installation

1. Download from assets plugin `zabbix-agent2-plugin-apt` to `/usr/sbin/zabbix-agent2-plugin/`
2. Download `apt.conf` to `/etc/zabbix/zabbix_agent2.d/plugins.d/` when installed by default and change preconfigured `Plugins.APT.System.Path=` if needed
