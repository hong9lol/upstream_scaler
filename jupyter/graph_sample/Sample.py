# -*- coding: utf-8 -*-


import matplotlib.pyplot as plt
import numpy as np

import sys

sys.path.append(r"H:\Python")

import Preset

Preset.setPlotStyle()

data = np.genfromtxt(r"data.txt", skip_header=1)

x = data[:, 0]
y = data[:, 1]

plt.plot(x, y, "^", label="Raw Data", color="red")

linfit = np.polyfit(x, y, 1)
f_1d = np.poly1d(linfit)

plt.plot(x, f_1d(x), label="Linear Fit", color="black")

plt.legend(loc="best")
# plt.legend(loc = (1.02, 0.02))

plt.xlabel("X [-]", fontsize=20)
plt.ylabel("Y [-]", fontsize=20)

# data = np.genfromtxt(r'data.txt', skip_header = 1)

# x = data[:,0]
# y = data[:,1]

# plt.plot(x,y, 'o',  label = 'Raw Data', color = 'red')

# linfit = np.polyfit(x,y,1)
# f_1d = np.poly1d(linfit)

# # plt.plot(x,f_1d(x), label = 'Linear Fit', color = 'black')

# plt.plot(x,f_1d(x), label = 'Linear Fit', color = 'black')

# plt.xlabel('X [-]', fontsize=20)
# plt.ylabel('Y [-]', fontsize=20)
# plt.legend(loc = 'best')
# plt.legend(loc = (1.02,0.02))


# https://dong2kim.tistory.com/entry/Python-Matplotlib%EC%9D%84-%EC%82%AC%EC%9A%A9%ED%95%B4%EC%84%9C-%EB%85%BC%EB%AC%B8%EC%97%90-%EB%93%A4%EC%96%B4%EA%B0%88-%EA%B7%B8%EB%9E%98%ED%94%84-%EC%99%84%EC%84%B1%EB%8F%84%EC%9E%88%EA%B2%8C-%EA%B7%B8%EB%A6%AC%EA%B8%B0-1
