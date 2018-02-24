#!/usr/bin/python3
# -*- encoding: utf-8 -*-

import csv
import openpyxl
import sys

wb = openpyxl.Workbook()
ws = wb.active

f = open(sys.argv[1])
reader = csv.reader(f, delimiter=';', quotechar='"')
for row in reader:
    ws.append(row)
f.close()

wb.save(sys.argv[2])
