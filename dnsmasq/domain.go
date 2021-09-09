package dnsmasq

import (
	"html/template"
	"os"
	"path"

	"dips.alta3.com/models"
	"github.com/lithammer/dedent"
)

func CreateDomainConfig(h models.Host) error {
	domainConfigPath := path.Join(dnsmasq.dhcpOptionsDir, h.Domain)

	f, err := os.OpenFile(domainConfigPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	DomainTemplate := dedent.Dedent(`
                              # Managed by DIPS
                              tag:{{ .Domain }},3,{{ .Gateway }}
                              tag:{{ .Domain }},119,{{ .Domain }} 
                              tag:{{ .Domain }},15,{{ .Domain }}`)
	tmpl, err := template.New("domain").Parse(DomainTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, h)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDomainConfig(h models.Host) error {
	domainConfigPath := path.Join(dnsmasq.dhcpOptionsDir, h.Domain)
	if _, err := os.Stat(domainConfigPath); os.IsNotExist(err) {
		return nil // already gone
	}
	err := os.Remove(domainConfigPath)
	if err != nil {
		return err
	}
	return nil
}
