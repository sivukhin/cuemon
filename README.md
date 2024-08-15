# cuemon

**cuemon** is a tool for managing Grafana dashboards and panels using CUE language modules. It provides functionality for bootstrapping and updating monitoring configurations by converting Grafana JSON dashboards and panels into CUE format.

## Installation

To install cuemon, follow these steps:

1. Install the Go programming language if you haven't already
2. Run the following command to install cuemon:
```bash
$> go install github.com/sivukhin/cuemon
```

## Usage

**cuemon** provides two main commands: *bootstrap* and *update*.

### bootstrap

The bootstrap command is used to initialize monitoring configurations from a Grafana dashboard JSON file. It converts the dashboard JSON into CUE format.

```bash
$> cuemon bootstrap -input <dashboard.json> -module <module-name> -dir <output-directory> -overwrite
```

- **-input**: Path to the input Grafana dashboard JSON file.
- **-module**: Name of the CUE module to create.
- **-dir**: Target directory where the cuemon setup will be initialized.
- **-overwrite**: Enable unsafe mode to allow overwriting files.

### update

The update command updates existing monitoring configurations with new panels. It updates CUE files by converting the new panel JSON into CUE format and inserting it into the appropriate location.

```bash
$> cuemon update -input <panel.json> -dir <cuemon-directory> -overwrite
```

- **-input**: Path to the input Grafana panel JSON file.
- **-dir**: Directory containing the cuemon setup to update.
- **-overwrite**: Enable unsafe mode to allow overwriting files.

## Maintenance

`grafanaVxx` is generated with the help of [cuebootstrap](https://github.com/sivukhin/cuebootstrap) tool from the set of cherry-picked Grafana JSONs like this:
```bash
$> cd bootstrap
$> make cuebootstrap-install
$> make bootstrap-grafana-v10
```

After this you can replace `#grafana` schema definition with the output of cuebootstrap written in the `grafanaV10.cue` file
