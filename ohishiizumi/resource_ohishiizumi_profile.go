package ohishiizumi

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"

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
				ForceNew:    true,
			},
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

	err := profile.Create()
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
	id := d.Id()

	profile := &Profile{}
	err := profile.Read(id)

	if err != nil {
		if os.IsNotExist(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", profile.Name)
	d.Set("age", profile.Age)
	d.Set("height", profile.Height)

	return nil
}

func resourceOhishiizumiProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	profile := &Profile{}
	err := profile.Read(id)

	if err != nil {
		if os.IsNotExist(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if d.HasChange("name") {
		profile.Name = d.Get("name").(string)
	}

	if d.HasChange("age") {
		profile.Age = d.Get("age").(int)
	}

	if d.HasChange("height") {
		profile.Height = d.Get("height").(int)
	}

	err = profile.Update()
	if err != nil {
		return err
	}

	return resourceOhishiizumiProfileRead(d, meta)
}

func resourceOhishiizumiProfileDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	profile := &Profile{}
	if err := profile.Delete(id); err != nil {
		if os.IsNotExist(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.SetId("")

	return nil
}

type Profile struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Height int    `json:"height"`
}

func (p *Profile) Create() error {
	body, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return fmt.Errorf("to json failed")
	}

	path := p.calcID() + ".json"

	err = ioutil.WriteFile(path, body, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file")
	}

	return nil
}

func (p *Profile) Read(id string) error {
	path := id + ".json"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	return nil
}

func (p *Profile) Update() error {
	return p.Create()
}

func (p *Profile) Delete(id string) error {
	path := id + ".json"

	if _, err := os.Stat(path); err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (p *Profile) calcID() string {
	h := fnv.New32a()
	h.Write([]byte(p.Name))
	return fmt.Sprintf("%d", h.Sum32())
}
