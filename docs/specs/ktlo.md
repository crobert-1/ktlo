# KTLO

## Goal

The goal is to use golang scripts to automate the KTLO process for the OpenTelemetry Collector and its components,
to reduce the manual workload currently required to fulfill KTLO requirements.

## Implementation

### Technologies

- golang - Used to implement scripts and tools for the KTLO process.
- GraphQL - Used to query the GitHub/GitLab GraphQL API for metadata information.

### Functionalities

#### Central GoLang script

- Refer to [golang-script.md](./golang-script.md) for the script's spec.

#### PR status checks



#### Open vulnerability checks



#### Flaky test checks



### Repository structure

- One top-level golang script that can handle all the KTLO functionalities, with subcommands for each functionality.
- Separate packages or modules for each functionality to keep the code organized and maintainable.
- Use a shared internal package for shared functionality