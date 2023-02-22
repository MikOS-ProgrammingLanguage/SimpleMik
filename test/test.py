import os
import sys
import subprocess

COUNT: int = 0
FAILED_COUNT: int = 0
FAILING_FILES = []

def TestFile(file_location: str):
    if file_location.endswith(".mik") or file_location.endswith(".milk"):
        global COUNT, FAILED_COUNT, FAILING_FILES
        return_code = subprocess.getoutput("mic -i "+file_location)
        if return_code.count("[SUCCESS]") != 1:
            print(f"❌: {COUNT+1} -> {file_location}")
            FAILING_FILES.append(file_location)
            FAILED_COUNT += 1
        else:
            print(f"✅: {COUNT+1} -> {file_location}")
        COUNT += 1

def TestDir(directory: str):
    for i in os.listdir(directory):
        if i == (__file__.split("/")[len(__file__.split("/"))-1]):
            continue
        if os.path.isdir(os.path.join(directory, i)):
            TestDir(directory+"/"+i)
        else:
            TestFile(directory+"/"+i)

if __name__=="__main__":
    if len(sys.argv) > 1:
        TestDir(sys.argv[1])
    else:
        TestDir(".")
    
    try:
        os.remove("mik.out")
        os.remove("ir.ll")
    except:
        pass
    
    print(f"\n\nFiles Scanned:        {COUNT}\n")
    print(f"Total failed tests: {FAILED_COUNT}")
    print(f"Total Successful:   {COUNT - FAILED_COUNT}")
    
    if len(FAILING_FILES) > 0:
        print("\nFailing Files:")
        for i in FAILING_FILES:
            print(f"\t- {i}")
