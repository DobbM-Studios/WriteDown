import os
import shutil
import subprocess
import sys

def main():
    print("=== 1. Compilando Go a WebAssembly en la Raíz ===")
    env = os.environ.copy()
    env["GOOS"] = "js"
    env["GOARCH"] = "wasm"

    # Compilamos directamente a ./main.wasm
    cmd_build = ["go", "build", "-o", "main.wasm", "src/main.go"]
    res_build = subprocess.run(cmd_build, env=env, capture_output=True, text=True)
    
    if res_build.returncode != 0:
        print(f"❌ Error al compilar Go:\n{res_build.stderr}")
        sys.exit(1)
    print("✔ 'main.wasm' generado en la raíz.")

    print("\n=== 2. Extrayendo archivo de soporte wasm_exec.js ===")
    try:
        goroot = subprocess.run(["go", "env", "GOROOT"], capture_output=True, text=True, check=True).stdout.strip()
        
        src_js = os.path.join(goroot, "lib", "wasm", "wasm_exec.js")
        if not os.path.exists(src_js):
            src_js = os.path.join(goroot, "misc", "wasm", "wasm_exec.js")

        # Copiamos directamente a ./wasm_exec.js
        shutil.copy(src_js, "wasm/wasm_exec.js")
        print("✔ 'wasm_exec.js' copiado a la raíz.")
    except Exception as e:
        print(f"❌ Error al buscar wasm_exec.js: {e}")
        sys.exit(1)

    print("\n🎉 ¡Proceso completado en la raíz!")

if __name__ == "__main__":
    main()
