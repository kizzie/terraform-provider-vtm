package main

import (
    "fmt"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/cwood/go-vtm"
)

func resourceProtection() *schema.Resource {
    return &schema.Resource {
        Create: resourceProtectionCreate,
        Read: resourceProtectionRead,
        Update: resourceProtectionUpdate,
        Delete: resourceProtectionDelete,

        Schema: map[string]*schema.Schema{
            "name": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
            },
            "debug": &schema.Schema{
                Type: schema.TypeBool,
                Default: false,
                Optional: true,
            },
            "enabled": &schema.Schema{
                Type: schema.TypeBool,
                Default: true,
                Optional: true,
            },
            "log_time": &schema.Schema{
                Type: schema.TypeInt,
                Default: 60,
                Optional: true,
            },
            "note": &schema.Schema{
                Type: schema.TypeString,
                Default: "",
                Optional: true,
            },
            "per_process_connection_count": &schema.Schema{
                Type: schema.TypeBool,
                Default: true,
                Optional: true,
            },
            "rule": &schema.Schema{
                Type: schema.TypeString,
                Default: "",
                Optional: true,
            },
            "testing": &schema.Schema{
                Type: schema.TypeBool,
                Default: false,
                Optional: true,
            },
            "allowed_addresses": &schema.Schema{
                Type: schema.TypeSet,
                Optional: true,
                Elem: &schema.Schema{Type: schema.TypeString},
                Set: schema.HashString,
            },
            "banned_addresses": &schema.Schema{
                Type: schema.TypeSet,
                Optional: true,
                Elem: &schema.Schema{Type: schema.TypeString},
                Set: schema.HashString,
            },
            "max_10_connections": &schema.Schema{
                Type: schema.TypeInt,
                Default: 200,
                Optional: true,
            },
            "max_1_connections": &schema.Schema{
                Type: schema.TypeInt,
                Default: 30,
                Optional: true,
            },
            "max_connection_rate": &schema.Schema{
                Type: schema.TypeInt,
                Default: 0,
                Optional: true,
            },
            "min_connections": &schema.Schema{
                Type: schema.TypeInt,
                Default: 4,
                Optional: true,
            },
            "rate_timer": &schema.Schema{
                Type: schema.TypeInt,
                Default: 60,
                Optional: true,
            },
            "check_rfc2396": &schema.Schema{
                Type: schema.TypeBool,
                Default: false,
                Optional: true,
            },
            "max_body_length": &schema.Schema{
                Type: schema.TypeInt,
                Default: 0,
                Optional: true,
            },
            "max_header_length": &schema.Schema{
                Type: schema.TypeInt,
                Default: 0,
                Optional: true,
            },
            "max_request_length": &schema.Schema{
                Type: schema.TypeInt,
                Default: 0,
                Optional: true,
            },
            "max_url_length": &schema.Schema{
                Type: schema.TypeInt,
                Default: 0,
                Optional: true,
            },
            "reject_binary": &schema.Schema{
                Type: schema.TypeBool,
                Default: false,
                Optional: true,
            },
            "send_error_page": &schema.Schema{
                Type: schema.TypeBool,
                Default: true,
                Optional: true,
            },
        },
    }
}

