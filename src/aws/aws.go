package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
)

type AWS struct {
	PrincipalKey   string
	PrincipalValue string
	Key            string
	Secret         string
	Region         string
	Endpoint       string
	Bucket         string
}

type RolePolicyDocument struct {
	Version   string
	Statement []RoleStatementEntry
}

type RoleStatementEntry struct {
	Effect    string
	Action    []string
	Principal map[string]string
}

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

type StatementEntry struct {
	Effect    string
	Action    []string
	Resource  string
	Condition map[string]interface{} `json:",omitempty"`
}

func makeRolePolicy(principal map[string]string) string {
	p := RolePolicyDocument{
		Version: "2012-10-17",
		Statement: []RoleStatementEntry{{
			Effect:    "Allow",
			Action:    []string{"sts:AssumeRole"},
			Principal: principal,
		}},
	}
	byte, _ := json.Marshal(p)
	return string(byte)
}

func makePolicy(bucket, path string) string {
	p := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			{
				Effect:   "Allow",
				Action:   []string{"s3:ListAllMyBuckets", "s3:GetBucketLocation"},
				Resource: "arn:aws:s3:::*",
			},
			{
				Effect:   "Allow",
				Action:   []string{"s3:ListBucket"},
				Resource: fmt.Sprintf("arn:aws:s3:::%s", bucket),
				Condition: map[string]interface{}{
					"StringEquals": map[string]interface{}{
						"s3:prefix":    []string{"", path},
						"s3:delimiter": []string{"/"},
					},
				},
			},
			{
				Effect:   "Allow",
				Action:   []string{"s3:ListBucket"},
				Resource: fmt.Sprintf("arn:aws:s3:::%s", bucket),
				Condition: map[string]interface{}{
					"StringLike": map[string]interface{}{
						"s3:prefix": []string{fmt.Sprintf("%s/*", path)},
					},
				},
			},
			{
				Effect:   "Allow",
				Action:   []string{"s3:*"},
				Resource: fmt.Sprintf("arn:aws:s3:::%s/%s/*", bucket, path),
			},
		},
	}
	byte, _ := json.Marshal(p)
	return string(byte)
}

func getRole(svcIam *iam.IAM, name string, principal map[string]string) (*iam.Role, error) {
	if re, err := svcIam.GetRole(&iam.GetRoleInput{RoleName: aws.String(name)}); err == nil {
		return re.Role, nil
	}
	if re, err := svcIam.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(makeRolePolicy(principal)),
		RoleName:                 aws.String(name),
	}); err == nil {
		return re.Role, nil
	} else {
		log.Debugf("policy:%s", string(makeRolePolicy(principal)))
		return nil, err
	}
}

func getPolicy(svcIam *iam.IAM, name, bucket, path string) (*iam.Policy, error) {
	reList, err := svcIam.ListPolicies(&iam.ListPoliciesInput{PathPrefix: aws.String("/" + name + "/")})
	if err == nil && len(reList.Policies) == 1 {
		return reList.Policies[0], nil
	}
	if rePol, err := svcIam.CreatePolicy(&iam.CreatePolicyInput{
		Path:           aws.String("/" + name + "/"),
		PolicyName:     aws.String(name),
		PolicyDocument: aws.String(makePolicy(bucket, path)),
	}); err == nil {
		return rePol.Policy, nil
	} else {
		return nil, err
	}
}

func (s *AWS) GetCreds(path string) (*credentials.Credentials, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(s.Key, s.Secret, ""),
	})
	if err != nil {
		log.Errorf("Error creating aws session, error:%s", err)
		return nil, err
	}

	name := fmt.Sprintf("%s-%s", s.Bucket, path)
	svcIam := iam.New(sess)
	role, err := getRole(svcIam, name, map[string]string{s.PrincipalKey: s.PrincipalValue})
	if err != nil {
		log.Errorf("Error getting role name:%s, principal:%s/%s error:%s", name, s.PrincipalKey, s.PrincipalValue, err)
		return nil, err
	}
	pol, err := getPolicy(svcIam, name, s.Bucket, path)
	if err != nil {
		log.Errorf("Error getting policy, error:%s", err)
		return nil, err
	}
	_, err = svcIam.AttachRolePolicy(&iam.AttachRolePolicyInput{
		PolicyArn: pol.Arn,
		RoleName:  role.RoleName,
	})
	if err != nil {
		log.Warnf("Attach policy error, err:%s", err)
	}
	return stscreds.NewCredentials(sess, *role.Arn), nil
}
