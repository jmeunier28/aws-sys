package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
)

//ListUsers returns all Iam users in the account
func (iamSess *Client) ListUsers() (*iam.ListUsersOutput, error) {

	svc := iam.New(iamSess.session)
	user := &iam.ListUsersInput{
		// get all the users at once
		MaxItems: aws.Int64(500),
	}
	resp, err := svc.ListUsers(user)
	if err != nil {
		return nil, err
	}

	return resp, err
}

//DeleteUser deletes user based on username
func (iamSess *Client) DeleteUser(uid *string) error {

	svc := iam.New(iamSess.session)

	// delete dependents first
	err := iamSess.deleteUserDependencies(uid)
	if err != nil {
		return err
	}
	user := &iam.DeleteUserInput{
		UserName: aws.String(*uid),
	}

	_, err = svc.DeleteUser(user)
	if err != nil {
		return err
	}
	return err
}

func (iamSess *Client) removeUserFromGroup(uid *string) (bool, error) {
	svc := iam.New(iamSess.session)
	group := &iam.ListGroupsForUserInput{
		UserName: aws.String(*uid),
	}
	groups, err := svc.ListGroupsForUser(group)
	if err == nil {
		// get the group names and remove them from them
		for _, g := range groups.Groups {
			removeUser := iam.RemoveUserFromGroupInput{
				GroupName: aws.String(*g.GroupName),
				UserName:  aws.String(*uid),
			}
			_, err = svc.RemoveUserFromGroup(&removeUser)
			if err != nil {
				return false, err
			}
		}

	}
	return true, nil
}

// private functions to delete all user pre reqs
func (iamSess *Client) deleteLoginProfile(uid *string) (bool, error) {
	svc := iam.New(iamSess.session)

	del := &iam.DeleteLoginProfileInput{
		UserName: aws.String(*uid),
	}
	_, err := svc.DeleteLoginProfile(del)
	if err != nil {
		if awsErr, ok := err.(awserr.RequestFailure); ok {
			if awsErr.StatusCode() == 404 {
				// user didnt exist anyways
				return true, nil
			}
		} else {

			return false, err
		}

	}

	return true, nil

}

func (iamSess *Client) deleteMFA(uid *string) (bool, error) {
	svc := iam.New(iamSess.session)
	mfa := &iam.ListMFADevicesInput{
		UserName: aws.String(*uid),
	}

	// find the device and delete it
	mfaResp, err := svc.ListMFADevices(mfa)
	if err == nil {
		for _, k := range mfaResp.MFADevices {
			// delete
			deleteMfa := &iam.DeactivateMFADeviceInput{
				UserName:     aws.String(*uid),
				SerialNumber: aws.String(*k.SerialNumber),
			}
			_, err = svc.DeactivateMFADevice(deleteMfa)
			if err != nil {
				return false, err
			}
		}

	}
	return true, nil
}

func (iamSess *Client) deleteAccessKeys(uid *string) (bool, error) {
	svc := iam.New(iamSess.session)
	key := &iam.ListAccessKeysInput{
		UserName: aws.String(*uid),
	}
	keys, err := svc.ListAccessKeys(key)
	// if keys exist delte them
	if err == nil {
		// delete em
		for _, l := range keys.AccessKeyMetadata {
			del := &iam.DeleteAccessKeyInput{
				UserName:    aws.String(*uid),
				AccessKeyId: aws.String(*l.AccessKeyId),
			}
			_, err = svc.DeleteAccessKey(del)
			if err != nil {
				return false, err
			}
		}

	}
	return true, nil
}

func (iamSess *Client) deleteSigningCert(uid *string) (bool, error) {
	svc := iam.New(iamSess.session)
	cert := &iam.ListSigningCertificatesInput{
		UserName: aws.String(*uid),
	}
	certResp, err := svc.ListSigningCertificates(cert)
	if err == nil {
		// delete em
		for _, c := range certResp.Certificates {
			deleteCert := &iam.DeleteSigningCertificateInput{
				CertificateId: aws.String(*c.CertificateId),
			}
			_, err = svc.DeleteSigningCertificate(deleteCert)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func (iamSess *Client) detachPolicy(uid *string) (bool, error) {
	svc := iam.New(iamSess.session)

	// Detach all the policies
	policy := &iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(*uid),
	}
	policyResp, err := svc.ListAttachedUserPolicies(policy)
	if err == nil {
		for _, p := range policyResp.AttachedPolicies {
			// delte
			detachPolicy := &iam.DetachUserPolicyInput{
				UserName:  aws.String(*uid),
				PolicyArn: aws.String(*p.PolicyArn),
			}
			_, err = svc.DetachUserPolicy(detachPolicy)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func (iamSess *Client) deleteUserDependencies(uid *string) error {
	var errors []error
	_, err := iamSess.deleteLoginProfile(uid) // get rid of login profile
	errors = append(errors, err)
	_, err = iamSess.deleteMFA(uid) // detach mfa device from user
	errors = append(errors, err)
	_, err = iamSess.deleteSigningCert(uid) // delete any signing certs
	errors = append(errors, err)
	_, err = iamSess.deleteAccessKeys(uid) // delete any access keys
	errors = append(errors, err)
	_, err = iamSess.detachPolicy(uid) // Detach any policies
	errors = append(errors, err)
	_, err = iamSess.removeUserFromGroup(uid) // remove the user from any groups theyre int
	errors = append(errors, err)

	// iterate through and make sure all was done flawlessly
	for _, e := range errors {
		if e != nil {
			return e
		}
	}
	return nil
}
