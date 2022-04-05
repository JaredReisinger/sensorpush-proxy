# sensorpush-proxy

A rate-limiting, authentication-hiding proxy for SensorPush data

----

_(from puppycam-sensor...)_

In order to get the config into the deployed container without putting it in
plaintext in the source, we end-run through the Drone "encrypted secret"
mechanism, and decode/expand it during the build.

* config.yaml -> base64 encoded (one line) -> drone encrypt

```bash
drone encrypt JaredReisinger/puppycam-sensor $(cat config.yaml | base64 -w0)
```

This (as of this moment) results in:

```text
0z6PxJVqtfoSWiUb+831gUC1jViO5zctZkwyqTJ7AfUY1c8Vmc1v3oiCgAWqAr/qH6ZXaC/H6CiOI9fjkVlOM6XafVeUi19kkq6dESzZxvdv1+y6MVA8jy9+7olrcQeagu4PQ0JjtbYdgIUTuPpBT5LR2lQEPKZqJHGVYeAzLsSOwIastH/M8oLihGbDYTu/cBsCxAEgAMTq11jz4vn9/eThEZQ5PbkmglP1ww==
```

Then, in [the drone file](./.drone.yml), we include the above string as a secret, put the value into an environment variable (called `CONFIG_YAML`), and re-hydrate as a local file:

```bash
printf "$${CONFIG_YAML}" | base64 -d > config.yaml
```

_Note that the double-`$` is needed to prevent Drone from attempting to expand the value itself._
