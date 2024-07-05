# GCP IAM Auditor

This Go program lists all admin roles assigned to user accounts (excluding specified member types) within specified GCP projects. The roles are grouped by project, member type, and user, and displayed in a stylish table with roles separated by commas. The program also includes a progress tracker that displays individual project processing lines along with an overall progress tracker.

## Prerequisites

- Ensure `gcloud` CLI is installed and authenticated.
- Ensure Go is installed.
- Install the required Go packages:

  ```sh
  go get -u github.com/jedib0t/go-pretty/v6/progress
  go get -u github.com/jedib0t/go-pretty/v6/table
  go get -u gopkg.in/yaml.v2
  ```

## Installation

1. Clone this repository or save the script to a directory named `gcpiamauditor`.

2. Build the Go program:

   ```sh
   go build -o gcpiamauditor gcpiamauditor/main.go gcpiamauditor/policy.go gcpiamauditor/roles.go gcpiamauditor/output.go
   ```

## Usage

Create a file named `projects.txt` with each project ID on a new line.

Example `projects.txt`:

```
project-id-1
project-id-2
project-id-3
```

Run the program with the project file as an argument:

```sh
./gcpiamauditor -projectFile=projects.txt [-exclude=serviceAccount,user,group] [-services=bigquery] [-member=admin,owner,editor] [-output=table,json,yaml,csv]
```

Use the `-exclude` argument to specify member types to exclude (e.g., `serviceAccount`, `user`, `group`). Use the `-services` argument to specify GCP services to check (e.g., `bigquery`). Use the `-member` argument to specify IAM roles to include (e.g., `admin`, `owner`, `editor`). Use the `-output` argument to specify the output format (e.g., `table`, `json`, `yaml`, `csv`).

## Script Description

- **Command-Line Argument Parsing**: The project file, member types to exclude, GCP services to check, IAM roles to include, and output format are passed as command-line arguments (`-projectFile`, `-exclude`, `-services`, `-member`, and `-output`).
- **Reading Project File**: The `projects.txt` file is read to get a list of project IDs.
- **Progress Tracker**: Uses the `github.com/jedib0t/go-pretty/v6/progress` package to display a progress tracker while fetching and processing IAM policies.
- **IAM Policy Retrieval and Role Collection**: For each project ID, the IAM policy is fetched using the `gcloud` CLI, and roles are collected into a map, excluding specified member types and filtering by specified services and IAM roles.
- **Stylish Table Output**: Uses the `github.com/jedib0t/go-pretty/v6/table` package to display the results in a stylish table with colored output.

## Example Output

The program displays the following output in table format:

```
+---------+-------------+-------------------+---------------------------------------------+
| PROJECT | MEMBER TYPE |       USER        |                   ROLES                    |
+---------+-------------+-------------------+---------------------------------------------+
| p1      | user        | user1@example.com | roles/editor, roles/owner, roles/viewer    |
| p1      | serviceAccount | sa1@example.com | roles/editor, roles/owner, roles/viewer    |
| p2      | user        | user2@example.com | roles/viewer                               |
+---------+-------------+-------------------+---------------------------------------------+
```

## Error Handling

If an error occurs while fetching the IAM policy for a project, the tracker will be marked as errored with a detailed message, and the program will continue processing the remaining projects.

## License

This project is licensed under the MIT License.