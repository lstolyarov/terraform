package azurerm

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/servicefabric"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jen20/riviera/azure"
)

func resourceArmServiceFabric() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceFabricCreate,
		Read:   resourceArmServiceFabricRead,
		Update: resourceArmServiceFabricUpdate,
		Delete: resourceArmServiceFabricDelete,

		Schema: map[string]*schema.Schema{
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reliability_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"upgrade_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"management_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"client_endpoint_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"http_endpoint_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protected_account_key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"blob_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"queue_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"table_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": locationSchema(),
			"tags":     tagsSchema(),
			"cluster_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmServiceFabricCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	ServiceFabricClient := client.serviceFabricClient

	log.Printf("[INFO] preparing arguments for Azure ARM Service Fabric creation.")

	resGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	reliabilityLevel := d.Get("reliability_level").(string)
	upgradeMode := d.Get("upgrade_mode").(string)
	managementEndpoint := d.Get("management_endpoint").(string)
	nodeName := d.Get("node_name").(string)
	instanceCount := azure.Int32(int32(d.Get("instance_count").(int)))
	isPrimary := true
	clientEndpointPort := azure.Int32(int32(d.Get("client_endpoint_port").(int)))
	httpEndpointPort := azure.Int32(int32(d.Get("http_endpoint_port").(int)))
	vmImage := d.Get("vm_image").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	storageAccountName := d.Get("storage_account_name").(string)
	protectedAccountKeyName := d.Get("protected_account_key_name").(string)
	blobEndpoint := d.Get("blob_endpoint").(string)
	queueEndpoint := d.Get("queue_endpoint").(string)
	tableEndpoint := d.Get("table_endpoint").(string)

	nodeTypeDescription := servicefabric.NodeTypeDescription{
		Name:                         &nodeName,
		VMInstanceCount:              instanceCount,
		IsPrimary:                    &isPrimary,
		ClientConnectionEndpointPort: clientEndpointPort,
		HTTPGatewayEndpointPort:      httpEndpointPort,
	}

	nodeTypes := make([]servicefabric.NodeTypeDescription, 0)

	nodeTypes = append(nodeTypes, nodeTypeDescription)

	diagnosticsStorageAccountConfig := servicefabric.DiagnosticsStorageAccountConfig{
		StorageAccountName:      &storageAccountName,
		ProtectedAccountKeyName: &protectedAccountKeyName,
		BlobEndpoint:            &blobEndpoint,
		QueueEndpoint:           &queueEndpoint,
		TableEndpoint:           &tableEndpoint,
	}

	clusterProperties := servicefabric.ClusterProperties{
		ReliabilityLevel:   servicefabric.ReliabilityLevel(reliabilityLevel),
		UpgradeMode:        servicefabric.UpgradeMode(upgradeMode),
		ManagementEndpoint: &managementEndpoint,
		NodeTypes:          &nodeTypes,
		VMImage:            &vmImage,
		DiagnosticsStorageAccountConfig: &diagnosticsStorageAccountConfig,
	}

	cluster := servicefabric.Cluster{
		Location:          &location,
		ClusterProperties: &clusterProperties,
		Tags:              expandTags(tags),
	}

	_, error := ServiceFabricClient.Create(resGroup, clusterName, cluster, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err2 := ServiceFabricClient.Get(resGroup, clusterName)
	if err2 != nil {
		//return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Service Fabric %s (resource group %s) ID", clusterName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceFabricRead(d, meta)
}

func resourceArmServiceFabricUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	ServiceFabricClient := client.serviceFabricClient

	log.Printf("[INFO] preparing arguments for Azure ARM Service Fabric creation.")

	resGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	reliabilityLevel := d.Get("reliability_level").(string)
	upgradeMode := d.Get("upgrade_mode").(string)
	nodeName := d.Get("node_name").(string)
	instanceCount := azure.Int32(int32(d.Get("instance_count").(int)))
	clientEndpointPort := azure.Int32(int32(d.Get("client_endpoint_port").(int)))
	httpEndpointPort := azure.Int32(int32(d.Get("http_endpoint_port").(int)))
	tags := d.Get("tags").(map[string]interface{})
	isPrimary := true

	nodeTypeDescription := servicefabric.NodeTypeDescription{
		Name:                         &nodeName,
		VMInstanceCount:              instanceCount,
		IsPrimary:                    &isPrimary,
		ClientConnectionEndpointPort: clientEndpointPort,
		HTTPGatewayEndpointPort:      httpEndpointPort,
	}

	nodeTypes := make([]servicefabric.NodeTypeDescription, 0)

	nodeTypes = append(nodeTypes, nodeTypeDescription)

	clusterProperties := servicefabric.ClusterPropertiesUpdateParameters{
		ReliabilityLevel: servicefabric.ReliabilityLevel(reliabilityLevel),
		UpgradeMode:      servicefabric.UpgradeMode(upgradeMode),
		NodeTypes:        &nodeTypes,
	}

	clusterUpdateParameters := servicefabric.ClusterUpdateParameters{
		ClusterPropertiesUpdateParameters: &clusterProperties,
		Tags: expandTags(tags),
	}

	_, error := ServiceFabricClient.Update(resGroup, clusterName, clusterUpdateParameters, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err2 := ServiceFabricClient.Get(resGroup, clusterName)
	if err2 != nil {
		//return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Service Fabric %s (resource group %s) ID", clusterName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceFabricRead(d, meta)
}

func resourceArmServiceFabricRead(d *schema.ResourceData, meta interface{}) error {
	ServiceFabricClient := meta.(*ArmClient).serviceFabricClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading service fabric cluster %s", id)

	resGroup := id.ResourceGroup
	clusterName := id.Path["clusters"]

	resp, err := ServiceFabricClient.Get(resGroup, clusterName)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Service Fabric %s: %s", clusterName, err)
	}

	clusterEndpoint := resp.ClusterProperties.ClusterEndpoint

	d.Set("cluster_name", clusterName)
	d.Set("resource_group_name", resGroup)
	d.Set("cluster_endpoint", clusterEndpoint)

	return nil
}

func resourceArmServiceFabricDelete(d *schema.ResourceData, meta interface{}) error {
	ServiceFabricClient := meta.(*ArmClient).serviceFabricClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["clusters"]

	log.Printf("[DEBUG] Deleting service fabric cluster %s: %s", resGroup, name)

	_, err = ServiceFabricClient.Delete(resGroup, name)

	return err
}
