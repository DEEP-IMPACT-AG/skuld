# Hyperdrive Skuld

## Introduction

The `skuld` command-line utility is meant to be used by developers interacting with the AWS SDK and wanting/needing to use temporary credentials generated with Two-Factor authentication. `skuld` uses the AWS Security Token Service together with a [named profile](https://docs.aws.amazon.com/cli/latest/userguide/cli-multiple-profiles.html) and a code from an MFA device to generate another named profile with temporary credentials.

Together with appropriate IAM policies, `skuld` can enforce the use of an MFA device to manipulate the AWS SDK from a developer's machine.

## Installation

On Mac OS X, you can use [brew](https://brew.sh) to install `skuld`.

```
$ brew tap DEEP-IMPACT-AG/hyperdrive
$ brew install skuld
```

On Windows, you can use [scoop](https://scoop.sh) to install `skuld`.

```
$ scoop bucket add hyperdrive https://github.com/DEEP-IMPACT-AG/scoop-hyperdrive.git
$ scoop install skuld
```

For Linux, you can install manually by downloading from [latest release page](https://github.com/DEEP-IMPACT-AG/skuld/releases/latest).

Finally, you can install it from the sources via `go get`. You will need Go 1.10.

## Preparation

Before using `skuld`, you must create an IAM user, assign an MFA device to it and create an Access Key for it. Refer to the [IAM documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users.html) of AWS.

Copy the Access Key to a named profiled in the credentials file. Choose the region according to your most frequent usage.

```ini
[<profile-name>]
aws_access_key_id     = ??????
aws_secret_access_key = ??????
region                = us-east-1
```

You can check the IAM user of the named profile by using the `aws` command-line utility.

```
 $ aws --profile=<profile-name> sts get-caller-identity
```

You can also check the existence of your MFA device.

```
$ aws --profile=<profile-name> iam list-mfa-devices --user-name <iam-user-name>
```

## Usage

To request temporary credentials, use `skuld` at the shell as follows:

```
$ skuld <profile-name>
Enter your token: 
```

When prompted by `Enter your token: `, enter the token of your MFA device and press the enter key.

`skuld` will fetch temporary credentials and create a new profile named  `<profile-name>-skuld` with them (i.e. the new profile's name is the original profile name with the suffix `-skuld`). If the skuld profile already exists, it will be overwritten with the new temporary credentials.

`skuld` also ouputs the expiring time of the temporary credentials in UTC:

```
Credentials valid until: 2018-01-02 20:00:01 +0000 UTC
```

The temporary credentials are valid for 10 hours but if the profile name ends with `-adm`; in that case, the temporary credentials are valid for 1 hour.

The new profile, respectively updated profile, can now be used normally. For instance, to describe ec2 instances:

```
$ aws --profile=<profile-name>-skuld ec2 describe-instances
```

Or used in the credentials file to be used as reference in other named profiles:

```ini
[other-profile-name]
source_profile = <profile-name>-skuld
role_arn       = arn:aws:iam::xxxxxx:role/admin
region         = us-east-1
```

The region of the skuld profile is given by the profile from which it is derived. For instance, if the base profile is in the `us-east-2` region, the skuld profile will be also configured to be in the `us-east-2`. Beside the configuration in the credentials file, `skuld` will also generate an entry in the configuraion `~/.aws/config` with the region.

The flag `-r <region>` can be used to override the region.

## Enforcing MFA Devices

`skuld` by itself does not enforce the use of MFA devices; it just simplifies the creation of temporary credentials with MFA devices.

To actually enforce the use of MFA Devices, you need to assign a proper IAM policy to your IAM user.

AWS has a tutorial to that purpose: [Enable Your Users to Configure Their Own Credentials and MFA Settings](https://docs.aws.amazon.com/IAM/latest/UserGuide/tutorial_users-self-manage-mfa-and-creds.html).