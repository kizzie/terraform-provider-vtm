package main

import (
	"fmt"
	"github.com/atlassian/go-vtm"
	"github.com/hashicorp/terraform/helper/schema"
)

func schemaStringList() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func resourceEventType() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventTypeCreate,
		Read:   resourceEventTypeRead,
		Update: resourceEventTypeUpdate,
		Delete: resourceEventTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"actions": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"built_in": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloudcredentials_event_tags": schemaStringList(),
			"cloudcredentials_objects":    schemaStringList(),
			"config_event_tags":           schemaStringList(),
			"faulttolerance_event_tags":   schemaStringList(),
			"general_event_tags":          schemaStringList(),
			"glb_event_tags":              schemaStringList(),
			"glb_objects":                 schemaStringList(),
			"java_event_tags":             schemaStringList(),
			"licensekeys_event_tags":      schemaStringList(),
			"licensekeys_objects":         schemaStringList(),
			"locations_event_tags":        schemaStringList(),
			"locations_objects":           schemaStringList(),
			"monitors_event_tags":         schemaStringList(),
			"monitors_objects":            schemaStringList(),
			"pools_event_tags":            schemaStringList(),
			"pools_objects":               schemaStringList(),
			"protection_event_tags":       schemaStringList(),
			"protection_objects":          schemaStringList(),
			"rules_event_tags":            schemaStringList(),
			"rules_objects":               schemaStringList(),
			"slm_event_tags":              schemaStringList(),
			"slm_objects":                 schemaStringList(),
			"ssl_event_tags":              schemaStringList(),
			"sslhw_event_tags":            schemaStringList(),
			"trafficscript_event_tags":    schemaStringList(),
			"vservers_event_tags":         schemaStringList(),
			"vservers_objects":            schemaStringList(),
			"zxtms_event_tags":            schemaStringList(),
			"zxtms_objects":               schemaStringList(),
		},
	}
}

