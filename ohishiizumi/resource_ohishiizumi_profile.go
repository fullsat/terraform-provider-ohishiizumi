package ohishiizumi

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOhishiizumiProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceOhishiizumiProfileCreate,
		Read:   resourceOhishiizumiProfileRead,
		Update: resourceOhishiizumiProfileUpdate,
		Delete: resourceOhishiizumiProfileDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "大石泉",
				Required:    true,
			},
			//			"friends": {
			//				Type:         schema.TypeMap,
			//				Optional:     true,
			//				Default:      make(map[string]interface{}),
			//				Description:  "Variables to substitute",
			//				ForceNew:     true,
			//			},
			"age": {
				Type:        schema.TypeInt,
				Description: "Path to the directory where the templated files will be written",
				Required:    true,
			},
			"height": {
				Type:        schema.TypeInt,
				Description: "Path to the directory where the templated files will be written",
				Required:    true,
			},
		},
	}
}

func resourceOhishiizumiProfileCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	age := d.Get("age").(int)
	height := d.Get("height").(int)

	profile := &Profile{
		Name:   name,
		Age:    age,
		Height: height,
	}

	_, err := profile.Create()
	if err != nil {
		return err
	}

	d.SetId(profile.calcID())
	d.Set("name", name)
	d.Set("age", age)
	d.Set("height", height)

	return nil
}

func resourceOhishiizumiProfileRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOhishiizumiProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOhishiizumiProfileDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

type Profile struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Height int    `json:"height"`
}

func (p *Profile) Create() (string, error) {
	path := "ohishiizumi.json"

	body, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return "", fmt.Errorf("to json failed")
	}

	err = ioutil.WriteFile(path, body, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file")
	}

	return string(body), nil
}

func (p *Profile) calcID() string {
	h := fnv.New32a()
	h.Write([]byte(p.Name))
	return fmt.Sprintf("%d", h.Sum32())
}

func (p *Profile) Read() (string, error) {
}
