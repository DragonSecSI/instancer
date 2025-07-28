import argparse
import sys
import yaml
import requests

chall_types = {
    "web": 0,
    "socket": 1,
}
flag_types = {
    "static": 0,
    "suffix": 1,
    "leetify": 2,
    "capitalize": 3,
    "combined": 4,
}

if __name__ == "__main__":
    argparser = argparse.ArgumentParser(description="Instancer CLI")
    argparser.add_argument("--api", type=str, default="https://instancer.vuln.si", help="API URL")
    argparser.add_argument("--token", type=str, default="admin", help="API token for authentication")
    argparser.add_argument("--config", type=str, required=True, help="Path to the configuration file")

    args = argparser.parse_args()

    try:
        challenge = open(args.config, "r").read()
        challenge = yaml.safe_load(challenge)
        if not challenge:
            print("Error: The configuration file is empty or not properly formatted.")
            sys.exit(1)
    except FileNotFoundError:
        print(f"Error: The configuration file '{args.config}' does not exist.")
        sys.exit(1)
    except Exception as e:
        print(f"An error occurred while reading the file: {e}")
        sys.exit(1)

    headers = {
        "Authorization": args.token,
    }
    payload = {
        "name": challenge["name"],
        "description": challenge["description"],
        "type": chall_types[challenge["type"]],
        "flag": challenge["flag"],
        "flag_type": flag_types[challenge["flag_type"]],
        "duration": challenge.get("duration", 1800),
        "repository": challenge.get("repository", "oci://registry:5000/charts"),
        "chart": challenge["chart"],
        "chart_version": challenge["chart_version"],
        "values": challenge.get("values", "").strip(),
    }
    try:
        response = requests.post(f"{args.api}/api/v1/challenge/", headers=headers, json=payload)
        response.raise_for_status()
        print("Challenge created successfully.")
    except requests.exceptions.HTTPError as http_err:
        print(f"HTTP error occurred: {http_err}")
    except requests.exceptions.RequestException as req_err:
        print(f"Request error occurred: {req_err}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
