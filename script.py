import os

def copy_go_files(folder_path, output_file):
    for root, dirs, files in os.walk(folder_path):
        for file in files:
            if file.endswith(".go") or file.endswith(".sql"):
                file_path = os.path.join(root, file)
                with open(file_path, "r") as go_file:
                    content = go_file.read()
                    output_file.write(f"------ {file_path}:\n\n")
                    output_file.write(content)
                    output_file.write("\n\n")

current_folder = os.getcwd()
output_file_path = os.path.join(current_folder, "output.txt")

with open(output_file_path, "w") as output_file:
    copy_go_files(current_folder, output_file)

print(f"Go files copied to {output_file_path}")
