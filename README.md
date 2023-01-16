> # THIS README IS A WORK-IN-PROGRESS!


![sensorpush-proxy](./docs/sensorpush-proxy-logo.png)

A rate-limiting, authentication-hiding proxy for [SensorPush](https://www.sensorpush.com) data.

## Usage

More than likely, you’ll want to use the already-built [Docker image](https://hub.docker.com/r/jaredreisinger/sensorpush-proxy) for this, so that you can just throw your credentials and sensor config at it and let it run behind something like Traefik or nginx.  Regardless of whether you use Docker or the raw binary, `sensorpush-proxy` includes both the proxy itself and also a `query` sub-command to help you discover the sensor IDs available to you.

### Get SensorPush credentials

If you have any SensorPush devices, you ought to already have an account with SensorPush.  In order to activate API access, you need to sign into their [Gateway Cloud Dashboard](https://dashboard.sensorpush.com/) at least once and agree to the terms of service.  Once you've done that, the username and password you use for that dashboard are the same ones you need for this proxy.  Unfortunately, SensorPush does not allow you to create limited-use (and revocable) tokens to use in place of your password for API access.  This proxy _**only**_ uses the password you provide in order to retrieve the sensor data you request.  Use secure best practices to ensure that the configuration you pass to this proxy is kept secret/encrypted until the last possible moment.

### Get sensorpush-proxy




### Discover sensor IDs

Use the `query` subcommand


### Docker image


### Binary (from release or source)




Uses the [SensorPush API](https://www.sensorpush.com/gateway-cloud-api) to fetch data.

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
