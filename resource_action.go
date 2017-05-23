package main

import (
	"fmt"
	"github.com/atlassian/go-vtm"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceActionCreate,
		Read:   resourceActionRead,
		Update: resourceActionUpdate,
		Delete: resourceActionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"verbose": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"syslog_msg_len_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1024,
			},
			"program": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"argument": &schema.Schema{
							Type: schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"value": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
							Optional: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"syslog": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sysloghost": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"soap": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_data": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"proxy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"log": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"from": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"trap": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_password": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"community": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"hash_algorithm": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "md5",
						},
						"priv_password": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"traphost": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "snmpv1",
						},
					},
				},
			},
			"email": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"to": &schema.Schema{
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceActionSet(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client
	r := stingray.NewAction(d.Get("name").(string))

	setBool(&r.Basic.Verbose, d, "verbose")
	setString(&r.Basic.Type, d, "type")
	setString(&r.Basic.Note, d, "note")
	setInt(&r.Basic.SyslogMsgLenLimit, d, "syslog_msg_len_limit")
	setInt(&r.Basic.Timeout, d, "timeout")

	if v, ok := d.GetOk("log"); ok {
		vL := v.(*schema.Set).List()

		for _, v := range vL {
			log := v.(map[string]interface{})
			r.Log.File = stingray.String(log["file"].(string))
			r.Log.From = stingray.String(log["from"].(string))
		}
	}

	if v, ok := d.GetOk("email"); ok {
		vL := v.(*schema.Set).List()

		for _, v := range vL {
			email := v.(map[string]interface{})
			r.Email.Server = stingray.String(email["server"].(string))
		}
	}

	if v, ok := d.GetOk("SOAP"); ok {
		vL := v.(*schema.Set).List()

		for _, v := range vL {
			soap := v.(map[string]interface{})

			r.SOAP.AdditionalData = stingray.String(soap["server"].(string))
			r.SOAP.Password = stingray.String(soap["password"].(string))
			r.SOAP.Proxy = stingray.String(soap["proxy"].(string))
			r.SOAP.Username = stingray.String(soap["username"].(string))
		}
	}

	if v, ok := d.GetOk("program"); ok {
		vL := v.(*schema.Set).List()

		for _, v := range vL {
			program := v.(map[string]interface{})
			r.Program.Program = stingray.String(program["name"].(string))
			if arguments, ok := program["argument"]; ok {
				for _, data := range arguments.(*schema.Set).List() {
					argument := data.(map[string]interface{})
					vtmArg := &stingray.ProgramArgument{}

					vtmArg.Name = stingray.String(argument["name"].(string))
					vtmArg.Value = stingray.String(argument["value"].(string))
					vtmArg.Description = stingray.String(argument["description"].(string))

					r.Program.Arguments = append(r.Program.Arguments, *vtmArg)
				}
			} else {
				r.Program.Arguments = make([]stingray.ProgramArgument, 0, 0)
			}
		}
	}

	if v, ok := d.GetOk("trap"); ok {
		vL := v.(*schema.Set).List()

		for _, v := range vL {
			trap := v.(map[string]interface{})
			r.Trap.AuthPassword = stingray.String(trap["auth_password"].(string))
			r.Trap.Community = stingray.String(trap["community"].(string))
			r.Trap.HashAlgorithm = stingray.String(trap["hash_algorithm"].(string))
			r.Trap.PrivPassword = stingray.String(trap["priv_password"].(string))
			r.Trap.TrapHost = stingray.String(trap["traphost"].(string))
			r.Trap.Username = stingray.String(trap["username"].(string))
		}
	}

	_, err := c.Set(r)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func resourceActionRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client

	r, resp, err := c.GetAction(d.Get("name").(string))
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// The resource doesn't exist anymore
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading resource: %s", err)
	}

	d.Set("verbose", bool(*r.Basic.Verbose))
	d.Set("note", string(*r.Basic.Note))
	d.Set("type", string(*r.Basic.Type))
	d.Set("timeout", int(*r.Basic.Timeout))
	d.Set("syslog_msg_len_limit", int(*r.Basic.SyslogMsgLenLimit))

	syslog, syslogList := nestSetHelper()
	syslog["sysloghost"] = string(*r.Syslog.SyslogHost)
	if "syslog" == d.Get("type") {
		d.Set("syslog", syslogList)
	}

	program, programList := nestSetHelper()
	program["name"] = string(*r.Program.Program)
	if "program" == d.Get("type") {
		argumentList := make([]map[string]interface{}, 0, len(r.Program.Arguments))
		for _, argument := range r.Program.Arguments {
			item := make(map[string]interface{})
			item["name"] = string(*argument.Name)
			item["value"] = string(*argument.Value)
			item["description"] = string(*argument.Description)
			argumentList = append(argumentList, item)
		}

		program["argument"] = argumentList
		d.Set("program", programList)
	}

	email, emailList := nestSetHelper()
	email["server"] = string(*r.Email.Server)
	if "email" == d.Get("type") {
		d.Set("email", emailList)
	}

	soap, soapList := nestSetHelper()
	soap["additional_data"] = string(*r.SOAP.AdditionalData)
	soap["password"] = string(*r.SOAP.Password)
	soap["proxy"] = string(*r.SOAP.Proxy)
	soap["username"] = string(*r.SOAP.Username)
	if "soap" == d.Get("type") {
		d.Set("soap", soapList)
	}

	log, logList := nestSetHelper()
	log["file"] = string(*r.Log.File)
	log["from"] = string(*r.Log.From)
	if "log" == d.Get("type") {
		d.Set("log", logList)
	}

	trap, trapList := nestSetHelper()
	trap["auth_password"] = string(*r.Trap.AuthPassword)
	trap["community"] = string(*r.Trap.Community)
	trap["hash_algorithm"] = string(*r.Trap.HashAlgorithm)
	trap["priv_password"] = string(*r.Trap.PrivPassword)
	trap["traphost"] = string(*r.Trap.TrapHost)
	trap["username"] = string(*r.Trap.Username)
	trap["version"] = string(*r.Trap.Version)
	if "trap" == d.Get("type") {
		d.Set("trap", trapList)
	}

	return nil
}

func resourceActionCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceActionSet(d, meta)
	if err != nil {
		return err
	}

	return resourceActionRead(d, meta)
}
func resourceActionUpdate(d *schema.ResourceData, meta interface{}) error {
	err := resourceActionSet(d, meta)
	if err != nil {
		return err
	}

	return resourceActionRead(d, meta)
}

func resourceActionDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client
	r := stingray.NewAction(d.Id())

	_, err := c.Delete(r)
	if err != nil {
		return err
	}

	return nil
}
