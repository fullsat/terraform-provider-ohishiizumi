package ohishiizumi

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ohishiizumi_profile": resourceOhishiizumiProfile(),
			//			"ohishiizumi_friend": schema.DataSourceResourceShim(
			//				"template_file",
			//				dataSourceFile(),
			//			),
		},
	}
}
