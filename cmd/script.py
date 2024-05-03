import os

def process_files(directory, output_file):
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith(".go"):
                file_path = os.path.join(root, file)
                with open(file_path, 'r') as f:
                    content = f.read()
                    output_file.write(f"__________ {file_path}:\n\n")
                    output_file.write(content)
                    output_file.write("\n\n")

# Directory to start the recursive search
start_directory = "."

# Output file path
output_file_path = "output.txt"

# Open the output file in write mode
with open(output_file_path, 'w') as output_file:
    # Start processing files recursively
    process_files(start_directory, output_file)

print(f"File contents have been written to {output_file_path}")
