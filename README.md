# Terraform Lock Manager

## Overview

The DynamoDB Lock Manager is a command-line tool designed to manage and unlock Terraform state locks directly in AWS DynamoDB. This utility allows users to interact with DynamoDB to unlock state locks without the need to initialize Terraform or be in a specific directory. It's particularly useful in scenarios where Terraform state locks need to be managed independently of Terraform CLI operations.

## Features

- **Unlock State Locks**: Allows users to delete Terraform state locks from DynamoDB based on a provided regex pattern.
- **Flexible Usage**: Can be used from any directory, independent of Terraform's current state or initialization.
- **Regex Support**: Supports regular expressions for identifying lock IDs, providing flexibility in selecting the locks to be removed.
- **Persistent Configuration**: Set your terraform config in the `~/.terraform.d/lock_manager.toml` file.


## Installation


```shell
git clone [repository-url]
cd [repository-directory]
go build -o lockmanager
```

## Usage

Your first run will configure the application.

```toml
table = "my-terraform-lock-table"
region = "us-west-2"
profile = "develeopment"
```

Run the tool using the following command format:

```shell
./lockmanager unlock [Regex-Pattern]
```

- `[Regex-Pattern]`: A regular expression pattern to match the lock IDs that need to be unlocked.

Example:

```shell
./lockmanager unlock "long/path-to-my/service/subkeys.*"
```

This command will remove all locks from the `my-terraform-lock-table` table that match the regex pattern `"long/path-to-my/service/subkeys.*"`
