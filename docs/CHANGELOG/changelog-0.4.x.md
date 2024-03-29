# Claudie `v0.4`

!!! warning "Due to a breaking change in the input manifest schema, the `v0.4.x` will not be backwards compatible with `v0.3.x`"

## Deployment

To deploy the Claudie `v0.4.X`, please:

1. Download the archive and checksums from the [release page](https://github.com/berops/claudie/releases)

2. Verify the archive with the `sha256` (optional)

    ```sh
    sha256sum -c --ignore-missing checksums.txt
    ```

    If valid, output is, depending on the archive downloaded

    ```sh
    claudie.tar.gz: OK
    ```

    or

    ```sh
    claudie.zip: OK
    ```

    or both.

3. Lastly, unpack the archive and deploy using `kubectl`

    > We strongly recommend changing the default credentials for MongoDB, MinIO and DynamoDB before you deploy it. To do this, change contents of the files in `mongo/secrets`, `minio/secrets` and `dynamo/secrets` respectively.

    ```sh
    kubectl apply -k .
    ```

## v0.4.0

### Features

- Input manifest definition now uses CRD instead of secret [#872](https://github.com/berops/claudie/pull/872)
- Various improvements in the overall documentation [#864](https://github.com/berops/claudie/pull/864) [#871](https://github.com/berops/claudie/pull/871) [#884](https://github.com/berops/claudie/pull/884) [#888](https://github.com/berops/claudie/pull/888)  [#891](https://github.com/berops/claudie/pull/891) [#893](https://github.com/berops/claudie/pull/893)

### Bugfixes

- Errors from the Scheduler are correctly saved under the clusters state [#868](https://github.com/berops/claudie/pull/868)
- Failure in the Terraformer will correctly saves the created state [#875](https://github.com/berops/claudie/pull/875)
- The clusters which previously resulted in error no longer block the workflow on input manifest reapply [#883](https://github.com/berops/claudie/pull/883)

### Known issues

- Single node pool definition cannot be used as control plane and as compute plane in the same cluster [#865](https://github.com/berops/claudie/issues/865)
- Input manifest status is not tracked during autoscaling [#886](https://github.com/berops/claudie/issues/886)

## v0.4.1

### Features

- Add support for static nodepools, enabling Claudie to create hybrid-cloud or private clusters [#902](https://github.com/berops/claudie/pull/902)
- Use new time format for Claudie generated logs [#919](https://github.com/berops/claudie/pull/919)
- Allow reuse of the same nodepool definition as a control plane and compute plane [#936](https://github.com/berops/claudie/pull/936)
- Rename Frontend service to Claudie-controller [#939](https://github.com/berops/claudie/pull/939)

### Bugfixes

- Fix Frontend not recognizing updates from Cluster Autoscaler [#901](https://github.com/berops/claudie/pull/901)
- Use keep-alive settings to minimise impact of network disturbance between Claudie services [#903](https://github.com/berops/claudie/pull/903)

### Known issues

No known issues since the last release


## v0.4.2

### Features

- Improved error logging when Claudie is deployed without debug logging.
- All resources created by claudie now have the `app.kubernetes.io/part-of: claudie` label.
- SecurityContext was added to Claudie deployed services.
- A limit was introduced to how many terraform, ansible, kubeone processes can be spawn at a given time when building multiple clusters.

### Bugfixes

- Fixed correct deletion of a InputManifests with multiple clusters where one of them failed.

### KnownIssues

No new known issues since the last release
