#!/usr/bin/env python3

from withoutbg import WithoutBG
from typing import Union
import sys

def printerr(v) -> None:
    print(v, file=sys.stderr)

def trim_ext(target: str) -> str:
    parts = target.split(".")
    if len(parts) != 2 or (parts[0] == "" or parts[1] == ""):
        raise ValueError(f"invalid input string:\n{target}")
    return parts[0]

def process(target) -> Union[str, None]:
    """process removes the background of the target image.
    Returns output_file_name: str if successful
    otherwise, returns None"""
    trimmed_target = ""
    try:
        trimmed_target = trim_ext(target)
    except ValueError as e:
        printerr(f"error processing target file:{e}")
        return None

    output_file_name = f"{trimmed_target}-trans.png"

    bg_remover = WithoutBG.opensource()
    result = bg_remover.remove_background(target)

    try:
        result.save(output_file_name)
    except ValueError as ve:
        printerr(ve)
        return None
    except OSError as ose:
        printerr(ose)
        return None

    return output_file_name

if __name__ == "__main__":
    num_args = len(sys.argv)

    if num_args != 2:
        printerr("invalid number of args")
        sys.exit(1)   

    input_file = sys.argv[1]

    res = process(input_file)

    if res is None:
        sys.exit(1)

    print(res)
