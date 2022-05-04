# ACK service controller for EMR on EKS

This repository contains source code for the AWS Controllers for Kubernetes
(ACK) service controller for EMR on EKS Service.

Please [log issues](https://github.com/aws-controllers-k8s/community/issues) and feedback on the main AWS Controllers for Kubernetes (ACK) Github project.

## Overview

The ACK service controller for EMR on EKS provides declarative way to run spark jobs on EKS clusters. EMR on EKS manages the lifecycle of these jobs and it uses highly optimized EMR runtime [3.5 times faster than open-source Spark](https://aws.amazon.com/blogs/big-data/amazon-emr-on-amazon-eks-provides-up-to-61-lower-costs-and-up-to-68-performance-improvement-for-spark-workloads/) when you run these jobs. For more information about EMR on EKS, please read our [documentation](https://docs.aws.amazon.com/emr/latest/EMR-on-EKS-DevelopmentGuide/emr-eks.html)

## Contributing

We welcome community contributions and pull requests.

See our [contribution guide](/CONTRIBUTING.md) for more information on how to
report issues, set up a development environment, and submit code.

We adhere to the [Amazon Open Source Code of Conduct][coc].

You can also learn more about our [Governance](/GOVERNANCE.md) structure.

[coc]: https://aws.github.io/code-of-conduct

## License

This project is [licensed](/LICENSE) under the Apache-2.0 License.
