import yaml
import sys

# The path to your YAML file
yaml_file_path = "DeathStarBench/socialNetwork/helm-chart/socialnetwork/values.yaml"

# The key you want to change and the new value
key_to_change = "limitCpu"

print(sys.argv[1])
new_value = str(sys.argv[1]) + "m"

# Step 1: Read the YAML file
with open(yaml_file_path, "r") as file:
    yaml_content = yaml.safe_load(file)

# Step 2: Update the value for the key
yaml_content["global"][key_to_change] = new_value

# Step 3: Write the updated content back to the YAML file
with open(yaml_file_path, "w") as file:
    yaml.safe_dump(yaml_content, file)
