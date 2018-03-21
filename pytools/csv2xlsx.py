#!/usr/bin/python3
# -*- encoding: utf-8 -*-

import csv
import openpyxl
import sys

LETTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
def excel_style(row, col):
    """ Convert given row and column number to an Excel-style cell name. """
    result = []
    while col:
        col, rem = divmod(col-1, 26)
        result[:0] = LETTERS[rem]
    return ''.join(result) + str(row)

def excel_column(col):
    result = []
    while col:
        col, rem = divmod(col-1, 26)
        result[:0] = LETTERS[rem]
    return ''.join(result)

wb = openpyxl.Workbook()
ws = wb.active

f = open(sys.argv[1])
reader = csv.reader(f, delimiter=';', quotechar='"')
row = next(reader)
width = len(row)
ws.append(row)
height = 1

column_widths = []
for i, cell in enumerate(row):
    if len(column_widths) > i:
       if len(cell) > column_widths[i]:
          column_widths[i] = len(cell)
    else:
       column_widths += [len(cell)]
for row in reader:
    height += 1
    for i, cell in enumerate(row):
        if len(column_widths) > i:
            if len(cell) > column_widths[i]:
                column_widths[i] = len(cell)
        else:
            column_widths += [len(cell)]
    ws.append(row)
f.close()

ws.auto_filter.ref = "A1:" + excel_style(height, width)
for i, column_width in enumerate(column_widths):
    ws.column_dimensions[excel_column(i+1)].width = column_width
ws.freeze_panes = ws['C2']
wb.save(sys.argv[2])
