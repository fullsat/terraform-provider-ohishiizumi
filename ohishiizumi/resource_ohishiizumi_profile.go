package ohishiizumi

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
)

// resource "hogehoge"に相当するもの
func resourceOhishiizumiProfile() *schema.Resource {
	return &schema.Resource{
		// Create
		//   リソースの作成
		//   tfstateに書き込むためにResourceDataでSetする
		//   IDは一意的なものをSetIdで設定する
		Create: resourceOhishiizumiProfileCreate,
		// Read
		//   引数から渡ってくるResourceDataに値をSetする
		//   値は実際のリソースの値を見て設定する
		Read: resourceOhishiizumiProfileRead,
		// Update
		//   Createと基本的には同様
		//   HasChangeヘルパーなどを駆使して差分だけ適用するなどする
		Update: resourceOhishiizumiProfileUpdate,
		// Delete
		//   リソースの削除
		//   SetId("")として削除したことをterraformに伝える
		//   ただし明示的に指定しなくても内部的に呼ばれているらしい
		Delete: resourceOhishiizumiProfileDelete,

		// resource "hogehoge" {
		//    hogehoge = ここに書くものを定義していく
		// }
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Ohishi Izumi",
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Cool or Cute or Passion",
				Required:    true,
			},
			"height": {
				Type:        schema.TypeInt,
				Description: "Height",
				Required:    true,
			},
			"weight": {
				Type:        schema.TypeInt,
				Description: "Weight",
				Required:    true,
			},
			"body": {
				Type:        schema.TypeMap,
				Description: "body",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"age": {
				Type:        schema.TypeInt,
				Description: "年齢",
				Required:    true,
			},
			"birthday": {
				Type:        schema.TypeString,
				Description: "age",
				Required:    true,
			},
		},
	}
}

func resourceOhishiizumiProfileCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	idol_type := d.Get("type").(string)
	age := d.Get("age").(int)
	height := d.Get("height").(int)
	weight := d.Get("weight").(int)
	birthday := d.Get("birthday").(string)

	body := d.Get("body").(map[string]interface{})
	hip := body["hip"].(int)
	waist := body["waist"].(int)
	bust := body["hip"].(int)

	threesize := &ThreeSize{
		Bust:  bust,
		Waist: waist,
		Hip:   hip,
	}

	profile := &Profile{
		Name:     name,
		Age:      age,
		Type:     idol_type,
		Height:   height,
		Weight:   weight,
		Body:     *threesize,
		Birthday: birthday,
	}

	err := profile.Create()
	if err != nil {
		return err
	}

	d.SetId(profile.calcID())
	d.Set("name", name)
	d.Set("age", age)
	d.Set("height", height)
	d.Set("type", idol_type)
	d.Set("age", age)
	d.Set("weight", weight)
	d.Set("birthday", birthday)
	d.Set("body", body)

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
	d.Set("type", profile.Type)
	d.Set("height", profile.Height)
	d.Set("weight", profile.Weight)
	d.Set("body", profile.Body)
	d.Set("age", profile.Age)
	d.Set("birthday", profile.Birthday)

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

	if d.HasChange("type") {
		profile.Age = d.Get("type").(int)
	}

	if d.HasChange("height") {
		profile.Height = d.Get("height").(int)
	}

	if d.HasChange("weight") {
		profile.Name = d.Get("weight").(string)
	}

	if d.HasChange("body") {
		body := d.Get("body").(map[string]interface{})
		profile.Body.Bust = body["bust"].(int)
		profile.Body.Waist = body["waist"].(int)
		profile.Body.Hip = body["hip"].(int)
	}

	if d.HasChange("age") {
		profile.Age = d.Get("age").(int)
	}

	if d.HasChange("birthday") {
		profile.Birthday = d.Get("birthday").(string)
	}

	err = profile.Update()
	if err != nil {
		return err
	}

	return resourceOhishiizumiProfileRead(d, meta)
}

func resourceOhishiizumiProfileDelete(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("大石泉ちゃんは何があろうが壊せない")
}

type Profile struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Weight   int       `json:"weight"`
	Height   int       `json:"height"`
	Age      int       `json:"age"`
	Birthday string    `json:"birthday"`
	Body     ThreeSize `json:"body"`
}

type ThreeSize struct {
	Bust  int `json:"bust"`
	Waist int `json:"waist"`
	Hip   int `json:"hip"`
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