func resourceProtectionSet(d *schema.ResourceData, meta interface{}) error {
    c := meta.(*providerConfig).client
    r := stingray.NewProtection(d.Get("name").(string))

    setBool(&r.Basic.Debug, d, "debug")
    setBool(&r.Basic.Enabled, d, "enabled")
    setInt(&r.Basic.LogTime, d, "log_time")
    setString(&r.Basic.Note, d, "note")
    setBool(&r.Basic.PerProcessConnectionCount, d, "per_process_connection_count")
    setBool(&r.Basic.Testing, d, "testing")
    setString(&r.Basic.Rule, d, "rule")

    setInt(&r.ConnectionLimiting.Max10Connections, d, "max_10_connections")
    setInt(&r.ConnectionLimiting.Max1Connections, d, "max_1_connections")
    setInt(&r.ConnectionLimiting.MaxConnectionRate, d, "max_connection_rate")
    setInt(&r.ConnectionLimiting.MinConnections, d, "min_connections")
    setInt(&r.ConnectionLimiting.RateTimer, d, "rate_timer")

    setBool(&r.Http.CheckRfc2396, d, "check_rfc2396")
    setInt(&r.Http.MaxBodyLength, d, "max_body_length")
    setInt(&r.Http.MaxHeaderLength, d, "max_header_length")
    setInt(&r.Http.MaxRequestLength, d, "max_request_length")
    setInt(&r.Http.MaxUrlLength, d, "max_url_length")
    setBool(&r.Http.RejectBinary, d, "reject_binary")
    setBool(&r.Http.SendErrorPage, d, "send_error_page")
    setStringSet(&r.AccessRestriction.Allowed, d, "allowed_addresses")
    setStringSet(&r.AccessRestriction.Banned, d, "banned_addresses")

    _, err := c.Set(r)
    if err != nil {
        return err
    }

    d.SetId(d.Get("name").(string))

    return nil
}

func resourceProtectionCreate(d *schema.ResourceData, meta interface{}) error {
    err := resourceProtectionSet(d, meta)
    if err != nil {
        return err
    }

    return resourceProtectionRead(d, meta)
}
func resourceProtectionRead(d *schema.ResourceData, meta interface{}) error {
    c := meta.(*providerConfig).client

    r, resp, err := c.GetProtection(d.Get("name").(string))
    if err != nil {
        if resp != nil && resp.StatusCode == 404 {
            // The resource doesn't exist anymore
            d.SetId("")
            return nil
        }

        return fmt.Errorf("Error reading resource: %s", err)
    }

    d.Set("debug", bool(*r.Basic.Debug))
    d.Set("enabled", bool(*r.Basic.Enabled))
    d.Set("log_time", int(*r.Basic.LogTime))
    d.Set("note", string(*r.Basic.Note))
    d.Set("per_process_connection_count", bool(*r.Basic.PerProcessConnectionCount))
    d.Set("rule", string(*r.Basic.Rule))
    d.Set("testing", bool(*r.Basic.Testing))

    d.Set("max_10_connections", int(*r.ConnectionLimiting.Max10Connections))
    d.Set("max_1_connections", int(*r.ConnectionLimiting.Max1Connections))
    d.Set("max_connection_rate", int(*r.ConnectionLimiting.MaxConnectionRate))
    d.Set("min_connections", int(*r.ConnectionLimiting.MinConnections))
    d.Set("rate_timer", int(*r.ConnectionLimiting.RateTimer))

    d.Set("check_rfc2396", bool(*r.Http.CheckRfc2396))
    d.Set("max_body_length", int(*r.Http.MaxBodyLength))
    d.Set("max_header_length", int(*r.Http.MaxHeaderLength))
    d.Set("max_request_length", int(*r.Http.MaxRequestLength))
    d.Set("max_url_length", int(*r.Http.MaxUrlLength))
    d.Set("reject_binary", bool(*r.Http.RejectBinary))
    d.Set("send_error_page", bool(*r.Http.SendErrorPage))

    d.Set("allowed_addresses", []string(*r.AccessRestriction.Allowed))
    d.Set("banned_addresses", []string(*r.AccessRestriction.Banned))

    return nil
}

func resourceProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
    err := resourceProtectionSet(d, meta)
    if err != nil {
        return err
    }

    return resourceProtectionRead(d, meta)
}

func resourceProtectionDelete(d *schema.ResourceData, meta interface{}) error {
    c := meta.(*providerConfig).client
    r := stingray.NewProtection(d.Id())

    _, err := c.Delete(r)
    if err != nil {
        return err
    }

    return nil
}
