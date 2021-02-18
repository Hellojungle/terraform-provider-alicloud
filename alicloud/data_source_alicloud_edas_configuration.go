package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEdasConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEdasConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"data_id_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"data_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_id_pattern": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEdasConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateListConfigCentersRequest()
	request.RegionId = client.RegionId
	request.LogicalRegionId = d.Get("logical_region_id").(string)
	request.Group = d.Get("group").(string)

	if v, ok := d.GetOk("data_id_pattern"); ok {
		request.DataIdPattern = v.(string)
	}

	if v, ok := d.GetOk("app_name"); ok {
		request.AppName = v.(string)
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			if id == nil {
				continue
			}
			idsMap[Trim(id.(string))] = Trim(id.(string))
		}
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListConfigCenters(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_configurations", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ListConfigCentersResponse)
	if response.Code != 200 {
		return WrapError(Error(response.Message))
	}
	var filteredConfigs []edas.ListConfigCenters
	nameRegex, ok := d.GetOk("data_id_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, conf := range response.ConfigCentersList.ListConfigCenters {
			if r != nil && !r.MatchString(conf.DataId) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[conf.Id]; !ok {
					continue
				}
			}
			filteredConfigs = append(filteredConfigs, conf)
		}
	} else {
		filteredConfigs = response.ConfigCentersList.ListConfigCenters
	}

	return edasConfigurationsAttributes(d, filteredConfigs)
}

func edasConfigurationsAttributes(d *schema.ResourceData, confs []edas.ListConfigCenters) error {
	var appIds []string
	var s []map[string]interface{}
	var names []string

	for _, conf := range confs {
		mapping := map[string]interface{}{
			"app_name":         conf.AppName,
			"data_id":          conf.DataId,
			"group": conf.Group,
			"id": conf.Id,
		}
		appIds = append(appIds, conf.Id)
		names = append(names, conf.DataId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(appIds))
	if err := d.Set("ids", appIds); err != nil {
		return WrapError(err)
	}
	if err := d.Set("data_ids", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("configurations", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
