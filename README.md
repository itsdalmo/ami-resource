##  AMI resource (work in progress)

Concourse resource for AWS AMIs that is very much based on
[this resource](https://github.com/jdub/ami-resource) by the same name. A 
new resource was created because the old one has not been maintained in the
last 12 months.

## Source Configuration

- `aws_access_key_id`: *Optional*: Access key id if you are passing credentials.
- `aws_secret_access_key`: *Optional*: See above.
- `aws_session_token`: *Optional*: Use if your access/secret keys are temporary (assumed role/MFA authenticated).
- `aws_region`: Region where the images of interest live.
- `filters`: A map (name: value) of filters for your AMI. See [AWS documentation](http://docs.aws.amazon.com/cli/latest/reference/ec2/describe-images.html) for a list of possible filters and values.

## Behaviour

#### `check`

Searches the provided region for AMIs that match the configured filters. Versions are determined by AMI ID, and ordered by creation date.

#### `get`

Fetches additional metadata about the AMI, in addition to two files:

- `id`: Plain text file with the AMI ID.
- `packer`: Packer friendly variable file: `{"source_ami": "<ami-id>"}`.

#### `put`

Not implemented.

## Example

The following example would check for new versions of Amazon Linux 2 every 1h
and trigger the `bake-concourse` job whenever a new AMI (in Ireland) was found:

```yaml
resource_types:
- name: ami
  type: docker-image
  source:
    repository: itsdalmo/ami-resource

resources:
- name: amazon2-ami
  type: ami
  check_every: 1h
  source:
    aws_access_key_id: ((aws-access-key))
    aws_secret_access_key: ((aws-secret-key))
    aws_session_token: ((aws-session-token))
    aws_region: eu-west-1
    filters:
      name: "amzn2-ami-hvm*gp2"
      owner-id: "137112412989"
      architecture: "x86_64"
      virtualization-type: "hvm"
      root-device-type: "ebs"

jobs:
- name: bake-concourse
  plan:
  - get: amazon2-ami
    trigger: true
    ...
```
