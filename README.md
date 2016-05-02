# snap collector plugin - users

This plugin has ability to gather the number of logged-in users. Though this type of information can be obtained from various files in the Linux system but in this plugin there is a command line utility 'who' to get the list of user logged in the system.
															
The plugin is used in the [snap framework] (http://github.com/intelsdi-x/snap).				

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements
* Linux system with available `who` tool at `/usr/bin/who`
 
### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation

#### Download the plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page. Download the plugins package from the latest release, unzip and store in a path you want `snapd` to access.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-users

Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-users.git
```

Build the snap users plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage

* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-users/blob/master/README.md#examples).

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Metric Namespace | Data Type | Description
------------ | ------------- | -------------
/intel/utmp/users/logged | uint64 | A number of users logged-in
/intel/utmp/users/logged_min | uint64 | A minimum number of logged-in users<sup>(*)</sup>
/intel/utmp/users/logged_max | uint64 | A maximum number of logged-in users<sup>(*)</sup>
/intel/utmp/users/logged_avg | float64 | An average number of logged-in users<sup>(*)</sup>

<sup>(*)</sup> since the task was started																																					
By default metrics are gathered once per second.

### Examples

Example of running snap users collector and writing data to file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```
In another terminal window:

Load users plugin for collecting:
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-users
Plugin loaded
Name: users
Version: 1
Type: collector
Signed: false
Loaded Time: Tue, 12 Jan 2016 05:25:35 EST
```

See available metrics:
```
$ $SNAP_PATH/bin/snapctl metric list
```

Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Tue, 12 Jan 2016 05:26:21 EST
```

Create a task manifest file to use snap-plugin-collector-users plugin (exemplary file in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-users/blob/master/examples/tasks/)):
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/utmp/users/logged": {},
                "/intel/utmp/users/logged_avg": {},
                "/intel/utmp/users/logged_max": {},
                "/intel/utmp/users/logged_min": {}
            },
            "config": {
            },
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_users"
                    }
                }
            ]
        }
    }
}
    
```

Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t users-file.json
Using task manifest to create task
Task created
ID: b0ba4a9c-9011-4117-bd25-803b732f8e5b
Name: Task-b0ba4a9c-9011-4117-bd25-803b732f8e5b
State: Running
```
See sample output from `snapctl task watch <task_id>`

```
$ $SNAP_PATH/bin/snapctl task watch  b0ba4a9c-9011-4117-bd25-803b732f8e5b
																								
Watching Task (b0ba4a9c-9011-4117-bd25-803b732f8e5b):
NAMESPACE                           DATA    TIMESTAMP                                   SOURCE
/intel/utmp/users/logged            6       2016-01-12 05:33:43.275317916 -0500 EST     gklab-108-166
/intel/utmp/users/logged_avg        5.48    2016-01-12 05:33:43.275317101 -0500 EST     gklab-108-166
/intel/utmp/users/logged_max        6       2016-01-12 05:33:43.275317675 -0500 EST     gklab-108-166
/intel/utmp/users/logged_min        5       2016-01-12 05:33:43.275317794 -0500 EST     gklab-108-166
```
(Keys `ctrl+c` terminate task watcher)


These data are published to file and stored there (in this example in /tmp/published_users).

Stop previously created task:
```
$ $SNAP_PATH/bin/snapctl task stop b0ba4a9c-9011-4117-bd25-803b732f8e5b
Task stopped:
ID: b0ba4a9c-9011-4117-bd25-803b732f8e5b
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-users/issues) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-users/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [snap Gitter channel](https://gitter.im/intelsdi-x/snap).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: 	[Izabella Raulin](https://github.com/IzabellaRaulin)

