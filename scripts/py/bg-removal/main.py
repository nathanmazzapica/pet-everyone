#!/usr/bin/env python3
from sympy.utilities.decorator import deprecated
from withoutbg import WithoutBG
from typing import Union
import sys

def printerr(v) -> None:
    print(v, file=sys.stderr)

def process(target, output_path):
    """process removes the background of the target image.
    :param target: the target image file path
    :param output_path: the output path"""

    bg_remover = WithoutBG.opensource()
    result = bg_remover.remove_background(target)

    result.save(output_path)


if __name__ == "__main__":
    num_args = len(sys.argv)

    if num_args != 3:
        printerr(f"invalid number of args: {num_args-1}, expected 2")
        sys.exit(1)   

    input_file = sys.argv[1]
    outputh_path = sys.argv[2]

    try:
        process(input_file, outputh_path)
    except Exception as e:
        printerr(f"error processing file:{e}")
        sys.exit(1)

    sys.exit(0)