func resourceEventTypeSet(d *schema.ResourceData, m interface{}) error {
	c := m.(*providerConfig).client
	r := stingray.NewEventType(d.Get("name").(string))

	setString(&r.Basic.Note, d, "note")
	setStringSet(&r.Basic.Actions, d, "actions")

	setStringSet(&r.CloudCredentials.EventTags, d, "cloudcredentials_event_tags")
	setStringSet(&r.CloudCredentials.Objects, d, "cloudcredentials_objects")

	setStringSet(&r.Config.EventTags, d, "config_event_tags")

	setStringSet(&r.FaultTolerance.EventTags, d, "faulttolerance_event_tags")

	setStringSet(&r.General.EventTags, d, "general_event_tags")

	setStringSet(&r.GLB.EventTags, d, "glb_event_tags")
	setStringSet(&r.GLB.Objects, d, "glb_objects")

	setStringSet(&r.Java.EventTags, d, "java_event_tags")

	setStringSet(&r.LicenseKeys.EventTags, d, "licensekeys_event_tags")
	setStringSet(&r.LicenseKeys.Objects, d, "licensekeys_objects")

	setStringSet(&r.Locations.EventTags, d, "locations_event_tags")
	setStringSet(&r.Locations.Objects, d, "locations_objects")

	setStringSet(&r.Monitors.EventTags, d, "monitors_event_tags")
	setStringSet(&r.Monitors.Objects, d, "monitors_objects")

	setStringSet(&r.Pools.EventTags, d, "pools_event_tags")
	setStringSet(&r.Pools.Objects, d, "pools_objects")

	setStringSet(&r.Protection.EventTags, d, "protection_event_tags")
	setStringSet(&r.Protection.Objects, d, "protection_objects")

	setStringSet(&r.Rules.EventTags, d, "rules_event_tags")
	setStringSet(&r.Rules.Objects, d, "rules_objects")

	setStringSet(&r.Slm.EventTags, d, "slm_event_tags")
	setStringSet(&r.Slm.Objects, d, "slm_objects")

	setStringSet(&r.SSL.EventTags, d, "ssl_event_tags")

	setStringSet(&r.SSLhw.EventTags, d, "sslhw_event_tags")

	setStringSet(&r.TrafficScript.EventTags, d, "trafficscript_event_tags")

	setStringSet(&r.Vservers.EventTags, d, "vservers_event_tags")
	setStringSet(&r.Vservers.Objects, d, "vservers_objects")

	setStringSet(&r.Zxtms.EventTags, d, "zxtms_event_tags")
	setStringSet(&r.Zxtms.Objects, d, "zxtms_objects")

	_, err := c.Set(r)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func resourceEventTypeCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceEventTypeSet(d, m)
	if err != nil {
		return err
	}

	return resourceEventTypeRead(d, m)
}
func resourceEventTypeUpdate(d *schema.ResourceData, m interface{}) error {
	err := resourceEventTypeSet(d, m)
	if err != nil {
		return err
	}

	return resourceEventTypeRead(d, m)
}
func resourceEventTypeDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*providerConfig).client
	r := stingray.NewEventType(d.Id())

	_, err := c.Delete(r)
	if err != nil {
		return err
	}

	return nil
}
func resourceEventTypeRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*providerConfig).client

	r, resp, err := c.GetEventType(d.Get("name").(string))
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// The resource doesn't exist anymore
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading resource: %s", err)
	}
	d.Set("actions", []string(*r.Basic.Actions))
	d.Set("note", string(*r.Basic.Note))
	d.Set("built_in", bool(*r.Basic.BuiltIn))

	d.Set("cloudcredentials_event_tags", []string(*r.CloudCredentials.EventTags))
	d.Set("cloudcredentials_objects", []string(*r.CloudCredentials.Objects))

	d.Set("config_event_tags", []string(*r.Config.EventTags))

	d.Set("faulttolerance_event_tags", []string(*r.FaultTolerance.EventTags))

	d.Set("general_event_tags", []string(*r.General.EventTags))

	d.Set("glb_event_tags", []string(*r.GLB.EventTags))
	d.Set("glb_objects", []string(*r.GLB.Objects))

	d.Set("java_event_tags", []string(*r.Java.EventTags))

	d.Set("licensekeys_event_tags", []string(*r.LicenseKeys.EventTags))
	d.Set("licensekeys_objects", []string(*r.LicenseKeys.Objects))

	d.Set("locations_event_tags", []string(*r.Locations.EventTags))
	d.Set("locations_objects", []string(*r.Locations.Objects))

	d.Set("monitors_event_tags", []string(*r.Monitors.EventTags))
	d.Set("monitors_objects", []string(*r.Monitors.Objects))

	d.Set("pools_event_tags", []string(*r.Pools.EventTags))
	d.Set("pools_objects", []string(*r.Pools.Objects))

	d.Set("protection_event_tags", []string(*r.Protection.EventTags))
	d.Set("protection_objects", []string(*r.Protection.Objects))

	d.Set("rules_event_tags", []string(*r.Rules.EventTags))
	d.Set("rules_objects", []string(*r.Rules.Objects))

	d.Set("slm_event_tags", []string(*r.Slm.EventTags))
	d.Set("slm_objects", []string(*r.Slm.Objects))

	d.Set("ssl_event_tags", []string(*r.SSL.EventTags))

	d.Set("sslhw_event_tags", []string(*r.SSLhw.EventTags))

	d.Set("trafficscript_event_tags", []string(*r.TrafficScript.EventTags))

	d.Set("vserservers_event_tags", []string(*r.Vservers.EventTags))
	d.Set("vserservervs_objects", []string(*r.Vservers.Objects))

	d.Set("zxtms_event_tags", []string(*r.Zxtms.EventTags))
	d.Set("zxtms_objects", []string(*r.Zxtms.Objects))

	return nil
}
