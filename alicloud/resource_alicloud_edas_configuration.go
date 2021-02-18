package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEdasConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: rresourceAlicloudEdasConfigurationCreate,
		Read:   resourceAlicloudEdasConfigurationRead,
		Delete: resourceAlicloudEdasConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func rresourceAlicloudEdasConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateInsertConfigCenterRequest()
	request.RegionId = client.RegionId
	request.DataId = d.Get("data_id").(string)
	request.Group = d.Get("group").(string)
	request.Data = d.Get("data").(string)
	request.LogicalRegionId = d.Get("logical_region_id").(string)
	if v, ok := d.GetOk("app_name"); ok {
		request.AppName = v.(string)
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.InsertConfigCenter(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_configuration", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.InsertConfigCenterResponse)
	if response.Code != 200 {
		return WrapError(Error("create configuration failed for " + response.Message))
	}

	return resourceAlicloudEdasConfigurationRead(d, meta)
}

func resourceAlicloudEdasConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	regionId := client.RegionId

	request := edas.CreateQueryConfigCenterRequest()
	request.RegionId = regionId
	request.DataId = d.Get("data_id").(string)
	request.Group = d.Get("group").(string)
	request.LogicalRegionId = d.Get("logical_region_id").(string)

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.QueryConfigCenter(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_configuration", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.QueryConfigCenterResponse)
	if response.Code != 200 {
		return WrapError(Error("create configuration failed for " + response.Message))
	}

	d.Set("app_name", response.ConfigCenterInfo.AppName)
	d.Set("data", response.ConfigCenterInfo.Content)
	d.Set("id", response.ConfigCenterInfo.Id)
	d.SetId(response.ConfigCenterInfo.Id)

	return nil
}

func resourceAlicloudEdasConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	regionId := client.RegionId

	request := edas.CreateDeleteConfigCenterRequest()
	request.RegionId = regionId
	request.DataId = d.Get("data_id").(string)
	request.Group = d.Get("group").(string)
	request.LogicalRegionId = d.Get("logical_region_id").(string)

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteConfigCenter(request)
		})
		response, _ := raw.(*edas.DeleteConfigCenterResponse)
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if response.Code != 200 {
			return resource.NonRetryableError(Error("delete configuration failed for " + response.Message))
		}

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}
