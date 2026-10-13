import os
import shutil
import subprocess
import sys

def main():
    print("=== Compiling Go to WASM ===")
    env = os.environ.copy()
    env["GOOS"] = "js"
    env["GOARCH"] = "wasm"

    cmd_build = ["go", "build", "-o", "main.wasm", "src/main.go"]
    res_build = subprocess.run(cmd_build, env=env, capture_output=True, text=True)
    
    if res_build.returncode != 0:
        print(f"Error compiling Go:\n{res_build.stderr}")
        sys.exit(1)
    print("main.wasm generated")

    print("\n=== Extracting wasm_exec.js ===")
    try:
        goroot = subprocess.run(["go", "env", "GOROOT"], capture_output=True, text=True, check=True).stdout.strip()
        
        src_js = os.path.join(goroot, "lib", "wasm", "wasm_exec.js")
        if not os.path.exists(src_js):
            src_js = os.path.join(goroot, "misc", "wasm", "wasm_exec.js")

        shutil.copy(src_js, "wasm/wasm_exec.js")
        print("Done")
    except Exception as e:
        print(f"Error searching wasm_exec.js: {e}")
        sys.exit(1)

    print("\nProcess finished")

if __name__ == "__main__":
    main()
