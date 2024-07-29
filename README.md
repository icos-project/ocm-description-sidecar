# OCM Descriptor Sidecar

OCM Descriptor Sidecar is a component responsible for triggering the execution of jobs in OCM Description Service and updating the status of all deployed resources into the Job Manager (JM). This sidecar is designed to work with OCM Description Service and uses Keycloak for authentication.

## Features

- Triggers the execution of jobs.
- Updates the status of all deployed resources into JM periodically.
- Caches authentication tokens in memory to avoid requesting multiple tokens at a time.

## Kind Installation

Please, refer to the helm suite in [ICOS Agent Repository](https://production.eng.it/gitlab/icos/suites/icos-agent)

## Usage

The sidecar triggers the execution of jobs and updates the status of deployed resources. The primary function `Schedule` is responsible for this process.

### Schedule Function

The `Schedule` function performs the following steps:

1. **Trigger Job Execution**: Sends a request to the deployment manager to start the execution of jobs.
2. **Fetch Token**: Obtains a Keycloak token for authentication.
3. **Debug Logging**: Logs the request and response for debugging purposes.
4. **Trigger Resource Sync**: Sends a request to the deployment manager to update the status of all deployed resources into JM periodically.


## Contributing

In order to contribute to this repository, feel free to open a pull request and assign `@x_alvolkov`or `x_magallar` as a reviewer.

## Legal

The OCM Descriptor Sidecar is released under the Apache 2.0 license.
Copyright © 2022-2024 Eviden. All rights reserved.

🇪🇺 This work has received funding from the European Union's HORIZON research and innovation programme under grant agreement No. 101070177.
