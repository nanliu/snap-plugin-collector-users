# snap collector plugin - users

This plugin has ability to gather the number of logged-in users. Though this type of information can be obtained from various files in the Linux system but in this plugin there is a command line utility 'who' to get the list of user logged in the system.
															
The plugin is used in the [snap framework] (http://github.com/intelsdi-x/snap).				

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
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

- Linux system

### Installation

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

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:
                                                                                                
Metric namespace is `/intel/users/<metric name>/`

Metric Name | Data Type | Description
------------ | ------------- | -------------
logged | uint64 | A number of users logged-in
logged_min | uint64 | A minimum number of logged-in users<sup>(*)</sup>
logged_max | uint64 | A maximum number of logged-in users<sup>(*)</sup>
logged_avg | float64 | An average number of logged-in users<sup>(*)</sup>

<sup>(*)</sup> since the task was started																																					
By default metrics are gathered once per second.

### Examples

Example of running snap users collector and writing data to file.

Run the snap daemon:
```
$ snapd -l 1 -t 0
```

Load users plugin for collecting:
```
$ snapctl plugin load $SNAP_USERS_PLUGIN_DIR/build/rootfs/snap-plugin-collector-users
Plugin loaded
Name: users
Version: 1
Type: collector
Signed: false
Loaded Time: Tue, 12 Jan 2016 05:25:35 EST
```

See available metrics:
```
$ snapctl metric list
```

Load file plugin for publishing:
```
$ snapctl plugin load $SNAP_DIR/build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Tue, 12 Jan 2016 05:26:21 EST
```

Create a task JSON file (exemplary file in examples/tasks/users-file.json):  
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
                "/intel/users/logged": {},
                "/intel/users/logged_avg": {},
                "/intel/users/logged_max": {},
                "/intel/users/logged_min": {}
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
$ snapctl task create -t $SNAP_USERS_PLUGIN_DIR/examples/tasks/users-file.json
Using task manifest to create task
Task created
ID: b0ba4a9c-9011-4117-bd25-803b732f8e5b
Name: Task-b0ba4a9c-9011-4117-bd25-803b732f8e5b
State: Running
```
See sample output from `snapctl task watch <task_id>`

```
$ snapctl task watch  b0ba4a9c-9011-4117-bd25-803b732f8e5b
																								
Watching Task (b0ba4a9c-9011-4117-bd25-803b732f8e5b):
NAMESPACE                        DATA    TIMESTAMP                                       SOURCE
/intel/users/logged              6       2016-01-12 05:33:43.275317916 -0500 EST         gklab-108-166
/intel/users/logged_avg          5.48    2016-01-12 05:33:43.275317101 -0500 EST         gklab-108-166
/intel/users/logged_max          6       2016-01-12 05:33:43.275317675 -0500 EST         gklab-108-166
/intel/users/logged_min          5       2016-01-12 05:33:43.275317794 -0500 EST         gklab-108-166
```
(Keys `ctrl+c` terminate task watcher)


These data are published to file and stored there (in this example in /tmp/published_users).

Stop task:
```
$ snapctl task stop b0ba4a9c-9011-4117-bd25-803b732f8e5b
Task stopped:
ID: b0ba4a9c-9011-4117-bd25-803b732f8e5b
```

### Roadmap

There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-users/issues).

## Community Support
This repository is one of **many** plugins in the **Snap Framework**: a powerful telemetry agent framework. To reach out on other use cases, visit:

* [Snap Gitter channel] (https://gitter.im/intelsdi-x/snap)

The full project is at http://github.com:intelsdi-x/snap.

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
List authors, co-authors and anyone you'd like to mention

* Author: 	[Izabella Raulin](https://github.com/IzabellaRaulin)

**Thank you!** Your contribution is incredibly important to us.
