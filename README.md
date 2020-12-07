# Perímener ✋

[![GitHub license](https://img.shields.io/badge/license-Apache%20License%202.0-blue.svg)](LICENSE)

Named after the Greek word **Perímene** _αναμονή_ for wait.

A golang application that waits for the expected number of pods to reach a `Ready` state before gracefully exiting.

# Problem statement

Have you ever needed an `initContainer` which can perform simple dependency management. 

Starting containers before their dependencies are ready can result in pods going into a status of `CrashLoopBackoff` these wastes resources which could have been used on the dependencies themselves.

An example of this might be ensuring that a number of zookeeper pods reach a `Ready` state before kafka is started.

Perímener also has functionality to avoid the thundering herd issue which is amplified when using GitOps principles. 

Setting  an environment variable inserts a random delay from 0 to `<RANDOM_WINDOW_SECONDS>` to application startup which can help prevent undesirable autoscaling events.

# Usage

Configuration is done via the use of ENVIRONMENT variables

|Environment variable Name|Parameter Function|Required/Default|
|------------------------|---------|-------------|
|DELAY\_SECONDS|Delay start of the main container by n seconds | Default: 0 seconds (Disabled) |
|EXPECTED\_READY\_POD\_COUNT| The number of pods to reach a `Ready` state before the init container will exit and allow the main container to start| Required |
|NAMESPACE|Which namespace to target| Required |
|POD\_LABEL|A label selector which matches the required pods| Required |
|RANDOM\_WINDOW\_SECONDS|Max seconds for window size | Default: 0 seconds (Disabled) |
|SLEEP\_COUNT|The time to wait between retries| Default: 5 seconds |
|USE\_LOCAL\_KUBECONFIG| Use a local kubeconfig or not. Defaults to in cluster auth| Default: false|
