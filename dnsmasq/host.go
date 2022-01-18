package dnsmasq

import (
	"html/template"
	"os"
	"path"

	"dips/models"
)

func createHostConfigs(h *models.Host) error {
	hostsConfigPath := path.Join(conf.HostsDir, h.Hostname+"."+h.Domain)
	hostsFile, err := os.OpenFile(hostsConfigPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer hostsFile.Close()

	hostsTemplateString := "{{ .IP }} {{ .Hostname }}.{{.Domain}}"
	hostsTemplate, err := template.New("hosts").Parse(hostsTemplateString)
	if err != nil {
		return err
	}

	err = hostsTemplate.Execute(hostsFile, h)
	if err != nil {
		return err
	}

	dhcpHostsConfigPath := path.Join(conf.DhcpHostsDir, h.Hostname+"."+h.Domain)
	dhcpHostsFile, err := os.OpenFile(dhcpHostsConfigPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer dhcpHostsFile.Close()

	dhcpHostsTemplateString := "{{ .MAC }},set:{{ .Domain }},{{ .IP }},{{ .Hostname }}.{{.Domain}}"
	dhcpHostsTemplate, err := template.New("dhcp_hosts").Parse(dhcpHostsTemplateString)
	if err != nil {
		return err
	}

	err = dhcpHostsTemplate.Execute(dhcpHostsFile, h)
	if err != nil {
		return err
	}
	return nil
}

func deleteHostConfigs(h *models.Host) error {
	hostsConfigPath := path.Join(conf.HostsDir, h.Hostname+"."+h.Domain)
	if _, err := os.Stat(hostsConfigPath); os.IsNotExist(err) {
		return nil // already gone
	}
	err := os.Remove(hostsConfigPath)
	if err != nil {
		return err
	}
	dhcpHostsConfigPath := path.Join(conf.DhcpHostsDir, h.Hostname+"."+h.Domain)
	if _, err := os.Stat(dhcpHostsConfigPath); os.IsNotExist(err) {
		return nil // already gone
	}
	err = os.Remove(dhcpHostsConfigPath)
	if err != nil {
		return err
	}
	return nil
}
