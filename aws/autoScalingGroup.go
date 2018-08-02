package aws

import (
					"os"
	"fmt"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var exitFn = os.Exit

/*func main() {

	updateAsgAmi("us-west-1", "hgokavarapu-test-asg", "ami-3b01ec58")
}*/

func updateAsgAmi(region string, asgName string, amiId string) {
	client, err := autoscalingClient(region)
	fatalfIfErr("Error in creating autoscalinggroup client : %v", err)

	// Get LaunchConfig from ASG
	asg := getAutoScalingGroup(client, asgName)

	if len(asg.AutoScalingGroups) == 0 {
		fmt.Println("ERROR: No autoscaling group found with the specified name")
	}
	launchConfigName := *asg.AutoScalingGroups[0].LaunchConfigurationName
	fmt.Println("Launch Config Name:", launchConfigName)

	launchConfig := getLaunchConfig(client, launchConfigName)
	fmt.Println(launchConfig)

	newLaunchConfigName := "hemanth-test"

	// copy data from existing launch config and create new launch config
	newInput := createNewLaunchConfigFromCopy(newLaunchConfigName, amiId, launchConfig)

	// validate it before using it
	err = newInput.Validate()
	fatalfIfErr("Error in new launch config created : %v", err)

	// create new launch config
	createNewLaunchConfig(client, newInput)

	// update ASG with new LaunchConfig
	updateAsgRequest := autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName:    &asgName,
		LaunchConfigurationName: &newLaunchConfigName,
	}

	_, err = client.UpdateAutoScalingGroup(&updateAsgRequest)
	fatalfIfErr("Error while updating the ASG with new Launch Config: %v", err)

	// check if asg is updated with new LC
	asg = getAutoScalingGroup(client, asgName)
	fmt.Println("Updated ", asgName, " with Launch Config Name:", *asg.AutoScalingGroups[0].LaunchConfigurationName)
}

func getAutoScalingGroup(client *autoscaling.AutoScaling, asgName string) *autoscaling.DescribeAutoScalingGroupsOutput {
	asgList := []*string{&asgName}
	asgRequestInput := autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: asgList,
	}

	asg, err := client.DescribeAutoScalingGroups(&asgRequestInput)
	fatalfIfErr("Error while getting ASG details : %v", err)
	return asg
}


func createNewLaunchConfig(client *autoscaling.AutoScaling, newInput autoscaling.CreateLaunchConfigurationInput) {
	req, _ := client.CreateLaunchConfigurationRequest(&newInput)
	err := req.Send()
	fatalfIfErr("Error creating new launch config: %v", err)
}

func getLaunchConfig(client *autoscaling.AutoScaling, name string) *autoscaling.DescribeLaunchConfigurationsOutput {
	list := []*string{&name}
	input := autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: list,
	}
	out, err := client.DescribeLaunchConfigurations(&input)
	fatalfIfErr("Error describing the launch config: %v", err)
	return out
}

func createNewLaunchConfigFromCopy(name string, amiId string, out *autoscaling.DescribeLaunchConfigurationsOutput) autoscaling.CreateLaunchConfigurationInput {
	// create new input
	if len(out.LaunchConfigurations) == 0 {
		fmt.Println("Error: No LaunchConfig was found with the given name")
		exitFn(1)
	}

	launchConfig := out.LaunchConfigurations[0]
	newInput := autoscaling.CreateLaunchConfigurationInput{}

	newInput.AssociatePublicIpAddress = launchConfig.AssociatePublicIpAddress
	newInput.BlockDeviceMappings = launchConfig.BlockDeviceMappings
	newInput.KeyName = launchConfig.KeyName
	newInput.SecurityGroups = launchConfig.SecurityGroups
	newInput.UserData = launchConfig.UserData
	newInput.InstanceType = launchConfig.InstanceType
	newInput.InstanceMonitoring = launchConfig.InstanceMonitoring
	newInput.EbsOptimized = launchConfig.EbsOptimized
	newInput.IamInstanceProfile = launchConfig.IamInstanceProfile
	if len(*launchConfig.KernelId) > 0 {
		newInput.KernelId = launchConfig.KernelId
	}
	if len(*launchConfig.RamdiskId) > 0 {
		newInput.RamdiskId = launchConfig.RamdiskId
	}

	// set ami id
	newInput.ImageId = &amiId

	// set new name for launchConfig
	newInput.LaunchConfigurationName = &name

	return newInput
}

func autoscalingClient(region string) (*autoscaling.AutoScaling, error) {
	awsSession, err := session.NewSession(aws.NewConfig().WithRegion(region))
	fatalfIfErr("aws error: %v", err)
	client := autoscaling.New(awsSession)
	return client, err
}

func fatalfIfErr(format string, err error) {
	if err != nil {
		fatalErr(format, err)
	}
}

func fatalErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	exitFn(1)
}

