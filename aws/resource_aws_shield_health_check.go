package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/shield"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsShieldHealthCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsShieldHealthCheckCreate,
		Read:   resourceAwsShieldHealthCheckRead,
		Delete: resourceAwsShieldHealthCheckDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"health_check_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protection_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArn,
			},
		},
	}
}

func resourceAwsShieldHealthCheckCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).shieldconn

	input := &shield.AssociateHealthCheckInput{
		HealthCheckArn: aws.String(d.Get("health_check_arn").(string)),
		ProtectionId:   aws.String(d.Get("protection_id").(string)),
	}

	_, err := conn.AssociateHealthCheck(input)
	if err != nil {
		return fmt.Errorf("Error creating Shield HealthCheck Association: %s", err)
	}
	d.SetId(d.Get("protecion_id").(string))
	return resourceAwsShieldHealthCheckRead(d, meta)
}

func resourceAwsShieldHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).shieldconn

	input := &shield.DescribeProtectionInput{
		ProtectionId: aws.String(d.Id()),
	}

	resp, err := conn.DescribeProtection(input)
	if err != nil {
		return fmt.Errorf("error reading Shield Protection (%s): %s", d.Id(), err)
	}
	d.Set("name", resp.Protection.Name)
	d.Set("resource_arn", resp.Protection.ResourceArn)
	return nil
}

func resourceAwsShieldHealthCheckDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).shieldconn

	input := &shield.DisassociateHealthCheckInput{
		HealthCheckArn: aws.String(d.Get("health_check_arn").(string)),
		ProtectionId:   aws.String(d.Get("protection_id").(string)),
	}

	_, err := conn.DisassociateHealthCheck(input)

	if isAWSErr(err, shield.ErrCodeResourceNotFoundException, "") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error disassociating Shield HealthCheck (%s): %s", d.Id(), err)
	}
	return nil
}
