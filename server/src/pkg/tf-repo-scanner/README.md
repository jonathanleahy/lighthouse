# Terraform Repository Scanner

This tool scans Terraform files to extract information about repository configurations and outputs the results in either JSON or CSV format. It can process both local directories and GitHub repositories.

## Features

- Scans Terraform files recursively in a directory
- Extracts repository names, team information, and descriptions
- Supports both local directories and GitHub repositories
- Outputs results in JSON or CSV format
- Automatically handles repository cloning and updates
- Uses directory names as team identifiers

## Prerequisites

- Go 1.16 or higher
- Git (for GitHub repository functionality)

## Getting Started Quickly

1. Build the program:
```bash
go build -o terraform-repo-scanner
```

2. Run it on a local directory:
```bash
./terraform-repo-scanner -path ./your-terraform-files
```

Or with a GitHub repository:
```bash
./terraform-repo-scanner -repo https://github.com/org/repo
```

## Installation

```bash
git clone https://github.com/yourusername/terraform-repo-scanner
cd terraform-repo-scanner
go build
```

## Usage

### Basic Command Structure

```bash
./terraform-repo-scanner [flags]
```

### Available Flags

- `-path`: Path to the root directory containing Terraform files (default: ".")
- `-format`: Output format, either "json" or "csv" (default: "json")
- `-repo`: GitHub repository URL (optional)
- `-tmp`: Temporary directory for cloning repositories (default: system temp directory)

### Examples

1. Process local directory:
```bash
./terraform-repo-scanner -path ./terraform-files -format json
```

2. Process GitHub repository:
```bash
./terraform-repo-scanner -repo https://github.com/org/repo -format json
```

3. Process GitHub repository with custom temp directory:
```bash
./terraform-repo-scanner -repo https://github.com/org/repo -tmp ./my-temp -format csv
```

### Example Output

JSON format:
```json
{
  "repositories": [
    {
      "repository_name": "example-repo",
      "team": "team-name",
      "description": "Repository description"
    }
  ],
  "total_count": 1
}
```

CSV format:
```
repository_name,team,description
------------------------------------
example-repo,team-name,Repository description
```

## How It Works

1. The tool walks through the specified directory (or cloned repository) recursively
2. For each `.tf` file found:
    - Extracts module blocks containing repository configurations
    - Uses the parent directory name as the team name
    - Parses repository name and description from the module block
3. Outputs the collected information in the specified format

## Notes

- When using the `-repo` flag, the tool will:
    - Clone the repository if it's not already present in the temp directory
    - Pull latest changes if the repository already exists
    - Process the files as it would for a local directory
- The tool uses the parent directory name of each Terraform file as the team name
- JSON output includes a total count of repositories found

## License

[Your chosen license]